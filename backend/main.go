package main

import (
	// "context"
	// "log"
	// "os"
	// "time"
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
	r.GET("/search", func(c * gin.Context) {
		el1 := c.Query(("e1"))
		el2 := c.Query(("e2"))

		el1 = strings.ToLower((el1))
		el2 = strings.ToLower((el2))

		result, ok := getCombination(el1, el2)
		// if !ok {
		// 	result, ok := comb[[2]string{el2, el1}]
		// }
		if ok {
			c.JSON(200, gin.H{
				"found": true,
				"result" : result,
			})
		} else {
			c.JSON(200, gin.H{
				"found" : false,
			})
		}
	}) 

	r.GET("/find", func(c *gin.Context) {
		target := strings.ToLower(c.Query("target"))
		path, found := bfsFindTarget(target)
		if found {
			c.JSON(200, gin.H{
				"found": true,
				"path":path,
			})
		} else {
			c.JSON(200, gin.H{
				"found":false,
			})
		}
	})
	r.Run(":8080")
}
