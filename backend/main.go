package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// worker function that processes jobs from the job channel
func worker(workerID int, jobChan <-chan Job, resChan chan<- JobResult, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Worker %d] Started\n", workerID)
	
	for job := range jobChan {
		fmt.Printf("[Worker %d] Processing job: %+v\n", workerID, job)
		startTime := time.Now()
		
		var steps []string
		var found bool
		var err error
		
		// Use defer and recover to catch panics within this job
		func() {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					err = fmt.Errorf("panic: %v\nStack trace:\n%s", r, stack)
					fmt.Printf("[Worker %d] PANIC in job %d: %v\n%s\n", 
						workerID, job.JobID, r, stack)
				}
			}()
			
			switch job.JobType {
			case "bfs":
				steps, found = bfsSinglePath(job.Target)
			default:
				err = fmt.Errorf("unsupported job type: %s", job.JobType)
				fmt.Printf("[Worker %d] ERROR: %v\n", workerID, err)
			}
		}()
		
		duration := time.Since(startTime)
		
		// Send result back through result channel
		resChan <- JobResult{
			JobID:    job.JobID,
			WorkerID: workerID,
			JobType:  job.JobType,
			Target:   job.Target,
			Found:    found,
			Steps:    steps,
			Err:      err,
			Duration: duration,
		}
		
		fmt.Printf("[Worker %d] Completed job %d in %v\n", 
			workerID, job.JobID, duration)
	}
	
	fmt.Printf("[Worker %d] Exiting\n", workerID)
}

func main() {
	fmt.Println("[Main] Starting Little Alchemy Pathfinder")
	
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
	workerCount := 4
	fmt.Printf("[Main] Starting %d workers\n", workerCount)
	
	jobChan := make(chan Job, workerCount*2) // Buffered channel for jobs
	resChan := make(chan JobResult, workerCount*2) // Buffered channel for results
	
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, jobChan, resChan, &wg)
	}
	
	// Send jobs
	targets := []string{"life", "human", "zombie", "energy", "metal", "plant"}
	go func() {
		for i, target := range targets {
			fmt.Printf("[Main] Sending job %d: bfs → %s\n", i+1, target)
			jobChan <- Job{
				JobID:   i + 1,
				JobType: "bfs",
				Target:  target,
			}
		}
		fmt.Println("[Main] All jobs sent, closing job channel")
		close(jobChan)
	}()
	
	// Wait for all workers to finish and close result channel
	go func() {
		wg.Wait()
		fmt.Println("[Main] All workers finished, closing result channel")
		close(resChan)
	}()
	
	// Process results
	fmt.Println("[Main] Processing results")
	resMap := make(map[int]JobResult)
	successCount := 0
	failCount := 0
	
	for res := range resChan {
		if res.Err != nil {
			fmt.Printf("[Result] Job %d (target: %s) FAILED by worker %d: %v\n", 
				res.JobID, res.Target, res.WorkerID, res.Err)
			failCount++
			continue
		}
		
		if res.Found {
			fmt.Printf("[Result] Job %d (target: %s) SUCCEEDED by worker %d in %v\n", 
				res.JobID, res.Target, res.WorkerID, res.Duration)
			for i, step := range res.Steps {
				fmt.Printf("  %d. %s\n", i+1, step)
			}
			successCount++
		} else {
			fmt.Printf("[Result] Job %d (target: %s) by worker %d: Not found (took %v)\n", 
				res.JobID, res.Target, res.WorkerID, res.Duration)
			failCount++
		}
		
		resMap[res.JobID] = res
	}
	
	fmt.Printf("[Main] All jobs processed. Success: %d, Failed/Not Found: %d\n", 
		successCount, failCount)
	fmt.Println("[Main] Program completed.")
}