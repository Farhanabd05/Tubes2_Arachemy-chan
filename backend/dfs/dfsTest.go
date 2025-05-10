// dfs_single_path.go
package main

import (
	"fmt"
	"strings"
)

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
		if ingrTier1 > elementTier || ingrTier2 > elementTier {
			if printCount < maxPrints {
				fmt.Printf("Skipping recipe due to tier: %s + %s = %s\n", ingr[0], ingr[1], element)
			}
			continue
		}
		if printCount < maxPrints {
			fmt.Printf("Trying: %s + %s = %s\n", ingr[0], ingr[1], element)
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
		combined = append(combined, fmt.Sprintf("%s + %s = %s", ingr[0], ingr[1], element))
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