package main

import "github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"

func main() {
	err := cs2gsi.StartupAndServe(":8000")
	if err != nil {
		return
	}
}
