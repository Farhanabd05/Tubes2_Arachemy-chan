package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"
)

// Recipe represents a single recipe from the JSON file
type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	ImgUrl1     string `json:"ImgUrl1"`
	ImgUrl2     string `json:"ImgUrl2"`
	Type        int    `json:"Type"`
}

// PathJob represents a search task for a specific path approach
type PathJob struct {
	JobID   int
	Target  string
	Method  string  // "dfs", "bfs", "bidirectional", etc.
	Depth   int     // Maximum search depth
	StartEl string  // Optional starting element for constraints
}

// Path represents a single creation path
type Path struct {
	Steps []string
	Elements map[string]bool
}

// JobResult contains the result of a processed job
type JobResult struct {
	JobID    int
	WorkerID int
	Target   string
	Method   string
	Found    bool
	Path     Path
	Err      error
	Duration time.Duration
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
	resultsMutex sync.Mutex // For thread-safe access to results collection
)

	type SearchState struct {
		Element   string
		Path      []string
		Elements  map[string]bool
		Depth     int
	}

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

// findPathDFS uses DFS to find a path with constraints on starting elements
func findPathDFS(target string, maxDepth int, startEl string) (Path, bool) {
	fmt.Printf("[DEBUG] Starting DFS for target: %s (max depth: %d, start: %s)\n", 
		target, maxDepth, startEl)
	
	target = strings.ToLower(target)
	
	// If target is a base element, return empty path
	if baseElements[target] {
		return Path{
			Steps: []string{},
			Elements: map[string]bool{target: true},
		}, true
	}
	
	// Set of valid starting elements
	startElements := make(map[string]bool)
	if startEl != "" {
		startElements[strings.ToLower(startEl)] = true
	} else {
		for base := range baseElements {
			startElements[base] = true
		}
	}
	
	// Stack for DFS
	type StackItem struct {
		Path      []string
		Current   string
		Depth     int
		Visited   map[string]bool
		Elements  map[string]bool
	}
	
	var stack []StackItem
	
	// Add starting points to stack
	for start := range startElements {
		visited := map[string]bool{start: true}
		elements := map[string]bool{start: true}
		stack = append(stack, StackItem{
			Path:     []string{},
			Current:  start,
			Depth:    0,
			Visited:  visited,
			Elements: elements,
		})
	}
	
	// Random seed to encourage path diversity among workers
	seed := time.Now().UnixNano() % 997
	
	// DFS
	for len(stack) > 0 {
		// Pop from stack
		idx := len(stack) - 1
		current := stack[idx]
		stack = stack[:idx]
		
		// Check depth limit
		if current.Depth >= maxDepth {
			continue
		}
		
		// Check if reached target
		if current.Current == target {
			return Path{
				Steps:    current.Path,
				Elements: current.Elements,
			}, true
		}
		
		// Try to create new elements by combining with known elements
		mutex.RLock()
		
		// Get list of known recipes that use current element
		producible := make([]string, 0)
		for element, recipes := range recipesMap {
			if current.Visited[element] {
				continue
			}
			
			for _, recipe := range recipes {
				ingredient1 := recipe[0]
				ingredient2 := recipe[1]
				
				// Skip if tier constraint is violated
				resultTier := tierMap[element]
				if tierMap[ingredient1] >= resultTier || tierMap[ingredient2] >= resultTier {
					continue
				}
				
				// Check if we can use current element to produce this new element
				if (ingredient1 == current.Current && current.Elements[ingredient2]) ||
				   (ingredient2 == current.Current && current.Elements[ingredient1]) {
					producible = append(producible, element)
					break
				}
			}
		}
		mutex.RUnlock()
		
		// Randomize order for diversity in worker paths (using seed for determinism per worker)
		for i := 0; i < len(producible); i++ {
			j := (int(seed) + i) % len(producible)
			if i != j {
				producible[i], producible[j] = producible[j], producible[i]
			}
		}
		
		// Add producible elements to stack
		for _, nextEl := range producible {
			mutex.RLock()
			// Find a valid recipe
			var validRecipe []string
			
			for _, recipe := range recipesMap[nextEl] {
				ingredient1 := recipe[0]
				ingredient2 := recipe[1]
				
				if (current.Elements[ingredient1] && current.Elements[ingredient2]) &&
				   (tierMap[ingredient1] < tierMap[nextEl] && tierMap[ingredient2] < tierMap[nextEl]) {
					validRecipe = recipe
					break
				}
			}
			mutex.RUnlock()
			
			// If we found a valid recipe, use it
			if validRecipe != nil {
				ingredient1 := validRecipe[0]
				ingredient2 := validRecipe[1]
				
				// Create new path and visited sets
				newPath := make([]string, len(current.Path))
				copy(newPath, current.Path)
				newPath = append(newPath, fmt.Sprintf("%s + %s = %s", 
					ingredient1, ingredient2, nextEl))
				
				newVisited := make(map[string]bool)
				for k, v := range current.Visited {
					newVisited[k] = v
				}
				newVisited[nextEl] = true
				
				newElements := make(map[string]bool)
				for k, v := range current.Elements {
					newElements[k] = v
				}
				newElements[nextEl] = true
				
				// Add to stack
				stack = append(stack, StackItem{
					Path:     newPath,
					Current:  nextEl,
					Depth:    current.Depth + 1,
					Visited:  newVisited,
					Elements: newElements,
				})
			}
		}
	}
	
	return Path{}, false
}

// findPathBFS uses BFS to find a path with constraints
func findPathBFS(target string, maxDepth int, startEl string) (Path, bool) {
	fmt.Printf("[DEBUG] Starting BFS for target: %s (max depth: %d, start: %s)\n", 
		target, maxDepth, startEl)
	
	target = strings.ToLower(target)
	
	// If target is a base element, return empty path
	if baseElements[target] {
		return Path{
			Steps: []string{},
			Elements: map[string]bool{target: true},
		}, true
	}
	
	// Set of valid starting elements
	startElements := make(map[string]bool)
	if startEl != "" {
		startElements[strings.ToLower(startEl)] = true
	} else {
		for base := range baseElements {
			startElements[base] = true
		}
	}
	
	// Queue for BFS
	type QueueItem struct {
		Path      []string
		Current   string
		Depth     int
		Visited   map[string]bool
		Elements  map[string]bool
	}
	
	var queue []QueueItem
	discovered := make(map[string]bool)
	
	// Add starting points to queue
	for start := range startElements {
		visited := map[string]bool{start: true}
		elements := map[string]bool{start: true}
		discovered[start] = true
		
		queue = append(queue, QueueItem{
			Path:     []string{},
			Current:  start,
			Depth:    0,
			Visited:  visited,
			Elements: elements,
		})
	}
	
	// Random seed to encourage path diversity among workers
	seed := time.Now().UnixNano() % 997
	
	// BFS
	for len(queue) > 0 {
		// Pop from queue
		current := queue[0]
		queue = queue[1:]
		
		// Check depth limit
		if current.Depth >= maxDepth {
			continue
		}
		
		// Check if reached target
		if current.Current == target {
			return Path{
				Steps:    current.Path,
				Elements: current.Elements,
			}, true
		}
		
		// Get potential combinations
		potentialCombos := make([]struct{
			Element    string
			Ingredient1 string
			Ingredient2 string
		}, 0)
		
		mutex.RLock()
		// Try to combine with other known elements
		for knownEl := range current.Elements {
			// Look for recipes that use both current.Current and knownEl
			for element, recipes := range recipesMap {
				if discovered[element] {
					continue
				}
				
				for _, recipe := range recipes {
					ingredient1 := recipe[0]
					ingredient2 := recipe[1]
					
					// Skip if tier constraint is violated
					resultTier := tierMap[element]
					if (tierMap[ingredient1] >= resultTier || tierMap[ingredient2] >= resultTier) {
						continue
					}
					
					// Check if we can use the combination
					if (ingredient1 == knownEl && ingredient2 == current.Current) ||
					   (ingredient1 == current.Current && ingredient2 == knownEl) {
						potentialCombos = append(potentialCombos, struct{
							Element    string
							Ingredient1 string
							Ingredient2 string
						}{
							Element:    element,
							Ingredient1: ingredient1,
							Ingredient2: ingredient2,
						})
					}
				}
			}
		}
		mutex.RUnlock()
		
		// Randomize order for diversity in worker paths
		for i := 0; i < len(potentialCombos); i++ {
			j := (int(seed) + i) % len(potentialCombos)
			if i != j {
				potentialCombos[i], potentialCombos[j] = potentialCombos[j], potentialCombos[i]
			}
		}
		
		// Process potential combinations
		for _, combo := range potentialCombos {
			nextEl := combo.Element
			ingredient1 := combo.Ingredient1
			ingredient2 := combo.Ingredient2
			
			if !discovered[nextEl] {
				discovered[nextEl] = true
				
				// Create new path and visited sets
				newPath := make([]string, len(current.Path))
				copy(newPath, current.Path)
				newPath = append(newPath, fmt.Sprintf("%s + %s = %s", 
					ingredient1, ingredient2, nextEl))
				
				newVisited := make(map[string]bool)
				for k, v := range current.Visited {
					newVisited[k] = v
				}
				newVisited[nextEl] = true
				
				newElements := make(map[string]bool)
				for k, v := range current.Elements {
					newElements[k] = v
				}
				newElements[nextEl] = true
				
				// Add to queue
				queue = append(queue, QueueItem{
					Path:     newPath,
					Current:  nextEl,
					Depth:    current.Depth + 1,
					Visited:  newVisited,
					Elements: newElements,
				})
				
				// Early exit if target found
				if nextEl == target {
					return Path{
						Steps:    newPath,
						Elements: newElements,
					}, true
				}
			}
		}
	}
	
	return Path{}, false
}

// findPathBidirectional uses bidirectional search to find paths
func findPathBidirectional(target string, maxDepth int, startEl string) (Path, bool) {
	fmt.Printf("[DEBUG] Starting bidirectional search for target: %s (max depth: %d, start: %s)\n", 
		target, maxDepth, startEl)
	
	target = strings.ToLower(target)
	
	// If target is a base element, return empty path
	if baseElements[target] {
		return Path{
			Steps: []string{},
			Elements: map[string]bool{target: true},
		}, true
	}
	
	// Set of valid starting elements
	startElements := make(map[string]bool)
	if startEl != "" {
		startElements[strings.ToLower(startEl)] = true
	} else {
		for base := range baseElements {
			startElements[base] = true
		}
	}

	
	// Forward search from base elements
	forwardQueue := make([]SearchState, 0)
	forwardVisited := make(map[string]SearchState)
	
	// Add base elements to forward queue
	for start := range startElements {
		state := SearchState{
			Element:  start,
			Path:     []string{},
			Elements: map[string]bool{start: true},
			Depth:    0,
		}
		forwardQueue = append(forwardQueue, state)
		forwardVisited[start] = state
	}
	
	// Backward search from target
	backwardQueue := make([]SearchState, 0)
	backwardVisited := make(map[string]SearchState)
	
	state := SearchState{
		Element:  target,
		Path:     []string{},
		Elements: map[string]bool{target: true},
		Depth:    0,
	}
	backwardQueue = append(backwardQueue, state)
	backwardVisited[target] = state
	
	// Bidirectional search
	for len(forwardQueue) > 0 && len(backwardQueue) > 0 {
		// Forward search step
		if len(forwardQueue) > 0 {
			current := forwardQueue[0]
			forwardQueue = forwardQueue[1:]
			
			// Check if we've reached max depth
			if current.Depth >= maxDepth/2 {
				continue
			}
			
			// Check if we've met the backward search
			if backState, exists := backwardVisited[current.Element]; exists {
				// We found a path! Combine forward and backward
				return combinePaths(current, backState, target), true
			}
			
			// Expand forward
			mutex.RLock()
			// Try combinations with other known elements
			for knownEl := range current.Elements {
				for element, recipes := range recipesMap {
					if forwardVisited[element].Element != "" {
						continue // Already visited
					}
					
					for _, recipe := range recipes {
						ingredient1 := recipe[0]
						ingredient2 := recipe[1]
						
						// Skip if tier constraint is violated
						resultTier := tierMap[element]
						if tierMap[ingredient1] >= resultTier || tierMap[ingredient2] >= resultTier {
							continue
						}
						
						// Check if we can make this recipe
						if (ingredient1 == knownEl && current.Elements[ingredient2]) ||
						   (ingredient2 == knownEl && current.Elements[ingredient1]) {
							
							// Create new state
							newPath := make([]string, len(current.Path))
							copy(newPath, current.Path)
							newPath = append(newPath, fmt.Sprintf("%s + %s = %s", 
								ingredient1, ingredient2, element))
							
							newElements := make(map[string]bool)
							for k, v := range current.Elements {
								newElements[k] = v
							}
							newElements[element] = true
							
							newState := SearchState{
								Element:  element,
								Path:     newPath,
								Elements: newElements,
								Depth:    current.Depth + 1,
							}
							
							forwardQueue = append(forwardQueue, newState)
							forwardVisited[element] = newState
							
							// Check if we've met the backward search
							if backState, exists := backwardVisited[element]; exists {
								mutex.RUnlock()
								return combinePaths(newState, backState, target), true
							}
							
							break // Found one recipe, no need to check others
						}
					}
				}
			}
			mutex.RUnlock()
		}
		
		// Backward search step
		if len(backwardQueue) > 0 {
			current := backwardQueue[0]
			backwardQueue = backwardQueue[1:]
			
			// Check if we've reached max depth
			if current.Depth >= maxDepth/2 {
				continue
			}
			
			// Check if we've met the forward search
			if fwdState, exists := forwardVisited[current.Element]; exists {
				return combinePaths(fwdState, current, target), true
			}
			
			// Expand backward
			mutex.RLock()
			// Find ingredients that can make current element
			recipes, exists := recipesMap[current.Element]
			if exists {
				for _, recipe := range recipes {
					ingredient1 := recipe[0]
					ingredient2 := recipe[1]
					
					// Skip if tier constraint is violated
					resultTier := tierMap[current.Element]
					if tierMap[ingredient1] >= resultTier || tierMap[ingredient2] >= resultTier {
						continue
					}
					
					// Add both ingredients to backward search
					for _, ingredient := range []string{ingredient1, ingredient2} {
						if backwardVisited[ingredient].Element != "" {
							continue // Already visited
						}
						
						// Create new state
						newPath := make([]string, len(current.Path))
						copy(newPath, current.Path)
						// Add in reverse order for backward search
						newPath = append([]string{fmt.Sprintf("%s + %s = %s", 
							ingredient1, ingredient2, current.Element)}, newPath...)
						
						newElements := make(map[string]bool)
						for k, v := range current.Elements {
							newElements[k] = v
						}
						newElements[ingredient] = true
						
						newState := SearchState{
							Element:  ingredient,
							Path:     newPath,
							Elements: newElements,
							Depth:    current.Depth + 1,
						}
						
						backwardQueue = append(backwardQueue, newState)
						backwardVisited[ingredient] = newState
						
						// Check if we've met the forward search
						if fwdState, exists := forwardVisited[ingredient]; exists {
							mutex.RUnlock()
							return combinePaths(fwdState, newState, target), true
						}
					}
				}
			}
			mutex.RUnlock()
		}
	}
	
	return Path{}, false
}

// combinePaths combines forward and backward paths
func combinePaths(forward, backward SearchState, target string) Path {
	result := Path{
		Steps:    make([]string, 0),
		Elements: make(map[string]bool),
	}
	
	// Add forward path steps
	result.Steps = append(result.Steps, forward.Path...)
	
	// Add backward path steps (need to reverse them)
	for i := len(backward.Path) - 1; i >= 0; i-- {
		result.Steps = append(result.Steps, backward.Path[i])
	}
	
	// Combine elements
	for k, v := range forward.Elements {
		result.Elements[k] = v
	}
	for k, v := range backward.Elements {
		result.Elements[k] = v
	}
	
	return result
}

// pathWorker is a worker that finds paths using different methods
func pathWorker(workerID int, jobChan <-chan PathJob, resultChan chan<- JobResult, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Worker %d] Started\n", workerID)
	
	for job := range jobChan {
		fmt.Printf("[Worker %d] Processing path job: target=%s, method=%s, depth=%d\n", 
			workerID, job.Target, job.Method, job.Depth)
		
		startTime := time.Now()
		var path Path
		var found bool
		var err error
		
		// Use recover to catch panics
		func() {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					err = fmt.Errorf("panic: %v\nStack trace:\n%s", r, stack)
					fmt.Printf("[Worker %d] PANIC in job %d: %v\n%s\n", 
						workerID, job.JobID, r, stack)
				}
			}()
			
			// Choose search method based on job
			switch job.Method {
			case "dfs":
				path, found = findPathDFS(job.Target, job.Depth, job.StartEl)
			case "bfs":
				path, found = findPathBFS(job.Target, job.Depth, job.StartEl)
			case "bidirectional":
				path, found = findPathBidirectional(job.Target, job.Depth, job.StartEl)
			default:
				err = fmt.Errorf("unsupported path finding method: %s", job.Method)
			}
		}()
		
		duration := time.Since(startTime)
		
		// Send result
		resultChan <- JobResult{
			JobID:    job.JobID,
			WorkerID: workerID,
			Target:   job.Target,
			Method:   job.Method,
			Found:    found,
			Path:     path,
			Err:      err,
			Duration: duration,
		}
		
		fmt.Printf("[Worker %d] Completed job %d in %v\n", workerID, job.JobID, duration)
	}
	
	fmt.Printf("[Worker %d] Exiting\n", workerID)
}

// formatPath returns a formatted string representation of a path
func formatPath(path Path, method string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Path found using %s (%d steps):\n", method, len(path.Steps)))
	
	for i, step := range path.Steps {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step))
	}
	
	// List elements used
	elements := make([]string, 0, len(path.Elements))
	for elem := range path.Elements {
		elements = append(elements, elem)
	}
	sort.Strings(elements)
	sb.WriteString("  Elements used: ")
	sb.WriteString(strings.Join(elements, ", "))
	
	return sb.String()
}

func main() {
	fmt.Println("[Main] Starting Little Alchemy Parallel Pathfinder")
	
	// Load recipes
	recipesFile := "test/data/recipes.json"
	recipes, err := loadRecipes(recipesFile)
	if err != nil {
		fmt.Printf("[FATAL] Error loading recipes: %v\n", err)
		return
	}
	
	// Build recipe data structures
	buildRecipeMap(recipes)
	buildReverseGraph()
	
	// Create channels for jobs and results
	workerCount := 8 // Use more workers for better parallelism
	fmt.Printf("[Main] Starting %d workers\n", workerCount)
	
	jobChan := make(chan PathJob, workerCount*2)
	resultChan := make(chan JobResult, workerCount*2)
	
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go pathWorker(i, jobChan, resultChan, &wg)
	}
	
	// Target elements to find paths for
	targets := []string{"brick"}
	
	// Create jobs for each target using different methods and starting points
	jobID := 1
	go func() {
		// For each target
		for _, target := range targets {
			fmt.Printf("[Main] Creating jobs for target: %s\n", target)
			
			// Different methods to find diverse paths
			methods := []string{"bfs"}
			depths := []int{15}
			
			// Different starting elements to force path diversity
			startingElements := []string{"", "fire", "water", "earth", "air"}
			
			// Create a variety of jobs to find different paths
			jobsCreated := 0
			for _, method := range methods {
				for _, depth := range depths {
					for _, start := range startingElements {
						// Skip some combinations to avoid too many jobs
						if jobsCreated >= 6 {
							break
						}
						
						fmt.Printf("[Main] Sending job %d: %s search for %s (depth: %d, start: %s)\n", 
							jobID, method, target, depth, start)
						
						jobChan <- PathJob{
							JobID:   jobID,
							Target:  target,
							Method:  method,
							Depth:   depth,
							StartEl: start,
						}
						
						jobID++
						jobsCreated++
					}
					if jobsCreated >= 6 {
						break
					}
				}
				if jobsCreated >= 6 {
					break
				}
			}
		}
		
		fmt.Println("[Main] All jobs sent, closing job channel")
		close(jobChan)
	}()
	
	// Wait for all workers to finish
	go func() {
		wg.Wait()
		fmt.Println("[Main] All workers finished, closing result channel")
		close(resultChan)
	}()
	
	// Process results and collect unique paths
	fmt.Println("[Main] Processing results")
	
	// Map to store unique paths for each target
	targetPaths := make(map[string]map[string]Path) // target -> path signature -> path
	for _, target := range targets {
		targetPaths[target] = make(map[string]Path)
	}
	
	// Process results as they come in
	for result := range resultChan {
		if result.Err != nil {
			fmt.Printf("[Result] Job %d (%s for %s) FAILED by worker %d: %v\n", 
				result.JobID, result.Method, result.Target, result.WorkerID, result.Err)
			continue
		}
		
		if result.Found {
			fmt.Printf("[Result] Job %d (%s for %s) SUCCEEDED by worker %d in %v\n", 
				result.JobID, result.Method, result.Target, result.WorkerID, result.Duration)
			
			// Create a signature for this path to detect duplicates
			pathSignature := strings.Join(result.Path.Steps, "|")
			
			// Add to paths map if unique
			resultsMutex.Lock()
			if _, exists := targetPaths[result.Target][pathSignature]; !exists {
				targetPaths[result.Target][pathSignature] = result.Path
				fmt.Printf("[Result] Found new unique path for %s using %s!\n", 
					result.Target, result.Method)
			}
			resultsMutex.Unlock()
		} else {
			fmt.Printf("[Result] Job %d (%s for %s) by worker %d: Not found (took %v)\n", 
				result.JobID, result.Method, result.Target, result.WorkerID, result.Duration)
		}
	}
	
	// Display all unique paths found
	fmt.Println("\n[Main] === Final Results ===")
	
	for _, target := range targets {
		paths := targetPaths[target]
		fmt.Printf("\n[Result] Found %d unique paths to create '%s':\n\n", len(paths), target)
		for _, path := range paths {
			fmt.Printf("[Result] %s\n", path)
		}
	}
}