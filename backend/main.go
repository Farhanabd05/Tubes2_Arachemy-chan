package main

import (
	// "context"
	// "log"
	// "os"
	// "time"
	"encoding/json"
	"strings"

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
	r.GET("/singlepath", func(c *gin.Context) {
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
		recipes, err := loadRecipes("test/data/recipes.json")
		if err != nil {
			c.JSON(500, gin.H{"error": "Error loading recipes: " + err.Error()})
			return
		}
		buildRecipeMap(recipes)
		if method == "bfs" {
			steps, ok := bfsSinglePath(strings.ToLower(target))
			result := Result{
				Found: ok,
				Steps: steps,
			}
			jsonResult, _ := json.Marshal(result)
			c.Data(200, "application/json", jsonResult)
			return
		} else if (method == "dfs") {
			steps, ok := dfsSinglePath(strings.ToLower(target), map[string]bool{}, []string{})
			result := Result{
				Found: ok,
				Steps: steps,
			}
			jsonResult, _ := json.Marshal(result)
			c.Data(200, "application/json", jsonResult)
			return
		}
	})
	r.Run(":8080")
}
