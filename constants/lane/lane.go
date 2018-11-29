// Package lane defines lane constants.
package lane

type Lane string
type Type string

const (
	Middle Lane = "MID"
	Top         = "TOP"
	Jungle      = "JUNGLE"
	Bottom      = "BOTTOM"
)

const (
	TypeMiddle Type = "MID_LANE"
	TypeTop         = "TOP_LANE"
	TypeBottom      = "BOT_LANE"
)
