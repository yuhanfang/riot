// Package region defines region constants.
package region

import "fmt"

// Region represents a Riot server region. Only constants defined in this
// package are valid inputs for the client.
type Region string

const (
	NA1   Region = "NA1"
	TRLH1        = "TRLH1"
)

// Host returns the full hostname corresponding to the region. This function
// panics if an invalid region is used.
func (r Region) Host() string {
	switch r {
	case NA1:
		return "https://na1.api.riotgames.com"
	default:
		panic(fmt.Sprintf("region %d does not have a configured host", r))
	}
}
