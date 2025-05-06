// backend/main.go
package main

import (
    "encoding/json"
    "net/http"
    "log"
)

type Message struct {
    Text string `json:"text"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Message{Text: "Halo dari Golang!"})
}


func main() {
    http.HandleFunc("/api/hello", HelloHandler)
    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
