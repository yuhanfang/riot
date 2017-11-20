package staticdata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yuhanfang/riot/constants/language"
)

type ChampionList struct {
	Keys    map[string]string
	Data    map[string]Champion
	Version string
	Type    string
	Format  string
}

type Champion struct {
	Info        ChampionInfo
	Enemytips   []string
	Stats       ChampionStats
	Name        string
	Title       string
	Image       Image
	Tags        []string
	Partype     string
	Skins       []ChampionSkin
	Passive     ChampionPassive
	Recommended []ChampionRecommended
	Allytips    []string
	Key         string
	Lore        string
	Id          string
	Blurb       string
	Spells      []ChampionSpell
}

type ChampionInfo struct {
	Difficulty int
	Attack     int
	Defense    int
	Magic      int
}

type ChampionStats struct {
	Armorperlevel        float64
	Hpperlevel           float64
	Attackdamage         float64
	Mpperlevel           float64
	Attackspeedoffset    float64
	Armor                float64
	Hp                   float64
	Hpregenperlevel      float64
	Spellblock           float64
	Attackrange          float64
	Movespeed            float64
	Attackdamageperlevel float64
	Mpregenperlevel      float64
	Mp                   float64
	Spellblockperlevel   float64
	Crit                 float64
	Mpregen              float64
	Attackspeedperlevel  float64
	Hpregen              float64
	Critperlevel         float64
}

type ChampionSkin struct {
	Num  int
	Name string
	Id   int
}

type ChampionPassive struct {
	Image                Image
	SanitizedDescription string
	Name                 string
	Description          string
}

type ChampionRecommended struct {
	Map      string
	Blocks   []ChampionBlock
	Champion string
	Title    string
	Priority bool
	Mode     string
	Type     string
}

type ChampionBlock struct {
	Items   []ChampionBlockItem
	RecMath bool
	Type    string
}

type ChampionBlockItem struct {
	Count int
	Id    int
}

// RangeOrSelf represents either an ability that targets self, or an ability
// with a given range. If Self is true, then Range should be ignored.
// Otherwise, Range holds the ability range.
type RangeOrSelf struct {
	Self  bool
	Range []int
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
	CooldownBurn         string
	Resource             string
	Leveltip             SpellLevelTip
	Vars                 []SpellVars
	CostType             string
	Image                Image
	SanitizedDescription string
	SanitizedTooltip     string
	Effect               [][]float64
	Tooltip              string
	Maxrank              int
	CostBurn             string
	RangeBurn            string
	Range                RangeOrSelf
	Cooldown             []float64
	Cost                 []int
	Key                  string
	Description          string
	EffectBurn           []string
	Altimages            []Image
	Name                 string
}

type SpellLevelTip struct {
	Effect []string
	Label  []string
}

type SpellVars struct {
	RanksWith string
	Dyn       string
	Link      string
	Coeff     []float64
	Key       string
}

func (c *Client) Champions(ctx context.Context, v Version, lang language.Language) (*ChampionList, error) {
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion.json", v, lang)
	var champList ChampionList
	err := c.getJSON(ctx, url, &champList)
	return &champList, err
}
