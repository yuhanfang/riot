// Package solver implements a strategy solver for optimal champion selection.
// The user is responsible for providing a pool of viable champions and a
// utility function for evaluating a set of champions.
//
// This package is still in heavy development, and is likely to be unstable.
package solver

import (
	"errors"
	"fmt"
	"sort"

	"github.com/yuhanfang/riot/constants/champion"
)

// ChampionSet is a set of champions. All champions in the map are assumed to
// be contained in the set.
type ChampionSet map[champion.Champion]bool

// ChampionSetOf constructs a ChampionSet from a list of champions.
func ChampionSetOf(vals []champion.Champion) ChampionSet {
	cs := make(map[champion.Champion]bool)
	for _, v := range vals {
		cs[v] = true
	}
	return cs
}

// Copy makes a copy of the ChampionSet.
func (c ChampionSet) Copy() ChampionSet {
	m := make(map[champion.Champion]bool)
	for k := range c {
		m[k] = true
	}
	return m
}

// Utility takes a ChampionSet and returns the utility.
type Utility func(ChampionSet) (float64, error)

// State is a serializable representation of all bans and picks up to the given
// evaluation point.
type State struct {
	// Actions are a sequence of picks and bans. Each player has one ban and one
	// pick. The order of actions is:
	//
	// Initial ban:
	//
	// 0  Blue ban
	// 1  Red ban
	// 2  Blue ban
	// 3  Red ban
	// 4  Blue ban
	// 5  Red ban
	//
	// Initial pick: for pairs of picks, champion IDs are sorted in ascending order.
	//
	// 6  Blue pick
	// 7  Red pick
	// 8  Red pick
	// 9  Blue pick
	// 10 Blue pick
	// 11 Red pick
	//
	// Second ban:
	//
	// 12 Red ban
	// 13 Blue ban
	// 14 Red ban
	// 15 Blue ban
	//
	// Second pick:
	//
	// 16 Red pick
	// 17 Blue pick
	// 18 Blue pick
	// 19 Red pick
	Actions [20]champion.Champion
}

// Normalize returns a copy of a State that is equivalent to the input state
// from the perspective of the next action. All champion IDs are sorted
// ascending when possible. For example, if blue side banned champions 2 and
// then 1, the bans will be normalized to 1 and then 2. Normalization improves
// cache hits by avoiding computation for strategically equivalent states. When
// returning the actual optimal play, the user must be careful to "denormalize"
// to get back the original pick order instead of the normalized order.
func (s State) Normalize() State {
	bans := []int{0, 1, 2, 3, 4, 5, 12, 13, 14, 15}
	blue := []int{6, 9, 10, 17, 18}
	red := []int{7, 8, 11, 16, 19}

	var banVals, blueVals, redVals champions

	for _, b := range bans {
		if s.Actions[b] == 0 {
			break
		}
		banVals = append(banVals, s.Actions[b])
	}
	for _, b := range blue {
		if s.Actions[b] == 0 {
			break
		}
		blueVals = append(blueVals, s.Actions[b])
	}
	for _, b := range red {
		if s.Actions[b] == 0 {
			break
		}
		redVals = append(redVals, s.Actions[b])
	}

	sort.Sort(banVals)
	sort.Sort(blueVals)
	sort.Sort(redVals)

	for i, v := range banVals {
		s.Actions[bans[i]] = v
	}
	for i, v := range blueVals {
		s.Actions[blue[i]] = v
	}
	for i, v := range redVals {
		s.Actions[red[i]] = v
	}
	return s
}

// Merge copies all nonzero states and copies them into the input state,
// returning the merged state as a copy.
func (s State) Merge(other State) State {
	cp := other
	for i, v := range s.Actions {
		if v == 0 {
			break
		}
		cp.Actions[i] = v
	}
	return cp
}

// Blue returns a set of all champions currently picked by Blue.
func (s State) Blue() ChampionSet {
	indices := []int{6, 9, 10, 17, 18}
	cs := make(map[champion.Champion]bool)
	for _, i := range indices {
		if s.Actions[i] != 0 {
			cs[s.Actions[i]] = true
		}
	}
	return cs
}

// Red returns a set of all champions currently picked by Red.
func (s State) Red() ChampionSet {
	indices := []int{7, 8, 11, 16, 19}
	cs := make(map[champion.Champion]bool)
	for _, i := range indices {
		if s.Actions[i] != 0 {
			cs[s.Actions[i]] = true
		}
	}
	return cs
}

// Bans returns a set of all champions currently banned by either side.
func (s State) Bans() ChampionSet {
	indices := []int{0, 1, 2, 3, 4, 5, 12, 13, 14, 15}
	cs := make(map[champion.Champion]bool)
	for _, i := range indices {
		if s.Actions[i] != 0 {
			cs[s.Actions[i]] = true
		}
	}
	return cs
}

// Unavailable returns a set of all champions currently picked or banned.
func (s State) Unavailable() ChampionSet {
	unavailable := ChampionSet(make(map[champion.Champion]bool))
	for _, v := range s.Actions {
		if v == 0 {
			break
		}
		unavailable[v] = true
	}
	return unavailable
}

// Available1 returns a slice containing all available champions within the pool.
func (s State) Available1(pool ChampionSet) []champion.Champion {
	var champs []champion.Champion
	unavail := s.Unavailable()
	for k := range pool {
		if !unavail[k] {
			champs = append(champs, k)
		}
	}
	return champs
}

// champions implements sort.Interface for sorting ascending by champion ID.
type champions []champion.Champion

// Len returns the slice length.
func (c champions) Len() int {
	return len(c)
}

// Less returns true if i has a lower champion ID than j.
func (c champions) Less(i, j int) bool {
	return c[i] < c[j]
}

// Swap swaps the two champion IDs.
func (c champions) Swap(i, j int) {
	temp := c[i]
	c[i] = c[j]
	c[j] = temp
}

// Available2 returns a slice containing all available pairs of champions
// within the pool, where ordering of the pair doesn't matter. The returned
// pair will be returned in ascending order of champion ID.
func (s State) Available2(pool ChampionSet) [][2]champion.Champion {
	var champs [][2]champion.Champion
	available := champions(s.Available1(pool))
	sort.Sort(available)
	for i, k1 := range available {
		for j := 0; j < i; j++ {
			champs = append(champs, [2]champion.Champion{available[j], k1})
		}
	}
	return champs
}

// BlueMinusRedUtility returns the utility difference of red and blue with
// respect to the given utility function.
func (s State) BlueMinusRedUtility(u Utility) (float64, error) {
	blue := s.Blue()
	red := s.Red()
	blueUtil, err := u(blue)
	if err != nil {
		return 0, err
	}
	redUtil, err := u(red)
	if err != nil {
		return 0, err
	}
	return blueUtil - redUtil, nil
}

// Payoff represents the value of an action from blue's perspective.
type Payoff struct {
	// Utility is blue side utility minus red side utility.
	Utility float64

	// NextState is the optimal next state given the current state.
	NextState State
}

// NewSolver initializes a solver for the given pool of champions and a utility
// function for evaluating a champion selection.
func NewSolver(champs ChampionSet, util Utility) *Solver {
	return &Solver{
		championPool: champs,
		utility:      util,
		cache:        make(map[State]Payoff),
	}
}

// Solver solves for bans and champion selections. It is illegal to use this
// directly. Use the NewSolver constructor to return a valid instance.
type Solver struct {
	// championPool is a set of all viable champions.
	championPool ChampionSet

	// cache maps from a given state to blue payoff for the previous action.
	// For example, if Actions is empty, then the value is the ex-ante blue
	// payoff with optimal selection from both sides. If Actions is a single
	// value, then that last action was a blue ban, so the payoff is from blue's
	// perspective with optimal selection from that point onward.
	cache map[State]Payoff

	// utility is the shared utility function used to evaluate a set of acquired
	// champions.
	utility Utility
}

// blueTwoActions solves for the payoff for two consecutive actions on blue
// side given a current state.
func (s *Solver) blueTwoActions(state State, action int, nextAction func(State) (Payoff, error)) (Payoff, error) {
	return s.twoActions(state, action, nextAction, func(a, b float64) bool { return a > b })
}

// redTwoActions solves for the payoff for two consecutive actions on red side
// given a current state.
func (s *Solver) redTwoActions(state State, action int, nextAction func(State) (Payoff, error)) (Payoff, error) {
	return s.twoActions(state, action, nextAction, func(a, b float64) bool { return a < b })
}

// twoActions solves for the payoff of two consecutive actions given a current
// state. greater takes two utility values from blue perspective and returns
// true if the first utility is greater than the second utility from the
// perspective of the player taking an action.
func (s *Solver) twoActions(state State, action int, nextAction func(State) (Payoff, error), greater func(float64, float64) bool) (Payoff, error) {
	fmt.Println("Solving twoActions on state:", state)
	ns := state.Normalize()
	cached, ok := s.cache[ns]
	if ok {
		cached.NextState = state.Merge(cached.NextState)
		return cached, nil
	}
	var (
		pay    Payoff
		looped bool
	)
	for _, champs := range state.Available2(s.championPool) {
		nextState := ns
		nextState.Actions[action] = champs[0]
		nextState.Actions[action+1] = champs[1]
		var (
			nextPay Payoff
			err     error
		)
		if nextAction != nil {
			nextPay, err = nextAction(nextState)
		} else {
			util, err := nextState.BlueMinusRedUtility(s.utility)
			if err != nil {
				return pay, err
			}
			nextPay = Payoff{
				Utility:   util,
				NextState: nextState,
			}
		}
		if err != nil {
			return pay, err
		}
		if !looped || greater(nextPay.Utility, pay.Utility) {
			looped = true
			pay.Utility = nextPay.Utility
			pay.NextState = nextPay.NextState
		}
	}
	if !looped {
		return pay, errors.New("no champions could be selected")
	}
	s.cache[ns] = pay
	// Overwrite the early history, which was normalized.
	pay.NextState = state.Merge(pay.NextState)
	return pay, nil

}

// oneAction is the same as twoActions, except for a single action.
func (s *Solver) oneAction(state State, action int, nextAction func(State) (Payoff, error), greater func(float64, float64) bool) (Payoff, error) {
	fmt.Println("Solving oneAction on state:", state)
	ns := state.Normalize()
	cached, ok := s.cache[ns]
	if ok {
		cached.NextState = state.Merge(cached.NextState)
		return cached, nil
	}
	var (
		pay    Payoff
		looped bool
	)
	for _, champ := range state.Available1(s.championPool) {
		nextState := ns
		nextState.Actions[action] = champ
		var (
			nextPay Payoff
			err     error
		)
		if nextAction != nil {
			nextPay, err = nextAction(nextState)
		} else {
			util, err := nextState.BlueMinusRedUtility(s.utility)
			if err != nil {
				return pay, err
			}
			nextPay = Payoff{
				Utility:   util,
				NextState: nextState,
			}
		}
		if err != nil {
			return pay, err
		}

		if !looped || greater(nextPay.Utility, pay.Utility) {
			looped = true
			pay.Utility = nextPay.Utility
			pay.NextState = nextPay.NextState
		}
	}
	if !looped {
		return pay, errors.New("unable to select a champion")
	}
	s.cache[ns] = pay
	// Overwrite the early history, which was normalized.
	pay.NextState = state.Merge(pay.NextState)
	return pay, nil
}

// blueOneAction is the same as blueTwoActions, except for a single action.
func (s *Solver) blueOneAction(state State, action int, nextAction func(State) (Payoff, error)) (Payoff, error) {
	return s.oneAction(state, action, nextAction, func(x, y float64) bool { return x > y })
}

// redOneAction is the same as redTwoActions, except for a single action.
func (s *Solver) redOneAction(state State, action int, nextAction func(State) (Payoff, error)) (Payoff, error) {
	return s.oneAction(state, action, nextAction, func(x, y float64) bool { return x < y })
}

// RedFifthPick returns the optimal state after red's fifth pick, given the
// input initial state.
func (s *Solver) RedFifthPick(state State) (Payoff, error) {
	return s.redOneAction(state, 19, nil)
}

// BlueFourthAndFifthPick returns the optimal state after blue's fourth and
// fifth pick, given the input initial state.
func (s *Solver) BlueFourthAndFifthPick(state State) (Payoff, error) {
	return s.blueTwoActions(state, 17, s.RedFifthPick)
}

// RedFourthPick returns the optimal state after red's fourth pick, given the
// input initial state.
func (s *Solver) RedFourthPick(state State) (Payoff, error) {
	return s.redOneAction(state, 16, s.BlueFourthAndFifthPick)
}

// BlueFifthBan returns the optimal state after blue's fifth ban, given the
// input initial state.
func (s *Solver) BlueFifthBan(state State) (Payoff, error) {
	return s.blueOneAction(state, 15, s.RedFourthPick)
}

// RedFifthBan returns the optimal state after red's fifth ban, given the
// input initial state.
func (s *Solver) RedFifthBan(state State) (Payoff, error) {
	return s.redOneAction(state, 14, s.BlueFifthBan)
}

// BlueFourthBan returns the optimal state after blue's fourth ban, given the
// input initial state.
func (s *Solver) BlueFourthBan(state State) (Payoff, error) {
	return s.blueOneAction(state, 13, s.RedFifthBan)
}

// RedFourthBan returns the optimal state after red's fourth ban, given the
// input initial state.
func (s *Solver) RedFourthBan(state State) (Payoff, error) {
	return s.redOneAction(state, 12, s.BlueFourthBan)
}

// RedThirdPick returns the optimal state after red's third pick, given the
// input initial state.
func (s *Solver) RedThirdPick(state State) (Payoff, error) {
	return s.redOneAction(state, 11, s.RedFourthBan)
}

// BlueSecondAndThirdPick returns the optimal state after blue's second and
// third pick, given the input initial state.
func (s *Solver) BlueSecondAndThirdPick(state State) (Payoff, error) {
	return s.blueTwoActions(state, 9, s.RedThirdPick)
}

// RedFirstAndSecondPick returns the optimal state after red's first and second
// pick, given the input initial state.
func (s *Solver) RedFirstAndSecondPick(state State) (Payoff, error) {
	return s.redTwoActions(state, 7, s.BlueSecondAndThirdPick)
}

// BlueFirstPick returns the optimal state after blue's first pick, given the
// input initial state.
func (s *Solver) BlueFirstPick(state State) (Payoff, error) {
	return s.blueOneAction(state, 6, s.RedFirstAndSecondPick)
}

// RedThirdBan returns the optimal state after red's third ban, given the input
// initial state.
func (s *Solver) RedThirdBan(state State) (Payoff, error) {
	return s.redOneAction(state, 5, s.BlueFirstPick)
}

// BlueThirdBan returns the optimal state after blue's third ban, given the
// input initial state.
func (s *Solver) BlueThirdBan(state State) (Payoff, error) {
	return s.blueOneAction(state, 4, s.RedThirdBan)
}

// RedSecondBan returns the optimal state after red's second ban, given the
// input initial state.
func (s *Solver) RedSecondBan(state State) (Payoff, error) {
	return s.redOneAction(state, 3, s.BlueThirdBan)
}

// BlueSecondBan returns the optimal state after blue's second ban, given the
// input initial state.
func (s *Solver) BlueSecondBan(state State) (Payoff, error) {
	return s.blueOneAction(state, 2, s.RedSecondBan)
}

// RedFirstBan returns the optimal state after red's first ban, given the input
// initial state.
func (s *Solver) RedFirstBan(state State) (Payoff, error) {
	return s.redOneAction(state, 1, s.BlueSecondBan)
}

// BlueFirstBan returns the optimal state after blue's first ban. There is no
// initial state, since this is the first move in the game.
func (s *Solver) BlueFirstBan() (Payoff, error) {
	return s.blueOneAction(State{}, 0, s.RedFirstBan)
}
