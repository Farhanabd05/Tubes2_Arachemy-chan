package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// DFSMultipleJob merepresentasikan tugas pencarian multiple path DFS
type DFSMultipleJob struct {
    Target   string
    MaxPaths int
    JobID    int
}

// DFSMultipleResult berisi hasil pencarian multiple path DFS
type DFSMultipleResult struct {
    Target   string
    Paths    [][]string
    Found    bool
    JobID    int
    Duration time.Duration
}

// StartDFSMultipleWorkerPool memulai worker pool untuk pencarian multiple path DFS
func StartDFSMultipleWorkerPool(numWorkers int) (chan<- DFSMultipleJob, <-chan DFSMultipleResult) {
    jobs := make(chan DFSMultipleJob, numWorkers*2)
    results := make(chan DFSMultipleResult, numWorkers*2)
    var wg sync.WaitGroup
    
    // Mulai workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go dfsMultiplePathsWorker(i, jobs, results, &wg)
    }
    
    // Tutup channel results ketika semua worker selesai
    go func() {
        wg.Wait()
        close(results)
    }()
    
    return jobs, results
}

func dfsMulPath(element string, visited map[string]bool, trace []string, maxPaths int, pathCollection *[][]string) bool {
    if baseElements[element] {
        // Untuk elemen dasar, tambahkan path kosong
        *pathCollection = append(*pathCollection, []string{})
        return true
    }

    if visited[element] {
        return false
    }

    visited[element] = true
    recipes, ok := recipesMap[element]
    if !ok {
        return false
    }

    foundAnyPath := false
    elementTier := tierMap[element]
    
    // Cek setiap resep yang mungkin
    for _, ingr := range recipes {
        // Skip jika tier bahan >= tier elemen
        ingrTier1 := tierMap[ingr[0]]
        ingrTier2 := tierMap[ingr[1]]
        if ingrTier1 > elementTier || ingrTier2 > elementTier {
            continue
        }
        
        // Kumpulkan path untuk ingredient pertama
        leftPaths := [][]string{}
        newVisited1 := copyMap(visited)
        if dfsMulPath(ingr[0], newVisited1, append(trace, element), maxPaths, &leftPaths) {
            // Kumpulkan path untuk ingredient kedua
            rightPaths := [][]string{}
            newVisited2 := copyMap(visited)
            if dfsMulPath(ingr[1], newVisited2, append(trace, element), maxPaths, &rightPaths) {
                // Kombinasikan path dari kedua ingredient
                for _, left := range leftPaths {
                    for _, right := range rightPaths {
                        if len(*pathCollection) >= maxPaths {
                            break
                        }
                        combined := append(append([]string{}, left...), right...)
                        combined = append(combined, fmt.Sprintf("%s + %s = %s", ingr[0], ingr[1], element))
                        *pathCollection = append(*pathCollection, combined)
                        foundAnyPath = true
                    }
                    if len(*pathCollection) >= maxPaths {
                        break
                    }
                }
            }
        }
        
        if len(*pathCollection) >= maxPaths {
            break
        }
    }
    
    return foundAnyPath
}


// dfsMultiplePathsTopLevel mencari multiple path menggunakan DFS terparalel di level atas
func dfsMultiplePathsTopLevel(target string, maxPaths int) ([][]string, bool) {
    target = strings.ToLower(target)
    
    // Cek apakah target adalah elemen dasar
    if baseElements[target] {
        return [][]string{{}}, true
    }
    
    // Kumpulkan semua path
    visited := make(map[string]bool)
    allPaths := [][]string{}
    
    // Panggil fungsi dfsSinglePath yang telah dimodifikasi
    success := dfsMulPath(target, visited, []string{}, maxPaths, &allPaths)
    
    return allPaths, success && len(allPaths) > 0
}


// dfsMultiplePathsWorker memproses job untuk mencari multiple path menggunakan DFS
func dfsMultiplePathsWorker(id int, jobs <-chan DFSMultipleJob, results chan<- DFSMultipleResult, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("[Worker %d] Started for DFS multiple paths\n", id)
    
    for job := range jobs {
        startTime := time.Now()
        fmt.Printf("[Worker %d] Processing job %d: Finding paths to %s (max: %d)\n", 
                  id, job.JobID, job.Target, job.MaxPaths)
        
        // Cari multiple path menggunakan DFS terparalel di level atas
        paths, found := dfsMultiplePathsTopLevel(job.Target, job.MaxPaths)
        
        duration := time.Since(startTime)
        results <- DFSMultipleResult{
            Target:   job.Target,
            Paths:    paths,
            Found:    found,
            JobID:    job.JobID,
            Duration: duration,
        }
    }
}
