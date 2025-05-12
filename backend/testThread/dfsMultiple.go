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

// Recipe represents a single recipe from the JSON file
type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Type        int    `json:"Type"`
}

// Job represents a single task to be processed by a worker
type Job struct {
	JobID   int
	JobType string
	Target  string
}

type Result struct {
	Found bool     `json:"found"`
	Steps []string `json:"steps"`
	Runtime time.Duration `json:"runtime"`
	NodesVisited int `json:"nodesVisited"`
}

// JobResult contains the result of a processed job
type JobResult struct {
	JobID    int
	WorkerID int
	JobType  string
	Target   string
	Found    bool
	Steps    []string
	Err      error
	Duration time.Duration // Added to track job execution time
}

// Global variables for recipe data
var (
	recipesMap   map[string][][]string
	tierMap      map[string]int
	revGraph     map[string][]string
	baseElements = map[string]bool{
		"fire": true, "water": true, "earth": true, "air": true, "time": true,
	}
	mutex sync.RWMutex // For thread-safe access to recipe maps
	printCount = 0
	maxPrints  = 200
)

// loadRecipes loads recipe data from a JSON file
func loadRecipes(file string) ([]Recipe, error) {
	fmt.Printf("[DEBUG] Loading recipes from %s\n", file)
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read file: %v\n", err)
		return nil, err
	}
	
	var recipes []Recipe
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		fmt.Printf("[ERROR] Failed to unmarshal JSON: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("[DEBUG] Successfully loaded %d recipes\n", len(recipes))
	return recipes, nil
}

// buildRecipeMap constructs the recipe and tier maps
func buildRecipeMap(recipes []Recipe) {
	fmt.Println("[DEBUG] Building recipe map")
	recipesMap = make(map[string][][]string)
	tierMap = make(map[string]int)
	
	for _, r := range recipes {
		element := strings.ToLower(r.Element)
		ingr1 := strings.ToLower(r.Ingredient1)
		ingr2 := strings.ToLower(r.Ingredient2)
		ingr := []string{ingr1, ingr2}
		
		recipesMap[element] = append(recipesMap[element], ingr)
		tierMap[element] = r.Type
	}
	
	fmt.Printf("[DEBUG] Recipe map built with %d elements\n", len(recipesMap))
}

// buildReverseGraph constructs a reverse lookup graph for efficient path finding
func buildReverseGraph() {
	fmt.Println("[DEBUG] Building reverse graph")
	revGraph = make(map[string][]string)
	
	for result, recipes := range recipesMap {
		for _, ingr := range recipes {
			// Add result to each ingredient's created elements list
			revGraph[ingr[0]] = appendUnique(revGraph[ingr[0]], result)
			revGraph[ingr[1]] = appendUnique(revGraph[ingr[1]], result)
		}
	}
	
	fmt.Printf("[DEBUG] Reverse graph built with %d elements\n", len(revGraph))
}

// appendUnique adds an element to a slice if it doesn't already exist
func appendUnique(slice []string, element string) []string {
	for _, e := range slice {
		if e == element {
			return slice
		}
	}
	return append(slice, element)
}

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
