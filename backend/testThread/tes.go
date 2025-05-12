package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"sort"
// )

// type Recipe struct {
// 	Element     string `json:"Element"`
// 	Ingredient1 string `json:"Ingredient1"`
// 	Ingredient2 string `json:"Ingredient2"`
// }

// var RecipesMap = map[string][][]string{}
// var TierMap = map[string]int{}

// func LoadRecipes(path string) ([]Recipe, error) {
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var recipes []Recipe
// 	err = json.Unmarshal(file, &recipes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return recipes, nil
// }

// func BuildRecipeMap(recipes []Recipe) {
// 	RecipesMap = map[string][][]string{}
// 	TierMap = map[string]int{}
// 	for _, r := range recipes {
// 		RecipesMap[r.Element] = append(RecipesMap[r.Element], []string{r.Ingredient1, r.Ingredient2})
// 		if _, ok := TierMap[r.Ingredient1]; !ok {
// 			TierMap[r.Ingredient1] = 0
// 		}
// 		if _, ok := TierMap[r.Ingredient2]; !ok {
// 			TierMap[r.Ingredient2] = 0
// 		}
// 		if _, ok := TierMap[r.Element]; !ok {
// 			TierMap[r.Element] = TierMap[r.Ingredient1] + TierMap[r.Ingredient2] + 1
// 		}
// 	}
// }

// func DebugPrintMaps() {
// 	fmt.Println("[Debug] RecipesMap Sample:")
// 	count := 0
// 	for k, v := range RecipesMap {
// 		fmt.Printf("  %s => %v\n", k, v)
// 		count++
// 		if count >= 10 {
// 			break
// 		}
// 	}
// 	fmt.Println("\n[Debug] TierMap Sample:")
// 	types := make([]string, 0, len(TierMap))
// 	for k := range TierMap {
// 		types = append(types, k)
// 	}
// 	sort.Strings(types)
// 	for i := 0; i < 10 && i < len(types); i++ {
// 		fmt.Printf("  %s = %d\n", types[i], TierMap[types[i]])
// 	}
// }

// func main() {
// 	recipes, err := LoadRecipes("recipes.json")
// 	if err != nil {
// 		fmt.Println("Failed to load recipes.json:", err)
// 		return
// 	}
// 	BuildRecipeMap(recipes)
// 	DebugPrintMaps()
// }
