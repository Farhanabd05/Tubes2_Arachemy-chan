// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"sync"
// )

// type Result struct {
// 	Found        bool     `json:"found"`
// 	Steps        []string `json:"steps"`
// 	WorkerID     int      `json:"worker_id"`
// 	NodesVisited int      `json:"nodes_visited"`
// 	//tambahkan time kalau bisa
// }

// var tierMap = map[string]int{
// 	"air":         0,
// 	"earth":       0,
// 	"fire":        0,
// 	"water":       0,
// 	"energy":      1,
// 	"steam":       1,
// 	"metal":       2,
// 	"cloud":       2,
// 	"pressure":    2,
// 	"steel":       3,
// 	"lightning":   3,
// 	"electricity": 4,
// }

// var recipesMap = map[string][][]string{
// 	"energy":      {{"fire", "air"}},
// 	"steam":       {{"water", "fire"}, {"fire", "fire"}},
// 	"metal":       {{"earth", "energy"}},
// 	"pressure":    {{"air", "air"}},
// 	"cloud":       {{"energy", "steam"}, {"steam", "steam"}},
// 	"steel":       {{"metal", "fire"}},
// 	"lightning":   {{"energy", "cloud"}},
// 	"electricity": {{"metal", "lightning"}, {"steel", "lightning"}},
// }

// var globalWorkerID = 1
// var mutex sync.Mutex
// var resultsChan = make(chan Result, 1000)
// var wg sync.WaitGroup
// var maxResults int = 5
// var resultCount int = 0

// func mapCopy(src map[string]bool) map[string]bool {
// 	dst := make(map[string]bool)
// 	for k, v := range src {
// 		dst[k] = v
// 	}
// 	return dst
// }

// func spawnDFSForEachRecipe(target string, visited map[string]bool, path []string) {
// 	recipes := recipesMap[target]
// 	for _, ingr := range recipes {
// 		mutex.Lock()
// 		if resultCount >= maxResults {
// 			mutex.Unlock()
// 			break
// 		}
// 		workerID := globalWorkerID
// 		globalWorkerID++
// 		mutex.Unlock()

// 		i1, i2 := ingr[0], ingr[1]

// 		wg.Add(1)
// 		go func(i1, i2 string, workerID int) {
// 			defer wg.Done()
// 			visited1 := mapCopy(visited)
// 			visited2 := mapCopy(visited)

// 			leftPaths := dfsCombinatorial(i1, visited1, workerID, []string{})
// 			rightPaths := dfsCombinatorial(i2, visited2, workerID, []string{})

// 			for _, l := range leftPaths {
// 				for _, r := range rightPaths {
// 					mutex.Lock()
// 					if resultCount >= maxResults {
// 						mutex.Unlock()
// 						return
// 					}
// 					resultCount++
// 					mutex.Unlock()

// 					combinedSteps := append([]string{}, l.Steps...)
// 					combinedSteps = append(combinedSteps, r.Steps...)
// 					step := fmt.Sprintf("%s + %s => %s (worker %d)", i1, i2, target, workerID)
// 					combinedSteps = append(combinedSteps, step)

// 					resultsChan <- Result{
// 						Found:        true,
// 						Steps:        append(path, combinedSteps...),
// 						WorkerID:     workerID,
// 						NodesVisited: l.NodesVisited + r.NodesVisited + 1,
// 					}
// 				}
// 			}
// 		}(i1, i2, workerID)
// 	}
// }

// func dfsCombinatorial(target string, visited map[string]bool, workerID int, path []string) []Result {
// 	if visited[target] {
// 		return nil
// 	}
// 	visited[target] = true

// 	recipes, ok := recipesMap[target]
// 	if !ok {
// 		return []Result{{
// 			Found:        true,
// 			Steps:        []string{},
// 			WorkerID:     workerID,
// 			NodesVisited: 1,
// 		}}
// 	}

// 	elementTier := tierMap[target]
// 	var results []Result

// 	for _, ingr := range recipes {
// 		i1, i2 := ingr[0], ingr[1]
// 		t1, t2 := tierMap[i1], tierMap[i2]
// 		if t1 >= elementTier || t2 >= elementTier {
// 			continue
// 		}

// 		mutex.Lock()
// 		subWorkerLeft := globalWorkerID
// 		globalWorkerID++
// 		subWorkerRight := globalWorkerID
// 		globalWorkerID++
// 		mutex.Unlock()

// 		visited1 := mapCopy(visited)
// 		visited2 := mapCopy(visited)

// 		leftPaths := dfsCombinatorial(i1, visited1, subWorkerLeft, []string{})
// 		rightPaths := dfsCombinatorial(i2, visited2, subWorkerRight, []string{})

// 		for _, l := range leftPaths {
// 			for _, r := range rightPaths {
// 				combinedSteps := append([]string{}, l.Steps...)
// 				combinedSteps = append(combinedSteps, r.Steps...)
// 				step := fmt.Sprintf("%s + %s => %s (worker %d)", i1, i2, target, workerID)
// 				combinedSteps = append(combinedSteps, step)

// 				results = append(results, Result{
// 					Found:        true,
// 					Steps:        append(path, combinedSteps...),
// 					WorkerID:     workerID,
// 					NodesVisited: l.NodesVisited + r.NodesVisited + 1,
// 				})
// 			}
// 		}
// 	}
// 	return results
// }

// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Println("Usage: go run dfsMultiple.go <target-element> <max-path>")
// 		return
// 	}
// 	target := os.Args[1]
// 	limit, err := strconv.Atoi(os.Args[2])
// 	if err != nil || limit <= 0 {
// 		fmt.Println("Invalid max-path value")
// 		return
// 	}
// 	maxResults = limit

// 	spawnDFSForEachRecipe(target, map[string]bool{}, []string{})

// 	go func() {
// 		wg.Wait()
// 		close(resultsChan)
// 	}()

// 	var allResults []Result
// 	for res := range resultsChan {
// 		allResults = append(allResults, res)
// 	}

// 	uniqueMap := make(map[string]bool)
// 	var uniqueResults []Result
// 	for _, res := range allResults {
// 		key := strings.Join(res.Steps, "|")
// 		if !uniqueMap[key] {
// 			uniqueMap[key] = true
// 			uniqueResults = append(uniqueResults, res)
// 		}
// 	}

// 	for _, res := range uniqueResults {
// 		fmt.Printf("=== Recipe by worker %d ===\n", res.WorkerID)
// 		fmt.Printf("Found: %v | Nodes Visited: %d\n", res.Found, res.NodesVisited)
// 		for _, step := range res.Steps {
// 			fmt.Println("  ", step)
// 		}
// 		fmt.Println()
// 	}

// 	jsonOut, _ := json.MarshalIndent(uniqueResults, "", "  ")
// 	fmt.Println("JSON Output:")
// 	fmt.Println(string(jsonOut))
// }