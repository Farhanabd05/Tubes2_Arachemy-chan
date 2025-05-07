package main

import "fmt"
func sliceToSet(slice []string) map[string]bool {
	set := map[string]bool{}
	for _,v := range slice {
		set[v] = true 
	}
	return set
}

func copySet(original  map[string]bool) map[string]bool {
	newSet := map[string]bool{}
	for k, v := range original {
		newSet[k] = v
	}
	return newSet
}

func keys(m map[string]bool) []string {
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func stateToKey(m map[string]bool) string {
	keys := keys(m)

	return fmt.Sprint(keys)
}