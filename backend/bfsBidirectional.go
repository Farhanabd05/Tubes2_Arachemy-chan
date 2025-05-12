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
						fmt.Printf("[DEBUG] Forward: %s -> %s\n", strings.Join(forward[ingr1].Path, " -> "), step)
						break
					}
				}
			}
		}

		// Expand backward frontier
		for elem := range backward {
			for _, ingr := range recipesMap[elem] {
				for _, component := range ingr {
					if _, exists := backward[component]; !exists {
						step := fmt.Sprintf("%s + %s = %s", ingr[0], ingr[1], elem)
						path := append([]string{step}, backward[elem].Path...)
						newBackward[component] = NodeInfo{Path: path}
						fmt.Printf("[DEBUG] Backward: %s <- %s\n", step, strings.Join(backward[elem].Path, " <- "))
					}
				}
			}
		}

		for k, v := range newForward {
			forward[k] = v
			nodesVisited++
			if _, ok := backward[k]; ok {
				fmt.Printf("[DEBUG] Meeting point at %s (discovered by forward)\n", k)
				meetingPoint = k
				goto reconstruct
			}
		}
		for k, v := range newBackward {
			backward[k] = v
			nodesVisited++
			if _, ok := forward[k]; ok {
				fmt.Printf("[DEBUG] Meeting point at %s (discovered by backward)\n", k)
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
	forwardPath := forward[meetingPoint].Path
	backwardPath := backward[meetingPoint].Path
	fullPath := append(forwardPath, backwardPath...)
	return fullPath, true, time.Since(startTime), nodesVisited
}
