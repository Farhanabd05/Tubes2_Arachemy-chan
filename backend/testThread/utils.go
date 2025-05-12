package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Recipe represents a single recipe from the JSON file
type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Type        int    `json:"Type"`
}

// Job represents a single task to be processed by a worker
type Job struct {
	JobID   int
	JobType string
	Target  string
}

type Result struct {
	Found bool     `json:"found"`
	Steps []string `json:"steps"`
	Runtime time.Duration `json:"runtime"`
	NodesVisited int `json:"nodesVisited"`
}

// JobResult contains the result of a processed job
type JobResult struct {
	JobID    int
	WorkerID int
	JobType  string
	Target   string
	Found    bool
	Steps    []string
	Err      error
	Duration time.Duration // Added to track job execution time
}

// Global variables for recipe data
var (
	recipesMap   map[string][][]string
	tierMap      map[string]int
	revGraph     map[string][]string
	baseElements = map[string]bool{
		"fire": true, "water": true, "earth": true, "air": true, "time": true,
	}
	mutex sync.RWMutex // For thread-safe access to recipe maps
	printCount = 0
	maxPrints  = 200
)

// loadRecipes loads recipe data from a JSON file
func loadRecipes(file string) ([]Recipe, error) {
	fmt.Printf("[DEBUG] Loading recipes from %s\n", file)
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read file: %v\n", err)
		return nil, err
	}
	
	var recipes []Recipe
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		fmt.Printf("[ERROR] Failed to unmarshal JSON: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("[DEBUG] Successfully loaded %d recipes\n", len(recipes))
	return recipes, nil
}

// buildRecipeMap constructs the recipe and tier maps
func buildRecipeMap(recipes []Recipe) {
	fmt.Println("[DEBUG] Building recipe map")
	recipesMap = make(map[string][][]string)
	tierMap = make(map[string]int)
	
	for _, r := range recipes {
		element := strings.ToLower(r.Element)
		ingr1 := strings.ToLower(r.Ingredient1)
		ingr2 := strings.ToLower(r.Ingredient2)
		ingr := []string{ingr1, ingr2}
		
		recipesMap[element] = append(recipesMap[element], ingr)
		tierMap[element] = r.Type
	}
	
	fmt.Printf("[DEBUG] Recipe map built with %d elements\n", len(recipesMap))
}

// buildReverseGraph constructs a reverse lookup graph for efficient path finding
func buildReverseGraph() {
	fmt.Println("[DEBUG] Building reverse graph")
	revGraph = make(map[string][]string)
	
	for result, recipes := range recipesMap {
		for _, ingr := range recipes {
			// Add result to each ingredient's created elements list
			revGraph[ingr[0]] = appendUnique(revGraph[ingr[0]], result)
			revGraph[ingr[1]] = appendUnique(revGraph[ingr[1]], result)
		}
	}
	
	fmt.Printf("[DEBUG] Reverse graph built with %d elements\n", len(revGraph))
}

// appendUnique adds an element to a slice if it doesn't already exist
func appendUnique(slice []string, element string) []string {
	for _, e := range slice {
		if e == element {
			return slice
		}
	}
	return append(slice, element)
}