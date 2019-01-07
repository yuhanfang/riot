// Package tier defines tier constants.
package tier

import (
	"encoding/json"
	"strings"
)

type Tier string

func (t *Tier) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = Tier(strings.ToUpper(s))
	return nil
}

const (
	Challenger  Tier = "CHALLENGER"
	Grandmaster      = "GRANDMASTER"
	Master           = "MASTER"
	Diamond          = "DIAMOND"
	Platinum         = "PLATINUM"
	Gold             = "GOLD"
	Silver           = "SILVER"
	Bronze           = "BRONZE"
	Iron             = "IRON"
)
