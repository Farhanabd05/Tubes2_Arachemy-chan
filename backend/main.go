package main

import (
	"fmt"
	"runtime"
)



func main() {
    // Load recipes dan build data structure
    recipesFile := "test/data/recipes.json"
    recipes, err := loadRecipes(recipesFile)
    if err != nil {
        fmt.Printf("[FATAL] Error loading recipes: %v\n", err)
        return
    }

    buildRecipeMap(recipes)
    buildReverseGraph()

    // Buat worker pool untuk DFS multiple paths
    numWorkers := runtime.NumCPU() // Gunakan jumlah CPU yang tersedia
    dfsJobs, dfsResults := StartDFSMultipleWorkerPool(numWorkers)

    // Submit jobs
    targets := []string{"blade"}
    maxPathsPerTarget := 5

    go func() {
        for i, target := range targets {
            dfsJobs <- DFSMultipleJob{
                Target:   target,
                MaxPaths: maxPathsPerTarget,
                JobID:    i + 1,
            }
        }
        close(dfsJobs)
    }()

    // Process results
    for result := range dfsResults {
        if result.Found {
            fmt.Printf("[Result] Target: %s - Found %d paths in %v using DFS\n",
                result.Target, len(result.Paths), result.Duration)
            // Tampilkan paths yang ditemukan
            for i, path := range result.Paths {
                fmt.Printf("Path %d: ", i+1)
                for j, step := range path {
                    if j > 0 {
                        fmt.Print(" -> ")
                    }
                    fmt.Print(step)
                }
                fmt.Println()
            }
        } else {
            fmt.Printf("[Result] Target: %s - No paths found using DFS\n", result.Target)
        }
    }
}
