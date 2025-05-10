package main

import (
	"fmt"
	"strings"
)

func bfsSinglePath(target string) ([]string, bool) {
	// Make target lowercase for consistent matching
	target = strings.ToLower(target)
	
	// Track discovered elements and their recipes
	discovered := make(map[string]bool)
	recipeUsed := make(map[string][]string)
	
	// Initialize with base elements
	for base := range baseElements {
		discovered[base] = true
	}
	
	// BFS implementation
	for {
		newDiscoveries := false
		
		// Try to create all possible new elements from currently discovered elements
		for resultElement, recipes := range recipesMap {
			// Skip if already discovered
			if discovered[resultElement] {
				continue
			}
			
			resultTier := tierMap[resultElement]
			
			// Try each recipe for this element
			for _, ingredients := range recipes {
				ingr1 := ingredients[0]
				ingr2 := ingredients[1]
				
				// Check if both ingredients are discovered
				if discovered[ingr1] && discovered[ingr2] {
					// Skip if ingredient tier >= result tier (prevent circular dependencies)
					if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
						if printCount < maxPrints {
							fmt.Printf("Skipping recipe due to tier: %s + %s => %s (tiers: %d,%d>=%d)\n", 
								ingr1, ingr2, resultElement, tierMap[ingr1], tierMap[ingr2], resultTier)
							printCount++
						}
						continue
					}
					
					if printCount < maxPrints {
						fmt.Printf("Discovered: %s + %s => %s\n", ingr1, ingr2, resultElement)
						printCount++
					}
					
					// Mark as discovered and record recipe used
					discovered[resultElement] = true
					recipeUsed[resultElement] = []string{ingr1, ingr2}
					newDiscoveries = true
					
					// Check if we found our target
					if resultElement == target {
						return reconstructPath(target, recipeUsed), true
					}
					
					// No need to check other recipes for this element
					break
				}
			}
		}
		
		// If no new elements were discovered in this iteration, we're done
		if !newDiscoveries {
			break
		}
	}
	
	return nil, false
}

func reconstructPath(target string, recipeUsed map[string][]string) []string {
	var steps []string
	
	// Helper function for recursively building the path
	var buildPath func(element string) []string
	buildPath = func(element string) []string {
		// Base case: this is a base element
		if baseElements[element] {
			return []string{}
		}
		
		ingredients, exists := recipeUsed[element]
		if !exists {
			return []string{}
		}
		
		// Get paths for both ingredients
		path1 := buildPath(ingredients[0])
		path2 := buildPath(ingredients[1])
		
		// Combine paths and add current step
		result := append(path1, path2...)
		result = append(result, fmt.Sprintf("%s + %s = %s", ingredients[0], ingredients[1], element))
		
		return result
	}
	
	steps = buildPath(target)
	return steps
}