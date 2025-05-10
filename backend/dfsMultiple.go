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

// dfsMultiplePathsTopLevel mencari multiple path menggunakan DFS terparalel di level atas
func dfsMultiplePathsTopLevel(target string, maxPaths int) ([][]string, bool) {
    target = strings.ToLower(target)
    
    // Cek apakah target adalah elemen dasar
    if baseElements[target] {
        return [][]string{{}}, true
    }
    
    // Dapatkan resep untuk target
    mutex.RLock()
    recipes, ok := recipesMap[target]
    elementTier := tierMap[target]
    mutex.RUnlock()
    
    if !ok {
        return [][]string{}, false
    }
    
    // Buat channel untuk mengumpulkan path
    pathChan := make(chan []string, maxPaths*2)
    
    // Buat wait group untuk goroutines
    var wg sync.WaitGroup
    
    // Buat mutex untuk melindungi pathCount
    var pathCountMutex sync.Mutex
    pathCount := 0
    
    // Filter resep valid
    validRecipes := [][]string{}
    for _, ingr := range recipes {
        // Skip jika tier ingredient >= tier elemen
        mutex.RLock()
        ingrTier1 := tierMap[ingr[0]]
        ingrTier2 := tierMap[ingr[1]]
        mutex.RUnlock()
        
        if ingrTier1 > elementTier || ingrTier2 > elementTier {
            continue
        }
        
        validRecipes = append(validRecipes, ingr)
    }
    
    // Mulai goroutine untuk setiap resep valid
    for _, recipe := range validRecipes {
        wg.Add(1)
        go func(recipe []string) {
            defer wg.Done()
            
            // Reset print counter
            printCount = 0
            
            // Coba mencari path untuk kedua ingredient
            visited1 := make(map[string]bool)
            path1, found1 := dfsSinglePath(recipe[0], visited1, []string{})
            
            if !found1 {
                return
            }
            
            visited2 := make(map[string]bool)
            path2, found2 := dfsSinglePath(recipe[1], visited2, []string{})
            
            if !found2 {
                return
            }
            
            // Gabungkan path
            combinedPath := append(path1, path2...)
            combinedPath = append(combinedPath, fmt.Sprintf("%s + %s = %s", recipe[0], recipe[1], target))
            
            // Cek apakah sudah cukup paths
            pathCountMutex.Lock()
            if pathCount < maxPaths {
                // Kirim path ke channel
                pathChan <- combinedPath
                pathCount++
            }
            pathCountMutex.Unlock()
        }(recipe)
    }
    
    // Tunggu semua goroutine selesai dan tutup path channel
    go func() {
        wg.Wait()
        close(pathChan)
    }()
    
    // Kumpulkan path
    var paths [][]string
    for path := range pathChan {
        paths = append(paths, path)
        if len(paths) >= maxPaths {
            break
        }
    }
    
    return paths, len(paths) > 0
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
