// dfs_bidirectional.go
package main

import (
	"fmt"
	"strings"
	"time"
)

func dfsBidirectionalPath(target string) ([]string, bool, time.Duration, int) {
	target = strings.ToLower(target)
	startTime := time.Now()

	fmt.Printf("[DEBUG] Starting bidirectional DFS for target: %s\n", target)

	if baseElements[target] {
		fmt.Printf("[DEBUG] Target '%s' is a base element, no path needed\n", target)
		return []string{}, true, time.Since(startTime), 1
	}

	// Initialize visited sets and stacks
	visitedFromStart := map[string][]string{target: {}}
	visitedFromGoal := map[string][]string{}

	stackStart := []string{target}
	stackGoal := []string{}
	
	// Seed goal stack with base elements
	for base := range baseElements {
		visitedFromGoal[base] = []string{base}
		stackGoal = append(stackGoal, base)
	}

	nodesVisited := 0

	for len(stackStart) > 0 && len(stackGoal) > 0 {
		// Expand from start
		currentStart := stackStart[len(stackStart)-1]
		stackStart = stackStart[:len(stackStart)-1]
		nodesVisited++
		fmt.Printf("[DEBUG] Expanding from start: %s\n", currentStart)

		if _, ok := visitedFromGoal[currentStart]; ok {
			// Path found
			startPath := visitedFromStart[currentStart]
			goalPath := visitedFromGoal[currentStart]
			fullPath := append([]string{}, startPath...)
			for i := len(goalPath) - 1; i > 0; i-- {
				fullPath = append(fullPath, goalPath[i])
			}
			fmt.Printf("[DEBUG] Path found: %s\n", strings.Join(fullPath, " -> "))
			return fullPath, true, time.Since(startTime), nodesVisited
		}

		if recipes, ok := recipesMap[currentStart]; ok {
			for _, ingr := range recipes {
				ingrTier1 := tierMap[ingr[0]]
				ingrTier2 := tierMap[ingr[1]]
				elementTier := tierMap[currentStart]
				if ingrTier1 >= elementTier || ingrTier2 >= elementTier {
					fmt.Printf("[DEBUG] Skipping recipe due to tier: %s + %s = %s\n", ingr[0], ingr[1], currentStart)
					continue
				}
				for _, ing := range ingr {
					if _, seen := visitedFromStart[ing]; !seen {
						visitedFromStart[ing] = append(visitedFromStart[currentStart], fmt.Sprintf("%s + %s = %s", ingr[0], ingr[1], currentStart))
						stackStart = append(stackStart, ing)
						fmt.Printf("[DEBUG] Adding %s to start stack\n", ing)
					}
				}
			}
		}

		// Expand from goal
		currentGoal := stackGoal[len(stackGoal)-1]
		stackGoal = stackGoal[:len(stackGoal)-1]
		nodesVisited++
		fmt.Printf("[DEBUG] Expanding from goal: %s\n", currentGoal)

		if _, ok := visitedFromStart[currentGoal]; ok {
			// Path found
			startPath := visitedFromStart[currentGoal]
			goalPath := visitedFromGoal[currentGoal]
			fullPath := append([]string{}, startPath...)
			for i := len(goalPath) - 1; i > 0; i-- {
				fullPath = append(fullPath, goalPath[i])
			}
			fmt.Printf("[DEBUG] Path found: %s\n", strings.Join(fullPath, " -> "))
			return fullPath, true, time.Since(startTime), nodesVisited
		}

		if nextElements, ok := revGraph[currentGoal]; ok {
			for _, parent := range nextElements {
				if _, seen := visitedFromGoal[parent]; !seen {
					visitedFromGoal[parent] = append(visitedFromGoal[currentGoal], fmt.Sprintf("%s + %s = %s", recipesMap[parent][0][0], recipesMap[parent][0][1], parent))
					stackGoal = append(stackGoal, parent)
					fmt.Printf("[DEBUG] Adding %s to goal stack\n", parent)
				}
			}
		}
	}

	return nil, false, time.Since(startTime), nodesVisited
}
