package staticdata

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/language"
)

type ItemList struct {
	Data    map[string]Item
	Version Version
	Tree    []ItemTree
	Groups  []ItemGroup
	Type    string
}

type ItemTree struct {
	Header string
	Tags   []string
}

type Item struct {
	Gold                 ItemGold
	Plaintext            string
	HideFromAll          bool
	InStore              bool
	Into                 []string
	Id                   int
	Stats                InventoryDataStats
	Colloq               string
	Maps                 map[string]bool
	SpecialRecipe        int
	Image                Image
	Description          string
	Tags                 []string
	Effect               map[string]string
	RequiredChampion     string
	RequiredAlly         string
	From                 []string
	Group                string
	ConsumeOnFull        bool
	Name                 string
	Consumed             bool
	SanitizedDescription string
	Depth                int
	Stacks               int
}

type ItemGold struct {
	Sell        int
	Total       int
	Base        int
	Purchasable bool
}

type InventoryDataStats struct {
	PercentCritDamageMod     float64
	PercentSpellBlockMod     float64
	PercentHPRegenMod        float64
	PercentMovementSpeedMod  float64
	FlatSpellBlockMod        float64
	FlatCritDamageMod        float64
	FlatEnergyPoolMod        float64
	PercentLifeStealMod      float64
	FlatMPPoolMod            float64
	FlatMovementSpeedMod     float64
	PercentAttackSpeedMod    float64
	FlatBlockMod             float64
	PercentBlockMod          float64
	FlatEnergyRegenMod       float64
	PercentSpellVampMod      float64
	FlatMPRegenMod           float64
	PercentDodgeMod          float64
	FlatAttackSpeedMod       float64
	FlatArmorMod             float64
	FlatHPRegenMod           float64
	PercentMagicDamageMod    float64
	PercentMPPoolMod         float64
	FlatMagicDamageMod       float64
	PercentMPRegenMod        float64
	PercentPhysicalDamageMod float64
	FlatPhysicalDamageMod    float64
	PercentHPPoolMod         float64
	PercentArmorMod          float64
	PercentCritChanceMod     float64
	PercentEXPBonus          float64
	FlatHPPoolMod            float64
	FlatCritChanceMod        float64
	FlatEXPBonus             float64
}

type Image struct {
	Full   string
	Group  string
	Sprite string
	H      int
	W      int
	Y      int
	X      int
}

type ItemGroup struct {
	MaxGroupOwnable string
	Key             string
}

// Items returns all items for the given game version and language.
func (c *Client) Items(ctx context.Context, v Version, lang language.Language) (*ItemList, error) {
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/item.json", v, lang)
	var itemList ItemList
	err := c.getJSON(ctx, url, &itemList)
	return &itemList, err
}
