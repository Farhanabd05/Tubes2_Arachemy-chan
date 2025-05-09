// dfs_single_path.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	// Type        string `json:"Type"`
	Type        int    `json:"Type"` 
}

type Result struct {
	Found bool     `json:"found"`
	Steps []string `json:"steps"`
}

var (
	recipesMap    map[string][][]string
	tierMap       map[string]int
	baseElements  = map[string]bool{
		"Fire": true, "Water": true, "Earth": true, "Air": true, "Time": true,
	}
	printCount    = 0
	maxPrints     = 200
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
		ingr := []string{r.Ingredient1, r.Ingredient2}
		recipesMap[r.Element] = append(recipesMap[r.Element], ingr)

		// if strings.HasPrefix(r.Type, "Tier") {
		// 	numStr := strings.TrimPrefix(r.Type, "Tier")
		// 	if val, err := strconv.Atoi(numStr); err == nil {
		// 		tierMap[r.Element] = val
		// 	}
		// }

		tierMap[r.Element] = r.Type
	}
}

func dfsSinglePath(element string, visited map[string]bool, trace []string) ([]string, bool) {
	if printCount < maxPrints {
		fmt.Println("Processing:", strings.Join(trace, " -> "), "->", element)
		printCount++
	}
	if baseElements[element] {
		return []string{}, true
	}
	if visited[element] {
		if printCount < maxPrints {
			fmt.Println("Cycle detected at:", element)
			printCount++
		}
		return nil, false
	}
	visited[element] = true

	recipes, ok := recipesMap[element]
	if !ok {
		if printCount < maxPrints {
			fmt.Println("No recipe found for:", element)
			printCount++
		}
		return nil, false
	}

	elementTier := tierMap[element]

	for _, ingr := range recipes {
		// skip if ingredient tier >= element tier
		ingrTier1 := tierMap[ingr[0]]
		ingrTier2 := tierMap[ingr[1]]
		if ingrTier1 >= elementTier || ingrTier2 >= elementTier {
			if printCount < maxPrints {
				fmt.Printf("Skipping recipe due to tier: %s + %s => %s\n", ingr[0], ingr[1], element)
			}
			continue
		}
		if printCount < maxPrints {
			fmt.Printf("Trying: %s + %s => %s\n", ingr[0], ingr[1], element)
		}
		newTrace := append([]string{}, trace...)
		newTrace = append(newTrace, element)
		leftSteps, ok1 := dfsSinglePath(ingr[0], copyMap(visited), newTrace)
		if !ok1 {
			continue
		}
		rightSteps, ok2 := dfsSinglePath(ingr[1], copyMap(visited), newTrace)
		if !ok2 {
			continue
		}
		combined := append(leftSteps, rightSteps...)
		combined = append(combined, fmt.Sprintf("%s + %s => %s", ingr[0], ingr[1], element))
		return combined, true
	}

	return nil, false
}

func copyMap(m map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dfs_single_path.go <element>")
		return
	}
	target := os.Args[1]

	recipes, err := loadRecipes("recipes.json")
	if err != nil {
		fmt.Println("Error loading recipes:", err)
		return
	}
	buildRecipeMap(recipes)

	steps, ok := dfsSinglePath(target, map[string]bool{}, []string{})
	result := Result{
		Found: ok,
		Steps: steps,
	}
	jsonResult, _ := json.Marshal(result)
	fmt.Println(string(jsonResult))
}