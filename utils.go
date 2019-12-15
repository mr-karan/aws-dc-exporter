package main

import "strings"

func replaceWithUnderscores(text string) string {
	replacer := strings.NewReplacer(" ", "_", ",", "_", "\t", "_", ",", "_", "/", "_", "\\", "_", ".", "_", "-", "_", ":", "_", "=", "_")
	return replacer.Replace(text)
}

func stateToFloat(state string) float64 {
	if state == "available" {
		return 1
	}
	return 0
}
