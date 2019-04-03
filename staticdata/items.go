package staticdata

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/language"
)

type ItemList struct {
	Data    map[string]Item `json:"data"`
	Version Version         `json:"version"`
	Tree    []ItemTree      `json:"tree"`
	Groups  []ItemGroup     `json:"groups"`
	Type    string          `json:"type"`
}

type ItemTree struct {
	Header string   `json:"header"`
	Tags   []string `json:"tags"`
}

type Item struct {
	Gold                 ItemGold           `json:"gold"`
	Plaintext            string             `json:"plaintext"`
	HideFromAll          bool               `json:"hideFromAll"`
	InStore              bool               `json:"inStore"`
	Into                 []string           `json:"into"`
	Id                   int                `json:"id"`
	Stats                InventoryDataStats `json:"stats"`
	Colloq               string             `json:"colloq"`
	Maps                 map[string]bool    `json:"maps"`
	SpecialRecipe        int                `json:"specialRecipe"`
	Image                Image              `json:"image"`
	Description          string             `json:"description"`
	Tags                 []string           `json:"tags"`
	Effect               map[string]string  `json:"effect"`
	RequiredChampion     string             `json:"requiredChampion"`
	RequiredAlly         string             `json:"requiredAlly"`
	From                 []string           `json:"from"`
	Group                string             `json:"group"`
	ConsumeOnFull        bool               `json:"consumeOnFull"`
	Name                 string             `json:"name"`
	Consumed             bool               `json:"consumed"`
	SanitizedDescription string             `json:"sanitizedDescription"`
	Depth                int                `json:"depth"`
	Stacks               int                `json:"stacks"`
}

type ItemGold struct {
	Sell        int  `json:"sell"`
	Total       int  `json:"total"`
	Base        int  `json:"base"`
	Purchasable bool `json:"purchasable"`
}

type InventoryDataStats struct {
	PercentCritDamageMod     float64 `json:"percentCritDamageMod"`
	PercentSpellBlockMod     float64 `json:"percentSpellBlockMod"`
	PercentHPRegenMod        float64 `json:"percentHPRegenMod"`
	PercentMovementSpeedMod  float64 `json:"percentMovementSpeedMod"`
	FlatSpellBlockMod        float64 `json:"flatSpellBlockMod"`
	FlatCritDamageMod        float64 `json:"flatCritDamageMod"`
	FlatEnergyPoolMod        float64 `json:"flatEnergyPoolMod"`
	PercentLifeStealMod      float64 `json:"percentLifeStealMod"`
	FlatMPPoolMod            float64 `json:"flatMPPoolMod"`
	FlatMovementSpeedMod     float64 `json:"flatMovementSpeedMod"`
	PercentAttackSpeedMod    float64 `json:"percentAttackSpeedMod"`
	FlatBlockMod             float64 `json:"flatBlockMod"`
	PercentBlockMod          float64 `json:"percentBlockMod"`
	FlatEnergyRegenMod       float64 `json:"flatEnergyRegenMod"`
	PercentSpellVampMod      float64 `json:"percentSpellVampMod"`
	FlatMPRegenMod           float64 `json:"flatMPRegenMod"`
	PercentDodgeMod          float64 `json:"percentDodgeMod"`
	FlatAttackSpeedMod       float64 `json:"flatAttackSpeedMod"`
	FlatArmorMod             float64 `json:"flatArmorMod"`
	FlatHPRegenMod           float64 `json:"flatHPRegenMod"`
	PercentMagicDamageMod    float64 `json:"percentMagicDamageMod"`
	PercentMPPoolMod         float64 `json:"percentMPPoolMod"`
	FlatMagicDamageMod       float64 `json:"flatMagicDamageMod"`
	PercentMPRegenMod        float64 `json:"percentMPRegenMod"`
	PercentPhysicalDamageMod float64 `json:"percentPhysicalDamageMod"`
	FlatPhysicalDamageMod    float64 `json:"flatPhysicalDamageMod"`
	PercentHPPoolMod         float64 `json:"percentHPPoolMod"`
	PercentArmorMod          float64 `json:"percentArmorMod"`
	PercentCritChanceMod     float64 `json:"percentCritChanceMod"`
	PercentEXPBonus          float64 `json:"percentEXPBonus"`
	FlatHPPoolMod            float64 `json:"flatHPPoolMod"`
	FlatCritChanceMod        float64 `json:"flatCritChanceMod"`
	FlatEXPBonus             float64 `json:"flatEXPBonus"`
}

type Image struct {
	Full   string `json:"full"`
	Group  string `json:"group"`
	Sprite string `json:"sprite"`
	H      int    `json:"h"`
	W      int    `json:"w"`
	Y      int    `json:"y"`
	X      int    `json:"x"`
}

type ItemGroup struct {
	MaxGroupOwnable string `json:"maxGroupOwnable"`
	Key             string `json:"key"`
}

// Items returns all items for the given game version and language.
func (c *Client) Items(ctx context.Context, v Version, lang language.Language) (*ItemList, error) {
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/item.json", v, lang)
	var itemList ItemList
	err := c.getJSON(ctx, url, &itemList)
	return &itemList, err
}
