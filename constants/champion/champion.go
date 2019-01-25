// Package champion defines champion constants.
package champion

import (
	"fmt"
	"sort"
)

type Champion int

const (
	Empty Champion = 0

	Annie        = 1
	Olaf         = 2
	Galio        = 3
	TwistedFate  = 4
	XinZhao      = 5
	Urgot        = 6
	Leblanc      = 7
	Vladimir     = 8
	Fiddlesticks = 9
	Kayle        = 10
	MasterYi     = 11
	Alistar      = 12
	Ryze         = 13
	Sion         = 14
	Sivir        = 15
	Soraka       = 16
	Teemo        = 17
	Tristana     = 18
	Warwick      = 19
	Nunu         = 20
	MissFortune  = 21
	Ashe         = 22
	Tryndamere   = 23
	Jax          = 24
	Morgana      = 25
	Zilean       = 26
	Singed       = 27
	Evelynn      = 28
	Twitch       = 29
	Karthus      = 30
	Chogath      = 31
	Amumu        = 32
	Rammus       = 33
	Anivia       = 34
	Shaco        = 35
	DrMundo      = 36
	Sona         = 37
	Kassadin     = 38
	Irelia       = 39
	Janna        = 40
	Gangplank    = 41
	Corki        = 42
	Karma        = 43
	Taric        = 44
	Veigar       = 45
	Trundle      = 48
	Swain        = 50
	Caitlyn      = 51
	Blitzcrank   = 53
	Malphite     = 54
	Katarina     = 55
	Nocturne     = 56
	Maokai       = 57
	Renekton     = 58
	JarvanIV     = 59
	Elise        = 60
	Orianna      = 61
	MonkeyKing   = 62
	Brand        = 63
	LeeSin       = 64
	Vayne        = 67
	Rumble       = 68
	Cassiopeia   = 69
	Skarner      = 72
	Heimerdinger = 74
	Nasus        = 75
	Nidalee      = 76
	Udyr         = 77
	Poppy        = 78
	Gragas       = 79
	Pantheon     = 80
	Ezreal       = 81
	Mordekaiser  = 82
	Yorick       = 83
	Akali        = 84
	Kennen       = 85
	Garen        = 86
	Leona        = 89
	Malzahar     = 90
	Talon        = 91
	Riven        = 92
	KogMaw       = 96
	Shen         = 98
	Lux          = 99
	Xerath       = 101
	Shyvana      = 102
	Ahri         = 103
	Graves       = 104
	Fizz         = 105
	Volibear     = 106
	Rengar       = 107
	Varus        = 110
	Nautilus     = 111
	Viktor       = 112
	Sejuani      = 113
	Fiora        = 114
	Ziggs        = 115
	Lulu         = 117
	Draven       = 119
	Hecarim      = 120
	Khazix       = 121
	Darius       = 122
	Jayce        = 126
	Lissandra    = 127
	Diana        = 131
	Quinn        = 133
	Syndra       = 134
	AurelionSol  = 136
	Kayn         = 141
	Zoe          = 142
	Zyra         = 143
	Kaisa        = 145
	Gnar         = 150
	Zac          = 154
	Yasuo        = 157
	Velkoz       = 161
	Taliyah      = 163
	Camille      = 164
	Braum        = 201
	Jhin         = 202
	Kindred      = 203
	Jinx         = 222
	TahmKench    = 223
	Lucian       = 236
	Zed          = 238
	Kled         = 240
	Ekko         = 245
	Vi           = 254
	Aatrox       = 266
	Nami         = 267
	Azir         = 268
	Thresh       = 412
	Illaoi       = 420
	RekSai       = 421
	Ivern        = 427
	Kalista      = 429
	Bard         = 432
	Rakan        = 497
	Xayah        = 498
	Ornn         = 516
	Sylas	     = 517
	Neeko        = 518
	Pyke         = 555
)

func (c Champion) String() string {
	switch c {
	case Empty:
		return ""
	case Annie:
		return "Annie"
	case Olaf:
		return "Olaf"
	case Galio:
		return "Galio"
	case TwistedFate:
		return "TwistedFate"
	case XinZhao:
		return "XinZhao"
	case Urgot:
		return "Urgot"
	case Leblanc:
		return "Leblanc"
	case Vladimir:
		return "Vladimir"
	case Fiddlesticks:
		return "Fiddlesticks"
	case Kayle:
		return "Kayle"
	case MasterYi:
		return "MasterYi"
	case Alistar:
		return "Alistar"
	case Ryze:
		return "Ryze"
	case Sion:
		return "Sion"
	case Sivir:
		return "Sivir"
	case Soraka:
		return "Soraka"
	case Teemo:
		return "Teemo"
	case Tristana:
		return "Tristana"
	case Warwick:
		return "Warwick"
	case Nunu:
		return "Nunu"
	case MissFortune:
		return "MissFortune"
	case Ashe:
		return "Ashe"
	case Tryndamere:
		return "Tryndamere"
	case Jax:
		return "Jax"
	case Morgana:
		return "Morgana"
	case Zilean:
		return "Zilean"
	case Singed:
		return "Singed"
	case Evelynn:
		return "Evelynn"
	case Twitch:
		return "Twitch"
	case Karthus:
		return "Karthus"
	case Chogath:
		return "Chogath"
	case Amumu:
		return "Amumu"
	case Rammus:
		return "Rammus"
	case Anivia:
		return "Anivia"
	case Shaco:
		return "Shaco"
	case DrMundo:
		return "DrMundo"
	case Sona:
		return "Sona"
	case Kassadin:
		return "Kassadin"
	case Irelia:
		return "Irelia"
	case Janna:
		return "Janna"
	case Gangplank:
		return "Gangplank"
	case Corki:
		return "Corki"
	case Karma:
		return "Karma"
	case Taric:
		return "Taric"
	case Veigar:
		return "Veigar"
	case Trundle:
		return "Trundle"
	case Swain:
		return "Swain"
	case Caitlyn:
		return "Caitlyn"
	case Blitzcrank:
		return "Blitzcrank"
	case Malphite:
		return "Malphite"
	case Katarina:
		return "Katarina"
	case Nocturne:
		return "Nocturne"
	case Maokai:
		return "Maokai"
	case Renekton:
		return "Renekton"
	case JarvanIV:
		return "JarvanIV"
	case Elise:
		return "Elise"
	case Orianna:
		return "Orianna"
	case MonkeyKing:
		return "MonkeyKing"
	case Brand:
		return "Brand"
	case LeeSin:
		return "LeeSin"
	case Vayne:
		return "Vayne"
	case Rumble:
		return "Rumble"
	case Cassiopeia:
		return "Cassiopeia"
	case Skarner:
		return "Skarner"
	case Heimerdinger:
		return "Heimerdinger"
	case Nasus:
		return "Nasus"
	case Nidalee:
		return "Nidalee"
	case Udyr:
		return "Udyr"
	case Poppy:
		return "Poppy"
	case Gragas:
		return "Gragas"
	case Pantheon:
		return "Pantheon"
	case Ezreal:
		return "Ezreal"
	case Mordekaiser:
		return "Mordekaiser"
	case Yorick:
		return "Yorick"
	case Akali:
		return "Akali"
	case Kennen:
		return "Kennen"
	case Garen:
		return "Garen"
	case Leona:
		return "Leona"
	case Malzahar:
		return "Malzahar"
	case Talon:
		return "Talon"
	case Riven:
		return "Riven"
	case KogMaw:
		return "KogMaw"
	case Shen:
		return "Shen"
	case Lux:
		return "Lux"
	case Xerath:
		return "Xerath"
	case Shyvana:
		return "Shyvana"
	case Ahri:
		return "Ahri"
	case Graves:
		return "Graves"
	case Fizz:
		return "Fizz"
	case Volibear:
		return "Volibear"
	case Rengar:
		return "Rengar"
	case Varus:
		return "Varus"
	case Nautilus:
		return "Nautilus"
	case Viktor:
		return "Viktor"
	case Sejuani:
		return "Sejuani"
	case Fiora:
		return "Fiora"
	case Ziggs:
		return "Ziggs"
	case Lulu:
		return "Lulu"
	case Draven:
		return "Draven"
	case Hecarim:
		return "Hecarim"
	case Khazix:
		return "Khazix"
	case Darius:
		return "Darius"
	case Jayce:
		return "Jayce"
	case Lissandra:
		return "Lissandra"
	case Diana:
		return "Diana"
	case Quinn:
		return "Quinn"
	case Syndra:
		return "Syndra"
	case AurelionSol:
		return "AurelionSol"
	case Kayn:
		return "Kayn"
	case Zoe:
		return "Zoe"
	case Zyra:
		return "Zyra"
	case Kaisa:
		return "Kaisa"
	case Gnar:
		return "Gnar"
	case Zac:
		return "Zac"
	case Yasuo:
		return "Yasuo"
	case Velkoz:
		return "Velkoz"
	case Taliyah:
		return "Taliyah"
	case Camille:
		return "Camille"
	case Braum:
		return "Braum"
	case Jhin:
		return "Jhin"
	case Kindred:
		return "Kindred"
	case Jinx:
		return "Jinx"
	case TahmKench:
		return "TahmKench"
	case Lucian:
		return "Lucian"
	case Zed:
		return "Zed"
	case Kled:
		return "Kled"
	case Ekko:
		return "Ekko"
	case Vi:
		return "Vi"
	case Aatrox:
		return "Aatrox"
	case Nami:
		return "Nami"
	case Azir:
		return "Azir"
	case Thresh:
		return "Thresh"
	case Illaoi:
		return "Illaoi"
	case RekSai:
		return "RekSai"
	case Ivern:
		return "Ivern"
	case Kalista:
		return "Kalista"
	case Bard:
		return "Bard"
	case Rakan:
		return "Rakan"
	case Xayah:
		return "Xayah"
	case Ornn:
		return "Ornn"
	case Sylas:
		return "Sylas"
	case Neeko:
		return "Neeko"
	case Pyke:
		return "Pyke"
	default:
		panic(fmt.Sprintf("invalid champion ID %d", c))
	}
}

type champions []Champion

func (c champions) Len() int {
	return len(c)
}
func (c champions) Swap(i, j int) {
	temp := c[i]
	c[i] = c[j]
	c[j] = temp
}
func (c champions) Less(i, j int) bool {
	return int64(c[i]) < int64(c[j])
}

// All returns a list of all champions in numerical ID order.
func All() []Champion {
	c := []Champion{
		Annie,
		Olaf,
		Galio,
		TwistedFate,
		XinZhao,
		Urgot,
		Leblanc,
		Vladimir,
		Fiddlesticks,
		Kayle,
		MasterYi,
		Alistar,
		Ryze,
		Sion,
		Sivir,
		Soraka,
		Teemo,
		Tristana,
		Warwick,
		Nunu,
		MissFortune,
		Ashe,
		Tryndamere,
		Jax,
		Morgana,
		Zilean,
		Singed,
		Evelynn,
		Twitch,
		Karthus,
		Chogath,
		Amumu,
		Rammus,
		Anivia,
		Shaco,
		DrMundo,
		Sona,
		Kassadin,
		Irelia,
		Janna,
		Gangplank,
		Corki,
		Karma,
		Taric,
		Veigar,
		Trundle,
		Swain,
		Caitlyn,
		Blitzcrank,
		Malphite,
		Katarina,
		Nocturne,
		Maokai,
		Renekton,
		JarvanIV,
		Elise,
		Orianna,
		MonkeyKing,
		Brand,
		LeeSin,
		Vayne,
		Rumble,
		Cassiopeia,
		Skarner,
		Heimerdinger,
		Nasus,
		Nidalee,
		Udyr,
		Poppy,
		Gragas,
		Pantheon,
		Ezreal,
		Mordekaiser,
		Yorick,
		Akali,
		Kennen,
		Garen,
		Leona,
		Malzahar,
		Talon,
		Riven,
		KogMaw,
		Shen,
		Lux,
		Xerath,
		Shyvana,
		Ahri,
		Graves,
		Fizz,
		Volibear,
		Rengar,
		Varus,
		Nautilus,
		Viktor,
		Sejuani,
		Fiora,
		Ziggs,
		Lulu,
		Draven,
		Hecarim,
		Khazix,
		Darius,
		Jayce,
		Lissandra,
		Diana,
		Quinn,
		Syndra,
		AurelionSol,
		Kayn,
		Zoe,
		Zyra,
		Kaisa,
		Gnar,
		Zac,
		Yasuo,
		Velkoz,
		Taliyah,
		Camille,
		Braum,
		Jhin,
		Kindred,
		Jinx,
		TahmKench,
		Lucian,
		Zed,
		Kled,
		Ekko,
		Vi,
		Aatrox,
		Nami,
		Azir,
		Thresh,
		Illaoi,
		RekSai,
		Ivern,
		Kalista,
		Bard,
		Rakan,
		Xayah,
		Ornn,
		Sylas,
		Neeko,
		Pyke,
	}
	sort.Sort(champions(c))
	return c
}
