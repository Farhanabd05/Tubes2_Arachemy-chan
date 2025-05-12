// package main

// import (
// 	"fmt"
// 	"strings"
// 	"sync"
// )

// type TreeNode struct {
// 	Value    string
// 	Children []*TreeNode
// }

// // Fungsi DFS yang menelusuri setiap node dan mengumpulkan path menuju leaf node
// func dfs(node *TreeNode, currentPath []string, result *[]string, wg *sync.WaitGroup, workerID int) {
// 	defer wg.Done() // Menandakan goroutine selesai

// 	// Menambahkan node saat ini ke dalam path dengan informasi worker yang mengerjakan
// 	currentPath = append(currentPath, fmt.Sprintf("%s (worker %d)", node.Value, workerID))

// 	// Jika node ini adalah leaf (tidak punya anak), tambahkan path ke result
// 	if len(node.Children) == 0 {
// 		// Gabungkan path menjadi string
// 		path := strings.Join(currentPath, " -> ")
// 		// Menambahkan path beserta jumlah node yang dikunjungi
// 		*result = append(*result, fmt.Sprintf("%s (Total nodes: %d)", path, len(currentPath)))
// 		return
// 	}

// 	// Jika node ini memiliki anak, lanjutkan untuk menelusuri anak-anaknya
// 	for i, child := range node.Children {
// 		wg.Add(1) // Menambah jumlah goroutine yang harus ditunggu
// 		go dfs(child, currentPath, result, wg, workerID+i+1) // Worker ID bertambah untuk setiap cabang baru
// 	}
// }

// func main() {
// 	// Membuat tree sesuai permintaan
// 	root := &TreeNode{
// 		Value: "Root",
// 		Children: []*TreeNode{
// 			{Value: "Child 1", Children: []*TreeNode{
// 				{Value: "Child 1.1", Children: []*TreeNode{
// 					{Value: "Child 1.1.1"}}},
// 				{Value: "Child 1.2"},
// 			}},
// 			{Value: "Child 2", Children: []*TreeNode{
// 				{Value: "Child 2.1", Children: []*TreeNode{
// 					{Value: "Child 2.1.1", Children: []*TreeNode{
// 						{Value: "Child 2.1.1.1", Children: []*TreeNode{
// 							{Value: "Child 2.1.1.1.1", Children: []*TreeNode{
// 								{Value: "Child 2.1.1.1.1.1"},
// 								{Value: "Child 2.1.1.1.1.2"},
// 							}},
// 						}}}},
// 				}}},
// 			},
// 		},
// 	}

// 	var result []string
// 	var wg sync.WaitGroup

// 	// Menambahkan goroutine pertama untuk memulai DFS dari root
// 	wg.Add(1)
// 	go dfs(root, []string{}, &result, &wg, 1)

// 	// Menunggu semua pekerjaan selesai
// 	wg.Wait()

// 	// Menampilkan hasil yang dikumpulkan
// 	fmt.Println("Paths to leaf nodes:")
// 	for _, path := range result {
// 		fmt.Println(path)
// 	}
// }
