package main


// import "strings"
import (
	"fmt"
)

type State struct {
	Elements map[string]bool
	Path []string
}
var combinations = map[[2]string]string{
	normalize("fire", "air"):    "energy",
	normalize("water", "earth"): "mud",
	normalize("energy", "mud"):  "life",
	normalize("life", "earth"):  "human",
	normalize("air", "water"):   "rain",
	normalize("earth", "rain"):  "plant",
	normalize("plant", "human"): "farmer",
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

					newPath := append([]string{}, current.Path...)
					newPath = append(newPath, fmt.Sprintf("%s + %s = %s", e1, e2, comb))

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