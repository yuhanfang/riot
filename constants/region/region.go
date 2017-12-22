// Package region defines region constants.
package region

import "fmt"

// Region represents a Riot server region. Only constants defined in this
// package are valid inputs for the client.
type Region string

const (
	// BR1 is Brazil.
	BR1 Region = "BR1"

	// EUN1 is Europe East.
	EUN1 = "EUN1"

	// EUW1 is Europe West.
	EUW1 = "EUW1"

	// JP1 is Japan.
	JP1 = "JP1"

	// KR is Korea.
	KR = "KR"

	// LA1 is Latin America North.
	LA1 = "LA1"

	// LA2 is Latin America South.
	LA2 = "LA2"

	// NA1 is North America.
	NA1 = "NA1"

	// OC1 is Oceania.
	OC1 = "OC1"

	// TR1 is Turkey.
	TR1 = "TR1"

	// RU is Russia.
	RU = "RU"
)

// All returns all supported regions.
func All() []Region {
	return []Region{
		BR1,
		EUN1,
		EUW1,
		JP1,
		KR,
		LA1,
		LA2,
		NA1,
		OC1,
		TR1,
		RU,
	}
}

// Host returns the full hostname corresponding to the region. This function
// panics if an invalid region is used.
func (r Region) Host() string {
	switch r {
	case BR1:
		return "https://br1.api.riotgames.com"
	case EUN1:
		return "https://eun1.api.riotgames.com"
	case EUW1:
		return "https://euw1.api.riotgames.com"
	case JP1:
		return "https://jp1.api.riotgames.com"
	case LA1:
		return "https://la1.api.riotgames.com"
	case LA2:
		return "https://la2.api.riotgames.com"
	case NA1:
		return "https://na1.api.riotgames.com"
	case OC1:
		return "https://oc1.api.riotgames.com"
	case TR1:
		return "https://tr1.api.riotgames.com"
	case KR:
		return "https://kr.api.riotgames.com"
	case RU:
		return "https://ru.api.riotgames.com"
	default:
		panic(fmt.Sprintf("region %s does not have a configured host", r))
	}
}
