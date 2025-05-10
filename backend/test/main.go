package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type RecipeType struct {
	Element     string
	ImgUrl1     string
	ImgUrl2     string
	Ingredient1 string
	Ingredient2 string
	Type        int
}

func getElementType(index int) int {
	switch index {
	case 1:
		return 0
	case 2:
		// Skip (Ruins/Archeologist)
		return -1
	case 3:
		return 1
	case 4:
		return 2
	case 5:
		return 3
	case 6:
		return 4
	case 7:
		return 5
	case 8:
		return 6
	case 9:
		return 7
	case 10:
		return 8
	case 11:
		return 9
	case 12:
		return 10
	case 13:
		return 11
	case 14:
		return 12
	case 15:
		return 13
	case 16:
		return 14
	case 17:
		return 15
	default:
		return -1
	}
}

func ScrapeHandler(ctx *gin.Context) {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	var recipes []RecipeType

	c := colly.NewCollector(colly.AllowedDomains("little-alchemy.fandom.com"),
	// Add timeout settings to avoid long wait times
	colly.MaxDepth(1),
	colly.Async(true),
)
		// Limit concurrent requests
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})
	tableIndex := 0

	c.OnHTML("table.list-table", func(table *colly.HTMLElement) {
		tableIndex++
		elementType := getElementType(tableIndex)
		if elementType == -1 {
			return
		}

		table.ForEach("tbody tr", func(_ int, h *colly.HTMLElement) {
			element := strings.TrimSpace(h.ChildText("td:first-of-type a"))
			if element == "" || element == "Time" || element == "Ruins" || element == "Archeologist" {
				return
			}

			h.ForEach("td:nth-of-type(2) li", func(_ int, li *colly.HTMLElement) {
				aTags := li.DOM.Find("a")

				if aTags.Length() < 4 {
					return
				}

				imgUrl1, _ := aTags.Eq(0).Find("img").Attr("data-src")
				imgUrl2, _ := aTags.Eq(2).Find("img").Attr("data-src")
				ingredient1 := strings.TrimSpace(aTags.Eq(1).Text())
				ingredient2 := strings.TrimSpace(aTags.Eq(3).Text())

				if ingredient1 == "Time" || ingredient2 == "Time" || ingredient1 == "Ruins" || ingredient2 == "Ruins" || ingredient1 == "Archeologist" || ingredient2 == "Archeologist" {
					return
				}

				r := RecipeType{
					Element:     strings.ToLower(element),
					ImgUrl1:     imgUrl1,
					ImgUrl2:     imgUrl2,
					Ingredient1: strings.ToLower(ingredient1),
					Ingredient2: strings.ToLower(ingredient2),
					Type:        elementType,
				}
				recipes = append(recipes, r)
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Print("Visiting ", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error:", e.Error())
				// Check for specific network/TCP errors
		if strings.Contains(e.Error(), "dial tcp") || 
		   strings.Contains(e.Error(), "context deadline exceeded") ||
		   strings.Contains(e.Error(), "i/o timeout") {
			
			// Log specific TCP error message
			fmt.Printf("TCP connection failed: %s. Retrying...\n", e.Error())
			
			// You could implement retry logic here
			// For example, try an alternative domain or proxy
			
			// Optional: Retry with a different transport
			transport := &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   15 * time.Second,
				ResponseHeaderTimeout: 15 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}
			
			c.WithTransport(transport)
			
			// Optional: Try with a backup URL if available
			// c.Visit(backupUrl)
		}
	})

	err := c.Visit(url)
	if err != nil {
		if strings.Contains(err.Error(), "dial tcp") || 
		   strings.Contains(err.Error(), "context deadline exceeded") ||
		   strings.Contains(err.Error(), "i/o timeout") {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Network connection failed. Please check your internet connection and try again later.",
				"details": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// Wait for all requests to finish
	c.Wait()
	if err := os.MkdirAll("data", 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create data directory"})
		return
	}

	filePath := "data/recipes.json"
	jsonBytes, err := json.Marshal(recipes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal recipes to JSON"})
		return
	}
	if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recipes"})
		return
	}

	ctx.SetCookie("scraped", "true", 86400, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"data": recipes})
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/scrape", ScrapeHandler)
	r.Run(":8080")
}

