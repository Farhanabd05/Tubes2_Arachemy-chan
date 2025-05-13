package main

import (
	// "context"
	// "log"
	// "os"
	// "time"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/go-redis/redis/v8"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/scrape", ScrapeHandler)
	r.GET("/find", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(400, gin.H{"error": "Target tidak boleh kosong"})
			return
		}
		method := c.Query("method")
		if method == "" {
			c.JSON(400, gin.H{"error": "Method tidak boleh kosong"})
			return
		}

		if method != "bfs" && method != "dfs" {
			c.JSON(400, gin.H{"error": "Method tidak valid"})
			return
		}

		numberRecipe := c.Query("numberRecipe")
		if numberRecipe == "" {
			c.JSON(400, gin.H{"error": "Number recipe tidak boleh kosong"})
			return
		}

		recipes, err := loadRecipes("data/recipes.json")
		if err != nil {
			c.JSON(500, gin.H{"error": "Error loading recipes: " + err.Error()})
			return
		}
		buildRecipeMap(recipes)

		numberRecipeInt, err := strconv.Atoi(numberRecipe)
		var maxPathsPerTarget = numberRecipeInt
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid numberRecipe value"})
			return
		}
		if numberRecipeInt == 1 {
			if method == "bfs" {
				if bidirectional := c.Query("bidirectional"); bidirectional == "true" {
					steps, ok, runtime, nodesVisited := bfsBidirectionalPath(strings.ToLower(target))
					result := Result{
						Found:        ok,
						Steps:        steps,
						Runtime:      runtime.String(),
						NodesVisited: nodesVisited,
					}
					jsonResult, _ := json.Marshal(result)
					c.Data(200, "application/json", jsonResult)
					return
				} else {
					steps, ok, runtimes, nodes := bfsSinglePath(strings.ToLower(target))
					result := Result{
						Found:        ok,
						Steps:        steps,
						Runtime:      runtimes.String(),
						NodesVisited: nodes,
					}
					jsonResult, _ := json.Marshal(result)
					c.Data(200, "application/json", jsonResult)
					return
				}
			} else if method == "dfs" {
				if bidirectional := c.Query("bidirectional"); bidirectional == "true" {
					steps, ok, runtime, nodesVisited := dfsBidirectionalPath(strings.ToLower(target))
					result := Result{
						Found:        ok,
						Steps:        steps,
						Runtime:      runtime.String(),
						NodesVisited: nodesVisited,
					}
					jsonResult, _ := json.Marshal(result)
					c.Data(200, "application/json", jsonResult)
					return
				} else {
					steps, ok, runtime, nodesVisited := DFSWrapper(strings.ToLower(target))
					result := Result{
						Found:        ok,
						Steps:        steps,
						Runtime:      runtime.String(),
						NodesVisited: nodesVisited,
					}
					jsonResult, _ := json.Marshal(result)
					c.Data(200, "application/json", jsonResult)
					return
				}
			}
		} else {
			if method == "bfs" {
				// Buat worker pool untuk BFS multiple paths

				bfsResults, found, runtime, nodes := bfsMultiplePaths(target, maxPathsPerTarget)

				resultsJSON := make([]map[string][]string, 0)
				if found {
					for i, path := range bfsResults {
						pathJSON := make(map[string][]string)
						pathJSON[fmt.Sprintf("Path %d", i+1)] = path
						resultsJSON = append(resultsJSON, pathJSON)
						resultsJSON[i]["Runtime"] = []string{runtime.String()}
						resultsJSON[i]["NodesVisited"] = []string{strconv.Itoa(nodes)}
					}
				}
				//
				jsonResult, _ := json.Marshal(resultsJSON)
				c.Data(200, "application/json", jsonResult)
				return
			} else if method == "dfs" {
				jobs := make(chan Job)
				results := make(chan JobResultDFS)
				var wg sync.WaitGroup
				maxResults = numberRecipeInt

				numWorkers := runtime.NumCPU()
				for i := 0; i < numWorkers; i++ {
					wg.Add(1)
					go worker(i, jobs, results, &wg)
				}

				go func() {
					for i := range recipesMap[target] {
						jobs <- Job{JobID: i + 1, JobType: "dfs", Target: target}
					}
					close(jobs)
				}()

				go func() {
					wg.Wait()
					close(results)
				}()

				printed := 0
				seenPaths := make(map[string]bool)
				var resultsJSON []map[string][]string // Ubah dari []Result ke []map[string][]string
				for res := range results {
					key := strings.Join(res.Steps, "|")
					if !seenPaths[key] {
						seenPaths[key] = true
						pathJSON := make(map[string][]string)
						pathJSON[fmt.Sprintf("Path %d", printed+1)] = res.Steps
						pathJSON["Runtime"] = []string{res.Duration.String()}
						pathJSON["NodesVisited"] = []string{strconv.Itoa(len(res.Steps))}
						resultsJSON = append(resultsJSON, pathJSON)
						printed++
						if printed >= maxResults {
							break
						}
					}
				}
				jsonResult, _ := json.Marshal(resultsJSON)
				c.Data(200, "application/json", jsonResult)
			}
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on 0.0.0.0:%s\n", port)
	log.Fatal(r.Run("0.0.0.0:" + port))
}
