package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuhanfang/riot/constants/language"
	"github.com/yuhanfang/riot/staticdata"
)

func prettyPrint(res interface{}, err error) {
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	js, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}
	fmt.Println(string(js))
}

func main() {
	c := staticdata.New(http.DefaultClient)
	ctx := context.Background()

	versions, err := c.Versions(ctx)
	prettyPrint(versions, err)

	languages, err := c.Languages(ctx)
	prettyPrint(languages, err)

	items, err := c.Items(ctx, versions[0], languages[0])
	prettyPrint(items, err)

	champs, err := c.Champions(ctx, versions[0], language.EnglishUnitedStates)
	prettyPrint(champs, err)
}
