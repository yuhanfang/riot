// Package queue defines queue constants.
package queue

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Queue int

func (q *Queue) UnmarshalJSON(b []byte) error {
	var (
		s string
		i int
	)

	// First see if it is stored as native int.
	err := json.Unmarshal(b, &i)
	if err == nil {
		*q = Queue(i)
		return nil
	}

	// Must be a string.
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch strings.ToUpper(s) {
	case "RANKED_SOLO_5X5":
		*q = RankedSolo5x5
	case "RANKED_FLEX_SR":
		*q = RankedFlexSR
	case "RANKED_FLEX_TT":
		*q = RankedFlexTT
	default:
		return fmt.Errorf("invalid queue %q", s)
	}
	return nil
}

func (q Queue) String() string {
	switch q {
	case RankedSolo5x5:
		return "RANKED_SOLO_5x5"
	case RankedFlexSR:
		return "RANKED_FLEX_SR"
	case RankedFlexTT:
		return "RANKED_FLEX_TT"
	default:
		panic(fmt.Sprintf("invalid Queue %d", q))
	}
}

const (
	RankedSolo5x5 Queue = 420
	RankedFlexSR        = 440
	RankedFlexTT        = 470
)
