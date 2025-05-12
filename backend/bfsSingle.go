package main

import (
	"fmt"
	"strings"
	"time"
)

func bfsSinglePath(target string) ([]string, bool, time.Duration, int) {
    startTime := time.Now()
    fmt.Printf("[DEBUG] Starting BFS for target: %s\n", target)
    target = strings.ToLower(target)
    
    // Check if target is already a base element
    if baseElements[target] {
        fmt.Printf("[DEBUG] Target '%s' is a base element, no path needed\n", target)
        return []string{}, true, time.Since(startTime), 1 // Base element = 1 node
    }
    
    // Maps to track discovered elements and the recipes used
    discovered := make(map[string]bool)
    recipeUsed := make(map[string][]string)
    
    // Start with base elements
    nodesVisited := 0 // Counter for nodes visited
    for base := range baseElements {
        discovered[base] = true
        nodesVisited++ // Count each base element as a visited node
        fmt.Printf("[DEBUG] Added base element: %s\n", base)
    }
    
    // BFS implementation
    iterationCount := 0
    for {
        iterationCount++
        fmt.Printf("[DEBUG] BFS iteration %d\n", iterationCount)
        newDiscoveries := false
        
        mutex.RLock() // Use read lock when accessing recipe maps
        for resultElement, recipes := range recipesMap {
            if discovered[resultElement] {
                continue // Skip already discovered elements
            }
            
            resultTier := tierMap[resultElement]
            
            for _, ingredients := range recipes {
                ingr1 := ingredients[0]
                ingr2 := ingredients[1]
                
                // Check if both ingredients are discovered
                if discovered[ingr1] && discovered[ingr2] {
                    // Check tier constraint
                    if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
                        continue
                    }
                    
                    discovered[resultElement] = true
                    nodesVisited++ // Count this newly discovered element
                    recipeUsed[resultElement] = []string{ingr1, ingr2}
                    newDiscoveries = true
                    
                    fmt.Printf("[DEBUG] Discovered new element: %s using %s + %s\n", 
                        resultElement, ingr1, ingr2)
                    
                    // Check if target is found
                    if resultElement == target {
                        mutex.RUnlock()
                        path := reconstructPath(target, recipeUsed)
                        fmt.Printf("[DEBUG] Target found! Path length: %d\n", len(path))
                        return path, true, time.Since(startTime), nodesVisited
                    }
                    break
                }
            }
        }
        mutex.RUnlock()
        
        // If no new discoveries were made, break the loop
        if !newDiscoveries {
            fmt.Printf("[DEBUG] No new discoveries in iteration %d, stopping BFS\n", iterationCount)
            break
        }
    }
    
    fmt.Printf("[DEBUG] Target '%s' not found after %d iterations\n", target, iterationCount)
    return nil, false, time.Since(startTime), nodesVisited
}

// reconstructPath builds the creation path from the target back to base elements
func reconstructPath(target string, recipeUsed map[string][]string) []string {
	fmt.Printf("[DEBUG] Reconstructing path for %s\n", target)
	var steps []string
	
	// Recursive function to build the path
	var buildPath func(element string) []string
	buildPath = func(element string) []string {
		if baseElements[element] {
			fmt.Printf("[DEBUG] Reached base element: %s\n", element)
			return []string{}
		}
		
		ingredients, exists := recipeUsed[element]
		if !exists {
			fmt.Printf("[DEBUG] Warning: No recipe found for %s\n", element)
			return []string{}
		}
		
		// Build paths for both ingredients recursively
		path1 := buildPath(ingredients[0])
		path2 := buildPath(ingredients[1])
		
		// Combine paths and add current step
		result := append(path1, path2...)
		step := fmt.Sprintf("%s + %s = %s", ingredients[0], ingredients[1], element)
		fmt.Printf("[DEBUG] Adding step: %s\n", step)
		result = append(result, step)
		
		return result
	}
	
	steps = buildPath(target)
	fmt.Printf("[DEBUG] Path reconstruction complete with %d steps\n", len(steps))
	return steps
}

