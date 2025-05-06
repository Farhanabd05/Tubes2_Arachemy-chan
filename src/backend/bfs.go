package main

import "fmt"

type treeNode struct {
	value    int
	children []*treeNode
}

func (n *treeNode) addChild(val *treeNode) {
	n.children = append(n.children, val)	
}

func traversePreOrder(node * treeNode){
	if node == nil {
		return
	}
	fmt.Print(node.value, " ")
	for _, child := range node.children {
		traversePreOrder(child)
	}
}

func printTree(node * treeNode, prefix string, isLast bool) {
	if node == nil {
		return
	}
	fmt.Print(prefix)
    if isLast {
        fmt.Print("└── ")
        prefix += "    "
    } else {
        fmt.Print("├── ")
        prefix += "│   "
    }
	fmt.Println(node.value)

    for i, child := range node.children {
        printTree(child, prefix, i == len(node.children)-1)
    }
}
// declare a list
func bfs(root *treeNode, value int) {
	if root == nil {
		fmt.Println("Tree is empty")
		return
	}
	queue := []*treeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:] // dequeue
		if node.value == value { // cek nilai
			fmt.Println("Found the value in the tree")
			fmt.Println("Depth: ", len(queue))
			return
		}
		// tambahkan anak-anak ke antrian
		for _, child := range node.children {
			queue = append(queue, child)
		}
	}
	fmt.Println("Value not found in the tree")
}
// main func
func main() {

	// Buat root node
    root := &treeNode{value: 1}

    // Tambah anak-anak untuk root
    child2 := &treeNode{value: 2}
    child3 := &treeNode{value: 3}
    child4 := &treeNode{value: 4}
    root.addChild(child2)
    root.addChild(child3)
    root.addChild(child4)

    // Tambah anak-anak untuk node 2
    child2.addChild(&treeNode{value: 5})
    child2.addChild(&treeNode{value: 6})

    // Tambah anak-anak untuk node 3
    child3.addChild(&treeNode{value: 7})

    // Tambah anak-anak untuk node 4
    child4.addChild(&treeNode{value: 8})
    child4.addChild(&treeNode{value: 9})
    child4.addChild(&treeNode{value: 10})

    // Visualisasi pohon
    fmt.Println("Visualisasi pohon:")
    printTree(root, "", true)

    // Traversal preorder
    fmt.Print("Preorder traversal: ")
    traversePreOrder(root)
    fmt.Println()

	bfs(root, 8);
}
