package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Tipe data untuk job pencarian BFS multiple path
type BFSMultipleJob struct {
    Target   string
    MaxPaths int
    JobID    int
}

// Tipe data untuk hasil pencarian
type BFSMultipleResult struct {
    Target   string
    Paths    [][]string
    Found    bool
    JobID    int
    Duration time.Duration
}

func bfsMultiplePaths(target string, maxPaths int) ([][]string, bool) {
    target = strings.ToLower(target)
    
    // Cek apakah target adalah elemen dasar
    if baseElements[target] {
        return [][]string{{}}, true
    }
    
    // Untuk menyimpan multiple path ke setiap elemen
    pathsToElement := make(map[string][][]string)
    discovered := make(map[string]bool)
    
    // Inisialisasi dengan elemen dasar
    for base := range baseElements {
        discovered[base] = true
        pathsToElement[base] = [][]string{{}} // Empty path untuk elemen dasar
    }
    
    // Implementasi BFS
    iteration := 0
    
    for {
        iteration++
        newDiscoveries := false
        mutex.RLock() // Lock untuk akses ke recipe maps
        
        for resultElement, recipes := range recipesMap {
            // Skip elemen yang sudah diproses
            if _, ok := pathsToElement[resultElement]; ok {
                continue
            }
            
            resultTier := tierMap[resultElement]
            pathsForElement := [][]string{}
            
            for _, ingredients := range recipes {
                ingr1 := ingredients[0]
                ingr2 := ingredients[1]
                
                // Cek apakah kedua ingredients sudah ditemukan
                if paths1, ok1 := pathsToElement[ingr1]; ok1 {
                    if paths2, ok2 := pathsToElement[ingr2]; ok2 {
                        // Cek constraint tier
                        if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
                            continue
                        }
                        
                        // Kombinasikan path dari kedua ingredients
                        for _, p1 := range paths1 {
                            for _, p2 := range paths2 {
                                newPath := make([]string, len(p1)+len(p2)+1)
                                copy(newPath, p1)
                                copy(newPath[len(p1):], p2)
                                
                                // Tambahkan langkah kombinasi
                                newPath[len(newPath)-1] = fmt.Sprintf("%s + %s = %s", ingr1, ingr2, resultElement)
                                
                                // Tambahkan path ke hasil
                                pathsForElement = append(pathsForElement, newPath)
                                
                                // Batasi jumlah path per elemen
                                if len(pathsForElement) >= maxPaths {
                                    break
                                }
                            }
                            
                            if len(pathsForElement) >= maxPaths {
                                break
                            }
                        }
                    }
                }
            }
            
            // Jika menemukan path untuk elemen ini, catat
            if len(pathsForElement) > 0 {
                pathsToElement[resultElement] = pathsForElement
                discovered[resultElement] = true
                newDiscoveries = true
                
                // Cek apakah ini target
                if resultElement == target {
                    mutex.RUnlock()
                    return pathsForElement[:min(len(pathsForElement), maxPaths)], true
                }
            }
        }
        
        mutex.RUnlock()
        
        // Jika tidak ada penemuan baru, hentikan loop
        if !newDiscoveries {
            break
        }
    }
    
    // Return path jika target ditemukan
    if paths, ok := pathsToElement[target]; ok {
        return paths[:min(len(paths), maxPaths)], true
    }
    
    return [][]string{}, false
}

// Worker yang memproses job BFS untuk multiple path
func bfsMultiplePathsWorker(id int, jobs <-chan BFSMultipleJob, results chan<- BFSMultipleResult, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("[Worker %d] Started for BFS multiple paths\n", id)
    
    for job := range jobs {
        startTime := time.Now()
        fmt.Printf("[Worker %d] Processing job %d: Finding paths to %s (max: %d)\n", 
                  id, job.JobID, job.Target, job.MaxPaths)
        
        // Pencarian multiple paths menggunakan BFS
        paths, found := bfsMultiplePaths(job.Target, job.MaxPaths)
        
        duration := time.Since(startTime)
        
        results <- BFSMultipleResult{
            Target:   job.Target,
            Paths:    paths,
            Found:    found,
            JobID:    job.JobID,
            Duration: duration,
        }
    }
}

func StartBFSMultipleWorkerPool(numWorkers int) (chan<- BFSMultipleJob, <-chan BFSMultipleResult) {
    jobs := make(chan BFSMultipleJob, numWorkers*2)     // Buffer untuk jobs
    results := make(chan BFSMultipleResult, numWorkers*2) // Buffer untuk results
    
    var wg sync.WaitGroup
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go bfsMultiplePathsWorker(i, jobs, results, &wg)
    }
    
    // Tutup channel results ketika semua worker selesai
    go func() {
        wg.Wait()
        close(results)
    }()
    
    return jobs, results
}