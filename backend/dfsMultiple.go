package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)


type ResultDFS struct {
	Found        bool          `json:"found"`
	Steps        []string      `json:"steps"`
	Runtime      time.Duration `json:"runtime"`
	NodesVisited int           `json:"nodesVisited"`
}

type JobResultDFS struct {
	JobID    int
	WorkerID int
	JobType  string
	Target   string
	Found    bool
	Steps    []string
	Err      error
	Duration time.Duration
}

// Globals
var (
	maxResults  = 3
	maxDepth    = 19
	resultMutex sync.Mutex
)



// Utility to copy map
func mapCopy(src map[string]bool) map[string]bool {
	dst := make(map[string]bool)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// DFS recursive with constraint
func dfsCombinatorial(target string, visited map[string]bool, path []string, depth int) []ResultDFS {
	if depth > maxDepth {
		return nil
	}

	if baseElements[target] {
		return []ResultDFS{{
			Found:        true,
			Steps:        []string{},
			NodesVisited: 1,
			Runtime:      0,
		}}
	}

	if visited[target] {
		return nil
	}

	visited[target] = true
	defer delete(visited, target)

	startTime := time.Now()
	elementTier := tierMap[target]
	var results []ResultDFS
	uniquePaths := make(map[string]bool)

	mutex.RLock()
	recipes, ok := recipesMap[target]
	mutex.RUnlock()
	if !ok {
		return nil
	}

	for _, ingr := range recipes {
		i1, i2 := ingr[0], ingr[1]
		t1, t2 := tierMap[i1], tierMap[i2]

		if t1 >= elementTier || t2 >= elementTier {
			continue
		}

		visited1 := mapCopy(visited)
		visited2 := mapCopy(visited)

		left := dfsCombinatorial(i1, visited1, path, depth+1)
		right := dfsCombinatorial(i2, visited2, path, depth+1)

		for _, l := range left {
			for _, r := range right {
				steps := append(append([]string{}, l.Steps...), r.Steps...)
				step := fmt.Sprintf("%s + %s = %s", i1, i2, target)
				steps = append(steps, step)
				if printCount < maxPrints {
					printCount++
					// fmt.Printf("[DEBUG] Processing: %s + %s => %s (Depth: %d)\n", i1, i2, target, depth)
				}

				key := strings.Join(steps, "|")
				if !uniquePaths[key] {
					results = append(results, ResultDFS{
						Found:        true,
						Steps:        steps,
						NodesVisited: l.NodesVisited + r.NodesVisited + 1,
						Runtime:      time.Since(startTime),
					})
					uniquePaths[key] = true
					if len(results) >= maxResults {
						return results
					}
				}
			}
		}
	}

	return results
}

// Worker goroutine
func worker(id int, jobs <-chan Job, results chan<- JobResultDFS, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		allResults := dfsCombinatorial(job.Target, make(map[string]bool), []string{}, 0)

		for _, r := range allResults {
			results <- JobResultDFS{
				JobID:    job.JobID,
				WorkerID: id,
				JobType:  job.JobType,
				Target:   job.Target,
				Found:    r.Found,
				Steps:    r.Steps,
				Duration: r.Runtime,
			}
		}
	}
}

