package solver

import (
	"testing"

	"github.com/yuhanfang/riot/constants/champion"
)

var (
	minPool        = ChampionSetOf([]champion.Champion{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	minPoolPlusOne = ChampionSetOf([]champion.Champion{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21})
)

func idUtility(c ChampionSet) (float64, error) {
	var util float64
	for k := range c {
		util += float64(k)
	}
	return util, nil
}

func TestRedFifthPick(t *testing.T) {
	pool := minPoolPlusOne
	state := State{
		Actions: [20]champion.Champion{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 0},
	}
	s := NewSolver(pool, idUtility)
	p, err := s.RedFifthPick(state)
	if err != nil {
		t.Fatal(err)
	}
	if p.NextState.Actions[19] != 21 {
		t.Errorf("red last pick = %d; want %d", p.NextState.Actions[19], 21)
	}
}

func TestBlueFirstBan(t *testing.T) {
	pool := minPoolPlusOne
	s := NewSolver(pool, idUtility)
	p, err := s.BlueFirstBan()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}
