package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ElementData struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Type        int    `json:"Type"`
}

func normalizeKey(s string) string {
	// 1) Lowercase
	// 2) Ganti spasi dengan underscore
	return strings.ReplaceAll(strings.ToLower(s), " ", "_")
}

func main() {
	// 1. Baca JSON input
	raw, err := os.ReadFile("recipes.json")
	if err != nil {
		panic(fmt.Errorf("gagal membaca JSON: %w", err))
	}
	var elems []ElementData
	if err := json.Unmarshal(raw, &elems); err != nil {
		panic(fmt.Errorf("gagal parse JSON: %w", err))
	}

	// 2. Scan folder images untuk file *_2.svg
	imageMap := make(map[string]string)
	err = filepath.Walk("images", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_2.svg") {
			// strip suffix "_2.svg" lalu lowercase key
			base := strings.TrimSuffix(info.Name(), "_2.svg")
			key := strings.ToLower(base)
			imageMap[key] = info.Name()
		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("gagal scanning folder images: %w", err))
	}

	// 3. Buat hasil mapping per entry
	type Mapped struct {
		Element           string `json:"Element"`
		ElementImage      string `json:"ElementImage"`
	}

	var out []Mapped
	for _, e := range elems {
		look := func(name string) string {
			k := normalizeKey(name)
			if img, ok := imageMap[k]; ok {
				return img
			}
			return "" // atau "NOT_FOUND"
		}

		out = append(out, Mapped{
			Element:          e.Element,
			ElementImage:     look(e.Element),
		})
	}

	// 4. Tulis hasil ke JSON baru
	result, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(fmt.Errorf("gagal marshal output: %w", err))
	}
	if err := os.WriteFile("mapped_elements.json", result, 0644); err != nil {
		panic(fmt.Errorf("gagal tulis output file: %w", err))
	}

	fmt.Println("Mapping selesai! Hasil tersimpan di mapped_elements.json")
}
