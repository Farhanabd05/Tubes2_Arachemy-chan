// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Recipe represents a recipe with inputs and output
type Recipe struct {
	Input  []string `json:"input"`
	Output string   `json:"output"`
}

// Node represents a graph node
type Node struct {
	Value string
}

// Graph represents our graph structure
type Graph struct {
	Nodes       map[string]*Node
	Connections map[string][]string
	mu          sync.RWMutex
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	return &Graph{
		Nodes:       make(map[string]*Node),
		Connections: make(map[string][]string),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(value string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if _, exists := g.Nodes[value]; !exists {
		g.Nodes[value] = &Node{Value: value}
	}
}

// AddConnection adds a connection from inputs to output
func (g *Graph) AddConnection(inputs []string, output string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	// Add nodes if they don't exist
	g.AddNode(output)
	for _, input := range inputs {
		g.AddNode(input)
		
		// Create connection from input to output
		if _, exists := g.Connections[input]; !exists {
			g.Connections[input] = []string{}
		}
		g.Connections[input] = append(g.Connections[input], output)
	}
}

// ProcessRecipes processes recipes and builds the graph
func (g *Graph) ProcessRecipes(recipes []Recipe) {
	for _, recipe := range recipes {
		for _, input := range recipe.Input {
			g.AddNode(input)
		}
		g.AddNode(recipe.Output)
		g.AddConnection(recipe.Input, recipe.Output)
	}
}

// BFS performs a breadth-first search to find a path from start to target
func (g *Graph) BFS(start, target string) ([]string, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	// Check if start and target nodes exist
	if _, exists := g.Nodes[start]; !exists {
		return nil, false
	}
	if _, exists := g.Nodes[target]; !exists {
		return nil, false
	}
	
	// BFS implementation
	queue := []string{start}
	visited := make(map[string]bool)
	parent := make(map[string]string)
	
	visited[start] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		// If we found the target, reconstruct the path
		if current == target {
			path := []string{}
			for current != start {
				path = append([]string{current}, path...)
				current = parent[current]
			}
			path = append([]string{start}, path...)
			return path, true
		}
		
		// Visit all neighbors
		for _, neighbor := range g.Connections[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}
	
	return nil, false
}

// Handler for processing recipes
func handleProcessRecipes(g *Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var recipes []Recipe
		err := json.NewDecoder(r.Body).Decode(&recipes)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		
		g.ProcessRecipes(recipes)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// Handler for BFS search
func handleBFS(g *Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		start := r.URL.Query().Get("start")
		target := r.URL.Query().Get("target")
		
		if start == "" || target == "" {
			http.Error(w, "Start and target parameters are required", http.StatusBadRequest)
			return
		}
		
		path, found := g.BFS(start, target)
		
		w.Header().Set("Content-Type", "application/json")
		if found {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "success",
				"path":   path,
				"found":  true,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "success",
				"found":  false,
			})
		}
	}
}

// Handler for getting all nodes and connections
func handleGetGraph(g *Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		g.mu.RLock()
		defer g.mu.RUnlock()
		
		// Extract node values
		nodes := make([]string, 0, len(g.Nodes))
		for nodeKey := range g.Nodes {
			nodes = append(nodes, nodeKey)
		}
		
		// Extract connections
		connections := make(map[string][]string)
		for source, targets := range g.Connections {
			connections[source] = targets
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "success",
			"nodes":       nodes,
			"connections": connections,
		})
	}
}

// Enable CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create a new graph
	graph := NewGraph()
	
	// Set up routes
	mux := http.NewServeMux()
	
	// API endpoints
	mux.Handle("/api/recipes", corsMiddleware(handleProcessRecipes(graph)))
	mux.Handle("/api/search", corsMiddleware(handleBFS(graph)))
	mux.Handle("/api/graph", corsMiddleware(handleGetGraph(graph)))
	
	// Start server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}