// Package event defines event constants.
package event

type Event string

const (
	ChampionKill     Event = "CHAMPION_KILL"
	WardPlaced             = "WARD_PLACED"
	WardKill               = "WARD_KILL"
	BuildingKill           = "BUILDING_KILL"
	EliteMonsterKill       = "ELITE_MONSTER_KILL"
	ItemPurchased          = "ITEM_PURCHASED"
	ItemSold               = "ITEM_SOLD"
	ItemDestroyed          = "ITEM_DESTROYED"
	ItemUndo               = "ITEM_UNDO"
	SkillLevelUp           = "SKILL_LEVEL_UP"
	AscendedEvent          = "ASCENDED_EVENT"
	CapturePoint           = "CAPTURE_POINT"
	PoroKingSummon         = "PORO_KING_SUMMON"
)
