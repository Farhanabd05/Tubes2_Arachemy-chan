package main

import (
	"encoding/json"
	"os"
	"strings"
)

type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	ImgUrl1     string `json:"ImgUrl1"`
	ImgUrl2     string `json:"ImgUrl2"`
	Type        int    `json:"Type"`
}

type Result struct {
	Found bool     `json:"found"`
	Steps []string `json:"steps"`
}

var (
	recipesMap   map[string][][]string
	tierMap      map[string]int
	baseElements = map[string]bool{
		"fire": true, "water": true, "earth": true, "air": true, "time": true,
	}
	printCount = 0
	maxPrints  = 200
)

func loadRecipes(file string) ([]Recipe, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var recipes []Recipe
	err = json.Unmarshal(data, &recipes)
	return recipes, err
}

func buildRecipeMap(recipes []Recipe) {
	recipesMap = make(map[string][][]string)
	tierMap = make(map[string]int)
	for _, r := range recipes {
		// Convert everything to lowercase for consistent matching
		element := strings.ToLower(r.Element)
		ingr1 := strings.ToLower(r.Ingredient1)
		ingr2 := strings.ToLower(r.Ingredient2)
		
		ingr := []string{ingr1, ingr2}
		recipesMap[element] = append(recipesMap[element], ingr)
		tierMap[element] = r.Type
	}
}