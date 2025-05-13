package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

func bfsMultiplePaths(target string, maxPaths int) ([][]string, bool,time.Duration, int) {
    target = strings.ToLower(target)
    start := time.Now()
    if baseElements[target] {
        return [][]string{{}}, true, 0, 0
    }

    type nodeInfo struct {
        predecessors [][2]string
        level        int
        sync.Mutex
    }

    var (
        elementInfo   = make(map[string]*nodeInfo)
        elementInfoMu sync.RWMutex // Mutex khusus untuk elementInfo
        queue         []string
        queueMu       sync.Mutex
        nodesVisited  int
        found         bool
        foundMu       sync.Mutex
    )

    // Inisialisasi base elements dengan lock
    elementInfoMu.Lock()
    for base := range baseElements {
        elementInfo[base] = &nodeInfo{level: 0}
        queue = append(queue, base)
        nodesVisited++
    }
    elementInfoMu.Unlock()

    // BFS level per level
    for len(queue) > 0 && !found {
        levelSize := len(queue)
        workCh := make(chan string, levelSize)
        
        for _, node := range queue {
            workCh <- node
        }
        close(workCh)
        
        queue = []string{}
        var wg sync.WaitGroup

        for i := 0; i < runtime.NumCPU(); i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for current := range workCh {
                    // Early termination check
                    foundMu.Lock()
                    if found {
                        foundMu.Unlock()
                        return
                    }
                    foundMu.Unlock()

                    // Dapatkan current level dengan lock
                    elementInfoMu.RLock()
                    currentInfo, exists := elementInfo[current]
                    elementInfoMu.RUnlock()

                    if !exists {
                        continue
                    }

                    currentInfo.Lock()
                    currentLevel := currentInfo.level
                    currentInfo.Unlock()

                    // Buat salinan keys untuk iterasi aman
                    elementInfoMu.RLock()
                    keys := make([]string, 0, len(elementInfo))
                    for k := range elementInfo {
                        keys = append(keys, k)
                    }
                    elementInfoMu.RUnlock()

                    // Proses kombinasi
                    for _, other := range keys {
                        mutex.RLock()
                        for resultElement, recipes := range recipesMap {
                            for _, ingredients := range recipes {
                                if (ingredients[0] == current && ingredients[1] == other) ||
                                    (ingredients[0] == other && ingredients[1] == current) {

                                    // Validasi tier
                                    if tierMap[ingredients[0]] >= tierMap[resultElement] || 
                                        tierMap[ingredients[1]] >= tierMap[resultElement] {
                                        continue
                                    }

                                    // Penggunaan lock untuk elementInfo
                                    elementInfoMu.Lock()
                                    if _, exists := elementInfo[resultElement]; !exists {
                                        elementInfo[resultElement] = &nodeInfo{
                                            level: currentLevel + 1,
                                        }
                                        queueMu.Lock()
                                        queue = append(queue, resultElement)
                                        queueMu.Unlock()
                                        nodesVisited++
                                    }

                                    resultInfo := elementInfo[resultElement]
                                    elementInfoMu.Unlock()

                                    resultInfo.Lock()
                                    if resultInfo.level == currentLevel+1 {
                                        resultInfo.predecessors = append(
                                            resultInfo.predecessors,
                                            [2]string{current, other},
                                        )
                                        
                                        // Update found dengan lock
                                        if resultElement == target && len(resultInfo.predecessors) >= maxPaths {
                                            foundMu.Lock()
                                            found = true
                                            foundMu.Unlock()
                                        }
                                    }
                                    resultInfo.Unlock()
                                }
                            }
                        }
                        mutex.RUnlock()
                    }
                }
            }()
        }
        wg.Wait()
    }

    // Helper untuk cek duplikasi path
	isPathExists := func(paths [][]string, newPath []string) bool {
		for _, p := range paths {
			if reflect.DeepEqual(p, newPath) {
				return true
			}
		}
		return false
	}
	// Rekonstruksi path dari predecessor
	var buildPaths func(element string) [][]string
	buildPaths = func(element string) [][]string {
		if elementInfo[element].level == 0 {
			return [][]string{{}}
		}

		var paths [][]string
		for _, predecessors := range elementInfo[element].predecessors {
			for _, p1 := range buildPaths(predecessors[0]) {
				for _, p2 := range buildPaths(predecessors[1]) {
					sorted1, sorted2 := normalizeIngredients(predecessors[0], predecessors[1])
					// Buat path baru
					newPath := make([]string, len(p1)+len(p2)+1)
					copy(newPath, p1)
					copy(newPath[len(p1):], p2)
					newPath[len(newPath)-1] = fmt.Sprintf("%s + %s = %s", sorted1, sorted2, element)
					// if sorted1 and sorted 2 already in
					// Cek duplikasi sebelum append
					if !isPathExists(paths, newPath) {
						paths = append(paths, newPath)
					}
                    if len(paths) >= maxPaths {
                        return paths
                    }
				}
			}
		}
		return paths
	}

	if _, exists := elementInfo[target]; !exists {
		return [][]string{}, false, time.Since(start), nodesVisited
	}

	allPaths := buildPaths(target)
	if len(allPaths) > maxPaths {
		allPaths = allPaths[:maxPaths]
	}

	return allPaths, true, time.Since(start), nodesVisited
}

// Tambahkan fungsi helper untuk mengurutkan nama bahan
func normalizeIngredients(a, b string) (string, string) {
	if a < b {
		return a, b
	}
	return b, a
}
