// dfsMultiple2.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DFSResult struct {
	Found       bool     `json:"found"`
	Steps       []string `json:"steps"`
	WorkerID    int      `json:"worker_id"`
	NodesVisited int     `json:"nodes_visited"`
	Runtime     float64  `json:"runtime"`
}

var (
	globalWorkerID = 1
	resultsChan    = make(chan DFSResult, 1000)
	wg             sync.WaitGroup
	maxResults     int = 5
	resultCount    int = 0
)

func mapCopy(src map[string]bool) map[string]bool {
	dst := make(map[string]bool)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func spawnDFSForEachRecipe(target string, visited map[string]bool, path []string) {
	recipes := recipesMap[target]
	for _, ingr := range recipes {
		mutex.Lock()
		if resultCount >= maxResults {
			mutex.Unlock()
			break
		}
		workerID := globalWorkerID
		globalWorkerID++
		mutex.Unlock()
		
		i1, i2 := ingr[0], ingr[1]
		wg.Add(1)
		
		go func(i1, i2 string, workerID int) {
			defer wg.Done()
			startTime := time.Now()
			
			visited1 := mapCopy(visited)
			visited2 := mapCopy(visited)
			leftPaths := dfsCombinatorial(i1, visited1, workerID, []string{})
			rightPaths := dfsCombinatorial(i2, visited2, workerID, []string{})
			
			for _, l := range leftPaths {
				for _, r := range rightPaths {
					mutex.Lock()
					if resultCount >= maxResults {
						mutex.Unlock()
						return
					}
					resultCount++
					mutex.Unlock()
					
					combinedSteps := append(append([]string{}, l.Steps...), r.Steps...)
					step := fmt.Sprintf("%s + %s => %s (worker %d)", i1, i2, target, workerID)
					combinedSteps = append(combinedSteps, step)
					
					resultsChan <- DFSResult{
						Found:       true,
						Steps:       append(path, combinedSteps...),
						WorkerID:    workerID,
						NodesVisited: l.NodesVisited + r.NodesVisited + 1,
						Runtime:     time.Since(startTime).Seconds(),
					}
				}
			}
		}(i1, i2, workerID)
	}
}

func dfsCombinatorial(target string, visited map[string]bool, workerID int, path []string) []DFSResult {
	if visited[target] {
		return nil
	}
	visited[target] = true
	
	startTime := time.Now()
	
	recipes, ok := recipesMap[target]
	if !ok {
		return []DFSResult{{
			Found:       true,
			Steps:       []string{},
			WorkerID:    workerID,
			NodesVisited: 1,
			Runtime:     time.Since(startTime).Seconds(),
		}}
	}
	
	elementTier := tierMap[target]
	var results []DFSResult
	for _, ingr := range recipes {
		i1, i2 := ingr[0], ingr[1]
		t1, t2 := tierMap[i1], tierMap[i2]
		
		if t1 >= elementTier || t2 >= elementTier {
			continue
		}
		
		mutex.Lock()
		subWorkerLeft := globalWorkerID
		globalWorkerID++
		subWorkerRight := globalWorkerID
		globalWorkerID++
		mutex.Unlock()
		
		visited1 := mapCopy(visited)
		visited2 := mapCopy(visited)
		leftPaths := dfsCombinatorial(i1, visited1, subWorkerLeft, []string{})
		rightPaths := dfsCombinatorial(i2, visited2, subWorkerRight, []string{})
		
		for _, l := range leftPaths {
			for _, r := range rightPaths {
				combinedSteps := append(append([]string{}, l.Steps...), r.Steps...)
				step := fmt.Sprintf("%s + %s => %s (worker %d)", i1, i2, target, workerID)
				combinedSteps = append(combinedSteps, step)
				
				results = append(results, DFSResult{
					Found:       true,
					Steps:       append(path, combinedSteps...),
					WorkerID:    workerID,
					NodesVisited: l.NodesVisited + r.NodesVisited + 1,
					Runtime:     time.Since(startTime).Seconds(),
				})
			}
		}
	}
	return results
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run dfsMultiple2.go <target_element> <max_path>")
		return
	}
	
	// Load recipes dari utils.go
	recipes, err := loadRecipes("recipes.json")
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}
	buildRecipeMap(recipes) // Panggil fungsi dari utils.go
	
	target := strings.ToLower(os.Args[1])
	maxPath, err := strconv.Atoi(os.Args[2])
	if err != nil || maxPath <= 0 {
		fmt.Println("Invalid max-path value")
		return
	}
	maxResults = maxPath
	
	spawnDFSForEachRecipe(target, make(map[string]bool), []string{})
	
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	var allResults []DFSResult
	for res := range resultsChan {
		allResults = append(allResults, res)
	}
	
	uniqueMap := make(map[string]bool)
	var uniqueResults []DFSResult
	for _, res := range allResults {
		key := strings.Join(res.Steps, "|")
		if !uniqueMap[key] {
			uniqueMap[key] = true
			uniqueResults = append(uniqueResults, res)
		}
	}
	
	// Output hasil
	for _, res := range uniqueResults {
		fmt.Printf("=== Path (Worker %d) ===\n", res.WorkerID)
		fmt.Printf("Runtime: %.4fs | Nodes: %d\n", res.Runtime, res.NodesVisited)
		for _, step := range res.Steps {
			fmt.Println("->", step)
		}
		fmt.Println()
	}
	
	// Output JSON
	jsonOut, _ := json.MarshalIndent(uniqueResults, "", "  ")
	fmt.Println("JSON Output:")
	fmt.Println(string(jsonOut))
}
