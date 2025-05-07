package main


// import "strings"
import (
	"fmt"
	"strings"
)
var combinations = map[[2]string]string{
	normalize("fire", "air"):    "energy",
	normalize("water", "earth"): "mud",
	normalize("energy", "mud"):  "life",
	normalize("life", "earth"):  "human",
	normalize("air", "water"):   "rain",
	normalize("earth", "rain"):  "plant",
	normalize("plant", "human"): "farmer",
}

// Normalisasi pasangan agar urutan tidak masalah
func normalize(a, b string) [2]string {
	if a < b {
		return [2]string{a, b}
	}
	return [2]string{b, a}
}

func getCombination(element1, element2 string) (string, bool) {
	key := normalize(element1, element2)
	result, exists := combinations[key]
	return result, exists
}

func bfsFindTarget(target string) ([]string, bool) {
	type State struct {
		Elements []string
		Path  []string
	}

	initial := []string{"fire", "water", "earth", "air"}
	visited := map[string]bool{}
	seenComb := map[string]bool{}
	queue := []State {
		{Known: initial, Path: []string{}},
	}

	fmt.Println("🔍 Mencari kombinasi menuju:", target)
	fmt.Println("🌱 Elemen dasar:", initial)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		fmt.Println("🔄 State sekarang:", curr.Known)
		for i := 0; i < len(curr.Known); i++ {
			for j := i + 1; j < len(curr.Known); j++ {
				a := curr.Known[i]
				b := curr.Known[j]
				key := normalize(a, b)
				keyStr := key[0] + "+" + key[1]

				if seenComb[keyStr] {
					continue
				}
				seenComb[keyStr] = true

				res, ok := getCombination(a, b)
				fmt.Printf("🔧 Mencoba: %s + %s → ", a, b)
				if ok {
					fmt.Printf("%s\n", res)
				} else {
					fmt.Println("tidak valid")
					continue
				}

				if visited[res] {
					fmt.Println("⛔ Sudah pernah ditemukan:", res)
					continue
				}

				if strings.EqualFold(res, target) {
					fmt.Println("🎯 Target ditemukan:", res)
					return append(curr.Path, fmt.Sprintf("%s + %s = %s", a, b, res)), true
				}

				visited[res] = true
				newKnown := append([]string{}, curr.Known...)
				newKnown = append(newKnown, res)

				newPath := append([]string{}, curr.Path...)
				newPath = append(newPath, fmt.Sprintf("%s + %s = %s", a, b, res))

				queue = append(queue, State{
					Known: newKnown,
					Path:  newPath,
				})
			}
		}
	}
	fmt.Println("❌ Tidak ditemukan jalur menuju target:", target)
	return nil, false
}
