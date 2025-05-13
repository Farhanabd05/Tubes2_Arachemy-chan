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
	meetingPoint := ""
	meetingPointPathStart := []string{}
	meetingPointPathGoal := []string{}
	
	for len(stackStart) > 0 && len(stackGoal) > 0 {
		// Expand from start
		currentStart := stackStart[len(stackStart)-1]
		stackStart = stackStart[:len(stackStart)-1]
		nodesVisited++
		fmt.Printf("[DEBUG] Expanding from start: %s\n", currentStart)
		
		if paths, ok := visitedFromGoal[currentStart]; ok {
			meetingPoint = currentStart
			meetingPointPathStart = visitedFromStart[currentStart]
			meetingPointPathGoal = paths
			break
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
		
		if paths, ok := visitedFromStart[currentGoal]; ok {
			meetingPoint = currentGoal
			meetingPointPathStart = paths
			meetingPointPathGoal = visitedFromGoal[currentGoal]
			break
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
	
	if meetingPoint == "" {
		return nil, false, time.Since(startTime), nodesVisited
	}
	
	fmt.Printf("[DEBUG] Path found via meeting point: %s\n", meetingPoint)
	
	// Build a recipe map from our search results
	recipeMap := make(map[string][]string)
	
	// Extract recipes from paths
	extractRecipes := func(paths []string) {
		for _, step := range paths {
			if strings.Contains(step, " = ") {
				parts := strings.Split(step, " = ")
				resultElem := parts[1]
				ingredients := strings.Split(parts[0], " + ")
				if len(ingredients) == 2 {
					recipeMap[resultElem] = ingredients
				}
			}
		}
	}
	
	// Extract recipes from both directions
	extractRecipes(meetingPointPathStart)
	extractRecipes(meetingPointPathGoal)
	
	// Recursively build the complete path
	var buildFullPath func(element string, visited map[string]bool) []string
	buildFullPath = func(element string, visited map[string]bool) []string {
		if baseElements[element] {
			return []string{}
		}
		
		if visited[element] {
			return []string{}  // Avoid cycles
		}
		visited[element] = true
		
		// Find recipe for this element
		ingredients, exists := recipeMap[element]
		if !exists {
			// If we don't have the recipe in our map, try to find one from all recipes
			if recipes, ok := recipesMap[element]; ok && len(recipes) > 0 {
				ingredients = recipes[0]
			} else {
				return []string{}
			}
		}
		
		// Recursively build paths for ingredients
		path1 := buildFullPath(ingredients[0], copyMap(visited))
		path2 := buildFullPath(ingredients[1], copyMap(visited))
		
		// Combine paths and add this step
		result := append(path1, path2...)
		step := fmt.Sprintf("%s + %s = %s", ingredients[0], ingredients[1], element)
		return append(result, step)
	}
	
	// Build the complete path from target to base elements
	completePath := buildFullPath(target, make(map[string]bool))
	
	return completePath, true, time.Since(startTime), nodesVisited
}
