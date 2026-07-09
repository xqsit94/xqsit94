// Command cards regenerates the profile SVGs offline (no GitHub token needed),
// using the real profile data and representative stats. Useful for previewing
// design changes locally: `go run ./cmd/cards`.
package main

import (
	"log"

	"github.com/xqsit94/xqsit94/internal/card"
	"github.com/xqsit94/xqsit94/internal/profile"
)

func main() {
	// representative stats (the weekly `cmd/gen` run overwrites with live values)
	stats := card.Stats{Commits: 1429, PullRequests: 2555, Stars: 133, Issues: 608, Contributed: 24}

	d := card.Assemble("manikandan", "xqsit.dev",
		profile.GetExperience(), profile.CalculateTotalExperience(), stats)

	if err := card.WriteAll("assets", d); err != nil {
		log.Fatalf("write cards: %v", err)
	}
	log.Println("wrote assets/dark_mode.svg, assets/light_mode.svg")
}
