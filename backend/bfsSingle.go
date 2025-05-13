package main

import (
	"fmt"
	"strings"
	"time"
)

// BFS Single Path dengan queue yang benar
func bfsSinglePath(target string) ([]string, bool, time.Duration, int) {
	startTime := time.Now()
	target = strings.ToLower(target)

	if baseElements[target] {
		return []string{}, true, time.Since(startTime), 1
	}

	discovered := make(map[string]bool)
	recipeUsed := make(map[string][]string)
	queue := []string{}
	nodesVisited := 0

	// Inisialisasi queue dengan base elements
	for base := range baseElements {
		discovered[base] = true
		queue = append(queue, base)
		nodesVisited++
	}

	// Proses BFS level demi level
	for len(queue) > 0 {
		levelSize := len(queue)

		// Proses semua node di level saat ini
		for i := 0; i < levelSize; i++ {
			current := queue[0]
			queue = queue[1:]

			// Coba kombinasi dengan semua elemen yang sudah ditemukan
			for other := range discovered {
				// Cek kombinasi current + other
				if result, exists := getCombinationResult(current, other); exists {
					resultTier := tierMap[result]
					if tierMap[current] >= resultTier || tierMap[other] >= resultTier {
						continue
					}
					if !discovered[result] {
						discovered[result] = true
						recipeUsed[result] = []string{current, other}
						queue = append(queue, result)
						nodesVisited++

						if result == target {
							return reconstructPath(target, recipeUsed), true,
								time.Since(startTime), nodesVisited
						}
					}
				}
			}
		}
	}

	return nil, false, time.Since(startTime), nodesVisited
}

// Helper function untuk kombinasi elemen
func getCombinationResult(a, b string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	for result, recipes := range recipesMap {
		for _, ingredients := range recipes {
			if (ingredients[0] == a && ingredients[1] == b) ||
				(ingredients[0] == b && ingredients[1] == a) {
				return result, true
			}
		}
	}
	return "", false
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
