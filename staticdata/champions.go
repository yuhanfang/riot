package staticdata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yuhanfang/riot/constants/language"
)

type ChampionList struct {
	Keys    map[string]string   `json:"keys"`
	Data    map[string]Champion `json:"data"`
	Version string              `json:"version"`
	Type    string              `json:"type"`
	Format  string              `json:"format"`
}

type Champion struct {
	Info        ChampionInfo          `json:"info"`
	Enemytips   []string              `json:"enemytips"`
	Stats       ChampionStats         `json:"stats"`
	Name        string                `json:"name"`
	Title       string                `json:"title"`
	Image       Image                 `json:"image"`
	Tags        []string              `json:"tags"`
	Partype     string                `json:"partype"`
	Skins       []ChampionSkin        `json:"skins"`
	Passive     ChampionPassive       `json:"passive"`
	Recommended []ChampionRecommended `json:"recommended"`
	Allytips    []string              `json:"allytips"`
	Key         string                `json:"key"`
	Lore        string                `json:"lore"`
	Id          string                `json:"id"`
	Blurb       string                `json:"blurb"`
	Spells      []ChampionSpell       `json:"spells"`
}

type ChampionInfo struct {
	Difficulty int `json:"difficulty"`
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
}

type ChampionStats struct {
	Armorperlevel        float64 `json:"armorperlevel"`
	Hpperlevel           float64 `json:"hpperlevel"`
	Attackdamage         float64 `json:"attackdamage"`
	Mpperlevel           float64 `json:"mpperlevel"`
	Attackspeedoffset    float64 `json:"attackspeedoffset"`
	Armor                float64 `json:"armor"`
	Hp                   float64 `json:"hp"`
	Hpregenperlevel      float64 `json:"hpregenperlevel"`
	Spellblock           float64 `json:"spellblock"`
	Attackrange          float64 `json:"attackrange"`
	Movespeed            float64 `json:"movespeed"`
	Attackdamageperlevel float64 `json:"attackdamageperlevel"`
	Mpregenperlevel      float64 `json:"mpregenperlevel"`
	Mp                   float64 `json:"mp"`
	Spellblockperlevel   float64 `json:"spellblockperlevel"`
	Crit                 float64 `json:"crit"`
	Mpregen              float64 `json:"mpregen"`
	Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	Hpregen              float64 `json:"hpregen"`
	Critperlevel         float64 `json:"critperlevel"`
}

type ChampionSkin struct {
	Num  int    `json:"num"`
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ChampionPassive struct {
	Image                Image  `json:"image"`
	SanitizedDescription string `json:"sanitizedDescription"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
}

type ChampionRecommended struct {
	Map      string          `json:"map"`
	Blocks   []ChampionBlock `json:"blocks"`
	Champion string          `json:"champion"`
	Title    string          `json:"title"`
	Priority bool            `json:"priority"`
	Mode     string          `json:"mode"`
	Type     string          `json:"type"`
}

type ChampionBlock struct {
	Items   []ChampionBlockItem `json:"items"`
	RecMath bool                `json:"recMath"`
	Type    string              `json:"type"`
}

type ChampionBlockItem struct {
	Count int    `json:"count"`
	Id    string `json:"id"`
}

// RangeOrSelf represents either an ability that targets self, or an ability
// with a given range. If Self is true, then Range should be ignored.
// Otherwise, Range holds the ability range.
type RangeOrSelf struct {
	Self  bool  `json:"self"`
	Range []int `json:"range"`
}

func (r *RangeOrSelf) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err == nil {
		if s == "self" {
			r.Self = true
			return nil
		}
		return fmt.Errorf("RangeOrSelf has string value %q", s)
	}
	var ints []int
	err = json.Unmarshal(b, &ints)
	if err == nil {
		r.Range = ints
		return nil
	}
	return fmt.Errorf("RangeOrSelf should be self or []int; got %s", string(b))
}

type ChampionSpell struct {
	CooldownBurn         string        `json:"cooldownBurn"`
	Resource             string        `json:"resource"`
	Leveltip             SpellLevelTip `json:"leveltip"`
	Vars                 []SpellVars   `json:"vars"`
	CostType             string        `json:"costType"`
	Image                Image         `json:"image"`
	SanitizedDescription string        `json:"sanitizedDescription"`
	SanitizedTooltip     string        `json:"sanitizedTooltip"`
	Effect               []SpellEffect `json:"effect"`
	Tooltip              string        `json:"tooltip"`
	Maxrank              int           `json:"maxrank"`
	CostBurn             string        `json:"costBurn"`
	RangeBurn            string        `json:"rangeBurn"`
	Range                RangeOrSelf   `json:"range"`
	Cooldown             []float64     `json:"cooldown"`
	Cost                 []int         `json:"cost"`
	Key                  string        `json:"key"`
	Description          string        `json:"description"`
	EffectBurn           []string      `json:"effectBurn"`
	Altimages            []Image       `json:"altimages"`
	Name                 string        `json:"name"`
}

type SpellEffect struct {
	Details []float64 `json:"details"`
}

func (s *SpellEffect) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.Details)
}

type SpellLevelTip struct {
	Effect []string `json:"effect"`
	Label  []string `json:"label"`
}

type SpellVars struct {
	RanksWith string      `json:"ranksWith"`
	Dyn       string      `json:"dyn"`
	Link      string      `json:"link"`
	Coeff     interface{} `json:"coeff"`
	Key       string      `json:"key"`
}

func (c *Client) Champions(ctx context.Context, v Version, lang language.Language) (*ChampionList, error) {
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion.json", v, lang)
	var champList ChampionList
	err := c.getJSON(ctx, url, &champList)
	return &champList, err
}

func (c *Client) Champion(ctx context.Context, v Version, lang language.Language, championName string) (*Champion, error) {
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json", v, lang, championName)
	var champList ChampionList
	err := c.getJSON(ctx, url, &champList)
	champ := champList.Data[championName]
	return &champ, err
}
