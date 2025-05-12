package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Type        int    `json:"Type"`
}

type RecipeStep struct {
	Element     string   `json:"element"`
	Ingredients []string `json:"ingredients"`
}

var (
	recipesMap   map[string][][]string
	tierMap      map[string]int
	baseElements = map[string]bool{
		"fire": true, "water": true, "earth": true, "air": true, "time": true,
	}
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
		tierMap[r.Element] = r.Type
		if _, ok := tierMap[r.Ingredient1]; !ok {
			tierMap[r.Ingredient1] = 0
		}
		if _, ok := tierMap[r.Ingredient2]; !ok {
			tierMap[r.Ingredient2] = 0
		}
	}
}

func cloneMap(m map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func dfsFullPath(
	element string,
	visited map[string]bool,
	maxCount int,
	nodeCount *int64,
	results chan<- []RecipeStep, // Correct channel use (send-only)
	path []RecipeStep,
	mu *sync.Mutex,
) {
	atomic.AddInt64(nodeCount, 1)

	// Base case for basic elements
	if baseElements[element] {
		results <- append([]RecipeStep{}, path...) // Send combined recipe path
		return
	}

	// Skip if the element has already been visited
	mu.Lock() // Lock to ensure exclusive access to visited map
	if visited[element] {
		mu.Unlock()
		return
	}
	visited[element] = true
	mu.Unlock() // Unlock after marking the element visited

	recipes := recipesMap[element]
	elementTier := tierMap[element]

	// Explore recipes for each ingredient
	for _, ingr := range recipes {
		ingr1, ingr2 := ingr[0], ingr[1]
		// Continue only if both ingredients have a valid tier
		if tierMap[ingr1] > elementTier || tierMap[ingr2] > elementTier {
			continue
		}

		// Channel and goroutine setup for parallel DFS
		leftChan := make(chan []RecipeStep, 8)
		rightChan := make(chan []RecipeStep, 8)
		var wg sync.WaitGroup
		wg.Add(2)

		// Parallel DFS for left and right ingredients
		go func() {
			defer wg.Done()
			dfsFullPath(ingr1, cloneMap(visited), maxCount, nodeCount, leftChan, append([]RecipeStep{}, path...), mu)
			close(leftChan)
		}()
		go func() {
			defer wg.Done()
			dfsFullPath(ingr2, cloneMap(visited), maxCount, nodeCount, rightChan, append([]RecipeStep{}, path...), mu)
			close(rightChan)
		}()

		// Wait for goroutines to finish
		wg.Wait()

		// Combine results from left and right channels
		for l := range leftChan {
			for r := range rightChan {
				// Stop if maxCount recipes have been found
				if len(results) >= maxCount {
					return
				}

				// Combine the paths and send them to the results channel
				combined := append([]RecipeStep{}, l...)
				combined = append(combined, r...)
				combined = append(combined, RecipeStep{
					Element:     element,
					Ingredients: []string{ingr1, ingr2},
				})

				// Send the combined result to the channel
				results <- combined
			}
		}
	}
}

// Function to find all DFS paths for a target element
func FindAllDFSPaths(target string, maxCount int) ([][]RecipeStep, int) {
	results := make(chan []RecipeStep, maxCount)
	var nodeCount int64
	var mu sync.Mutex // Mutex to synchronize visited map updates

	// Start DFS for the target element
	go func() {
		dfsFullPath(target, map[string]bool{}, maxCount, &nodeCount, results, []RecipeStep{}, &mu)
		close(results)
	}()

	// Collect results from the channel
	var all [][]RecipeStep
	for result := range results {
		if len(all) >= maxCount {
			break
		}
		all = append(all, result)
	}

	// Print the found recipes
	for i, path := range all {
		fmt.Printf("Recipe %d:\n", i+1)
		for _, step := range path {
			fmt.Printf("  %s = %s + %s\n", step.Element, step.Ingredients[0], step.Ingredients[1])
		}
		fmt.Println()
	}

	return all, int(nodeCount)
}

// Main function
func main() {
	recipes, err := loadRecipes("recipes.json")
	if err != nil {
		fmt.Println("Failed to load recipes:", err)
		return
	}
	buildRecipeMap(recipes)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run multithread_dfs.go <element>")
		return
	}
	target := os.Args[1]
	maxCount := 5

	// Find and print the recipes for the target element
	_, visited := FindAllDFSPaths(target, maxCount)
	fmt.Println("Nodes visited:", visited)
}
