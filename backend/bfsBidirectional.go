package main

import (
	"fmt"
	"strings"
	"time"
)

func bfsBidirectionalPath(target string) ([]string, bool, time.Duration, int) {
	startTime := time.Now()
	target = strings.ToLower(target)

	if baseElements[target] {
		return []string{}, true, time.Since(startTime), 1
	}

	type NodeInfo struct {
		Path []string
	}

	forward := make(map[string]NodeInfo)
	backward := make(map[string]NodeInfo)

	mutex.RLock()
	defer mutex.RUnlock()

	// Initialize forward (from base) and backward (from target)
	for base := range baseElements {
		forward[base] = NodeInfo{Path: []string{}}
	}
	backward[target] = NodeInfo{Path: []string{}}

	nodesVisited := len(forward) + 1

	meetingPoint := ""
	iteration := 0

	for {
		iteration++
		fmt.Printf("[DEBUG] Iteration %d\n", iteration)

		newForward := make(map[string]NodeInfo)
		newBackward := make(map[string]NodeInfo)

		// Expand forward frontier
		for result, recipes := range recipesMap {
			if _, seen := forward[result]; seen {
				continue
			}

			for _, ingr := range recipes {
				ingr1, ingr2 := ingr[0], ingr[1]

				if _, ok1 := forward[ingr1]; ok1 {
					if _, ok2 := forward[ingr2]; ok2 {
						if tierMap[ingr1] >= tierMap[result] || tierMap[ingr2] >= tierMap[result] {
							continue
						}
						path := append(append([]string{}, forward[ingr1].Path...), forward[ingr2].Path...)
						step := fmt.Sprintf("%s + %s = %s", ingr1, ingr2, result)
						path = append(path, step)
						newForward[result] = NodeInfo{Path: path}
						fmt.Printf("[DEBUG] Forward discovered: %s\n", result)
						break
					}
				}
			}
		}

		// Expand backward frontier
		for elem, _ := range backward {
			for _, ingr := range recipesMap[elem] {
				for _, component := range ingr {
					if _, exists := backward[component]; !exists {
						step := fmt.Sprintf("%s + %s = %s", ingr[0], ingr[1], elem)
						path := append([]string{step}, backward[elem].Path...)
						newBackward[component] = NodeInfo{Path: path}
						fmt.Printf("[DEBUG] Backward discovered: %s\n", component)
					}
				}
			}
		}

		for k, v := range newForward {
			forward[k] = v
			nodesVisited++
			if _, ok := backward[k]; ok {
				meetingPoint = k
				goto reconstruct
			}
		}
		for k, v := range newBackward {
			backward[k] = v
			nodesVisited++
			if _, ok := forward[k]; ok {
				meetingPoint = k
				goto reconstruct
			}
		}

		if len(newForward) == 0 && len(newBackward) == 0 {
			break
		}
	}

	// Not found
	fmt.Println("[DEBUG] No meeting point found")
	return nil, false, time.Since(startTime), nodesVisited

reconstruct:
	fmt.Printf("[DEBUG] Found meeting point: %s\n", meetingPoint)

	// Build a complete recipe map from our search
	recipeMap := make(map[string][]string)

	// Add recipes from both directions to our recipe map
	for _, recipes := range recipesMap {
		for _, ingr := range recipes {
			resultElement := ""
			for _, step := range append(forward[meetingPoint].Path, backward[meetingPoint].Path...) {
				parts := strings.Split(step, " = ")
				if len(parts) == 2 {
					resultElem := parts[1]
					ingredients := strings.Split(parts[0], " + ")
					if len(ingredients) == 2 {
						recipeMap[resultElem] = ingredients
					}
				}
			}
			if resultElement != "" {
				recipeMap[resultElement] = ingr
			}
		}
	}

	var buildFullPath func(element string, visited map[string]bool, depth int) []string
	buildFullPath = func(element string, visited map[string]bool, depth int) []string {
    
    if baseElements[element] {
        return []string{}
    }
    
    if visited[element] {
        return []string{}
    }
    visited[element] = true

    // Prioritize recipe dari hasil search
    ingredients, exists := recipeMap[element]
    if !exists {
        if recipes, ok := recipesMap[element]; ok && len(recipes) > 0 {
            // Ambil resep dengan tier terendah
            ingredients = findLowestTierRecipe(recipes, element)
        } else {
            return []string{}
        }
    }

    // Bangun path untuk kedua bahan dengan tracking depth
    path1 := buildFullPath(ingredients[0], copyMap(visited), depth+1)
    path2 := buildFullPath(ingredients[1], copyMap(visited), depth+1)

    // Gabungkan path dan tambahkan step
    combined := append(path1, path2...)
    step := fmt.Sprintf("%s + %s = %s", ingredients[0], ingredients[1], element)
    return append(combined, step)
	}

	// Di bagian pemanggilan:
	completePath := buildFullPath(target, make(map[string]bool), 0)
	return completePath, true, time.Since(startTime), nodesVisited
	}

// Tambahkan fungsi helper
func findLowestTierRecipe(recipes [][]string, element string) []string {
    minTier := int(^uint(0) >> 1)
    var bestRecipe []string
    
    for _, recipe := range recipes {
        tier1 := tierMap[recipe[0]]
        tier2 := tierMap[recipe[1]]
        if tier1 < minTier || tier2 < minTier {
            minTier = min(tier1, tier2)
            bestRecipe = recipe
        }
    }
    return bestRecipe
}

