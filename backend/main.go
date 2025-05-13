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
	"time"

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
		numWorkers := runtime.NumCPU()
		if numberRecipeInt == 1 {
			if method == "bfs" {
				if bidirectional := c.Query("bidirectional"); bidirectional == "true" {
					steps, ok, runtime, nodesVisited := bfsBidirectionalPath(strings.ToLower(target))
					result := Result{
						Found:        ok,
						Steps:        steps,
						Runtime:      runtime,
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
						Runtime:      runtimes,
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
						Runtime:      runtime,
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
						Runtime:      runtime,
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

				bfsJobs, bfsResults := StartBFSMultipleWorkerPool(numWorkers)
				// Submit jobs
				go func() {
					bfsJobs <- BFSMultipleJob{
						Target:   target,
						MaxPaths: maxPathsPerTarget,
						JobID:    1,
					}
					close(bfsJobs)
				}()

				resultsJSON := make([]map[string][]string, 0)
				for result := range bfsResults {
					if result.Found {
						for i, path := range result.Paths {
							pathJSON := make(map[string][]string)
							pathJSON[fmt.Sprintf("Path %d", i+1)] = path
							resultsJSON = append(resultsJSON, pathJSON)
							resultsJSON[i]["Runtime"] = []string{result.Runtime.String()}
							resultsJSON[i]["NodesVisited"] = []string{strconv.Itoa(result.NodesVisited)}
						}
					}
				}
				//
				jsonResult, _ := json.Marshal(resultsJSON)
				c.Data(200, "application/json", jsonResult)
				return
			} else if method == "dfs" {
				dfsJobs, dfsResults := StartDFSMultipleWorkerPool(numWorkers)
				go func() {
					dfsJobs <- DFSMultipleJob{
						Target:   target,
						MaxPaths: maxPathsPerTarget,
						JobID:    3,
					}
					close(dfsJobs)
				}()
				resultsJSON := make([]map[string][]string, 0)
				for result := range dfsResults {
					if result.Found {
						var totalRuntime time.Duration = result.Runtime
						var totalNodes int = result.NodesVisited
						for i, path := range result.Paths {
							pathJSON := make(map[string][]string)
							pathJSON[fmt.Sprintf("Path %d", i+1)] = path
							resultsJSON = append(resultsJSON, pathJSON)
							resultsJSON[i]["Runtime"] = []string{totalRuntime.String()}
							resultsJSON[i]["NodesVisited"] = []string{strconv.Itoa(totalNodes)}
						}
					}
				}
				jsonResult, _ := json.Marshal(resultsJSON)
				c.Data(200, "application/json", jsonResult)
				return
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
