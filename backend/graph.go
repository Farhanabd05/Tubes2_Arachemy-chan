package main

// import "strings"
import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type State struct {
	Elements map[string]bool
	Path []string
}
var combinations map[[2]string]string

func init() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		combinations = readRecipes("./test/data/recipes.json")
	}()
	wg.Wait()
}

func loadJson(filename string, data interface{}) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		panic(err)
	}
}

func readRecipes(filename string) map[[2]string]string {
	var recipes []map[string]interface{}
	loadJson(filename, &recipes)
	result := make(map[[2]string]string)
	for _, recipe := range recipes {
		ingredient1 := recipe["Ingredient1"].(string)
		ingredient2 := recipe["Ingredient2"].(string)
		result[normalize(ingredient1, ingredient2)] = recipe["Element"].(string)
	}
	return result
}

// Normalisasi pasangan agar urutan tidak masalah
func normalize(a, b string) [2]string {
	if a < b {
		return [2]string{a, b}
	}
	return [2]string{b, a}
}

func getCombination(element1, element2 string) (string, bool) {
	key := normalize(element1, element2)
	result, exists := combinations[key]
	return result, exists
}

// Fungsi utama BFS
func findWithBFS(target string) ([]string, bool) {
	initialElements := []string{"fire", "water", "earth", "air"}

	// Antrian BFS
	queue := []State{
		{
			Elements: sliceToSet(initialElements),
			Path:     []string{},
		},
	}

	// Set untuk mengecek state yang sudah dikunjungi
	visitedStates := map[string]bool{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Jika sudah ditemukan
		if current.Elements[target] {
			return current.Path, true
		}

		stateKey := stateToKey(current.Elements)
		if visitedStates[stateKey] {
			continue
		}
		visitedStates[stateKey] = true

		elements := keys(current.Elements)

		// Coba semua kombinasi dari elemen yang ada saat ini
		for i := 0; i < len(elements); i++ {
			for j := i + 1; j < len(elements); j++ {
				e1, e2 := elements[i], elements[j]
				comb, ok := getCombination(e1, e2)
				if ok && !current.Elements[comb] {
					newElements := copySet(current.Elements)
					newElements[comb] = true

					newPath := make([]string, len(current.Path), len(current.Path)+1)
					copy(newPath, current.Path)
					newPath[len(current.Path)] = fmt.Sprintf("%s + %s = %s", e1, e2, comb)

					queue = append(queue, State{
						Elements: newElements,
						Path:     newPath,
					})
				}
			}
		}
	}

	return nil, false
}

// ✅ DFS implementation
func findWithDFS(target string) ([]string, bool) {
	startElements := []string{"fire", "water", "earth", "air"}
	visited := make(map[string]bool)
	steps := []string{}
	pathFound := dfsHelper(startElements, target, visited, &steps)
	return steps, pathFound
}

func dfsHelper(current []string, target string, visited map[string]bool, steps *[]string) bool {
	// Check if target already in current state
	for _, el := range current {
		if el == target {
			return true
		}
	}

	for i := 0; i < len(current); i++ {
		for j := i + 1; j < len(current); j++ {
			a, b := current[i], current[j]
			combined, ok := getCombination(a, b)
			if ok && !visited[combined] {
				visited[combined] = true
				*steps = append(*steps, fmt.Sprintf("%s + %s => %s", a, b, combined))
				newState := append([]string{combined}, current...)
				if dfsHelper(newState, target, visited, steps) {
					return true
				}
				// Backtrack
				*steps = (*steps)[:len(*steps)-1]
			}
		}
	}

	return false
}

// BFS dengan multi-threading
func findWithBFSParallel(target string, numWorkers int) ([]string, bool) {
	initialElements := []string{"fire", "water", "earth", "air"}
	
	// Antrian BFS dengan mutex untuk akses bersamaan
	var queueMutex sync.Mutex
	queue := []State{
		{
			Elements: sliceToSet(initialElements),
			Path:     []string{},
		},
	}
	
	// Set untuk state yang sudah dikunjungi dengan mutex
	var visitedMutex sync.Mutex
	visitedStates := map[string]bool{}
	
	// Channel untuk hasil
	resultChan := make(chan []string, numWorkers)
	foundChan := make(chan bool, numWorkers)
	
	// WaitGroup untuk menunggu semua worker selesai
	var wg sync.WaitGroup
	
	// Flag untuk menandakan pencarian sudah selesai
	done := false
	var doneMutex sync.Mutex
	
	isDone := func() bool {
		doneMutex.Lock()
		defer doneMutex.Unlock()
		return done
	}
	
	setDone := func() {
		doneMutex.Lock()
		done = true
		doneMutex.Unlock()
	}
	
	// Fungsi worker untuk BFS
	worker := func(id int) {
		defer wg.Done()
		
		for !isDone() {
			// Ambil state dari queue dengan aman
			queueMutex.Lock()
			if len(queue) == 0 {
				queueMutex.Unlock()
				return // Queue kosong, worker selesai
			}
			
			current := queue[0]
			queue = queue[1:]
			queueMutex.Unlock()
			
			// Cek apakah target sudah ditemukan
			if current.Elements[target] {
				resultChan <- current.Path
				foundChan <- true
				setDone()
				return
			}
			
			// Cek visited state
			stateKey := stateToKey(current.Elements)
			visitedMutex.Lock()
			visited := visitedStates[stateKey]
			if !visited {
				visitedStates[stateKey] = true
			}
			visitedMutex.Unlock()
			
			if visited {
				continue
			}
			
			elements := keys(current.Elements)
			
			// Coba semua kombinasi
			for i := 0; i < len(elements); i++ {
				for j := i + 1; j < len(elements); j++ {
					if isDone() {
						return
					}
					
					e1, e2 := elements[i], elements[j]
					comb, ok := getCombination(e1, e2)
					
					if ok && !current.Elements[comb] {
						newElements := copySet(current.Elements)
						newElements[comb] = true
						
						newPath := append([]string{}, current.Path...)
						newPath = append(newPath, fmt.Sprintf("%s + %s = %s", e1, e2, comb))
						
						queueMutex.Lock()
						queue = append(queue, State{
							Elements: newElements,
							Path:     newPath,
						})
						queueMutex.Unlock()
					}
				}
			}
		}
	}
	
	// Mulai worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i)
	}
	
	// Goroutine untuk menunggu semua worker selesai
	go func() {
		wg.Wait()
		// Jika semua worker selesai dan tidak ada hasil
		if !isDone() {
			foundChan <- false
		}
		close(resultChan)
		close(foundChan)
	}()
	
	// Tunggu hasil
	found := <-foundChan
	if found {
		return <-resultChan, true
	}
	
	return nil, false
}