package main

import (
	"fmt"
	"sort"
	"sync"
)

type Job struct {
	JobID int
	JobType string // mul, sum, dll
	Target string
}

type JobResult struct {
	JobID  int
	WorkerID int
	JobType string
	Target string
	Found bool
	Steps []string
	Err error
}

func main() {
	var wg sync.WaitGroup
	workerCount := 4
	jobCount := 10
	jobChan := make(chan Job, jobCount) // channel lokal, tidak akan pernah diisi
	resChan := make(chan JobResult) // channel lokal, tidak akan pernah diisi
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("[Worker %d] Started\n", workerID)
			for job := range jobChan {
				fmt.Printf("[Worker %d] Processing job: %+v\n", workerID, job)
				var steps []string
				var found bool
				var err error
				switch job.JobType {
				case "bfs":
					defer func() {
						if r := recover(); r!=nil {
							err = fmt.Errorf("panic: %v", r)
						}
					}()
					steps, found = bfsSinglePath(job.Target)
				default:
					err = fmt.Errorf("unsupported  job type: %s", job.JobType)
				}
				resChan <- JobResult{
					JobID:  job.JobID,
					WorkerID: workerID,
					JobType: job.JobType,
					Target: job.Target,
					Found:   found,
					Steps:   steps,
					Err:     err,
				}
			}

			fmt.Printf("[Worker %d] Exiting\n", workerID)
		}(i)
	}

	targets := []string{"life", "human", "zombie", "energy"}
	// Kirim job dari main thread
	go func() {
		for i, target := range targets {
			fmt.Printf("[Main] Submitting job %d: %s\n", i+1, target)
			jobChan <- Job{
				JobID:   i + 1,
				JobType: "bfs",
				Target:  target,
			}
		}
		close(jobChan)
	}()

	//ambil hasil dari resChan
	resMap := make(map[int]JobResult)
	go func() {
		wg.Wait()
		close(resChan)
	}()

	// kumpulkan hasil
	for res := range resChan {
		if res.Err != nil {
			fmt.Printf("[Result] Job %d (target: %s) FAILED by worker %d: %v\n", res.JobID, res.Target, res.WorkerID, res.Err)
			continue
		}
		fmt.Printf("[Result] Job %d (target: %s) by worker %d: Found = %v\n", res.JobID, res.Target, res.WorkerID, res.Found)
		if res.Found {
			for _, step := range res.Steps {
				fmt.Printf("  %s\n", step)
			}
		}
	}
	// cetak hasil terurut
	fmt.Println("[Main] Sorted Results:")
	keys := make([]int, 0, len(resMap))
	for k := range resMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("  Job %d → %d\n (by worker %d)\n", k, resMap[k].Target, resMap[k].WorkerID)
	}

	fmt.Println("[Main] All workers done.")
}
