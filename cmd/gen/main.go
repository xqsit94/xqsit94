package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xqsit94/xqsit94/internal/github"
	"github.com/xqsit94/xqsit94/internal/posts"
	"github.com/xqsit94/xqsit94/internal/profile"
)

type TemplateData struct {
	Posts     []posts.Post
	ShowPosts bool
	Github    *github.Stats
	Profile   struct {
		Experience      []profile.Experience
		Education       []profile.Education
		TotalExperience string
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	showPosts := strings.EqualFold(os.Getenv("SHOW_LATEST_POSTS"), "true")

	var blogPosts []posts.Post
	if showPosts {
		var err error
		blogPosts, err = posts.GetPosts()
		if err != nil {
			log.Fatalf("Error getting blogPosts: %v", err)
		}
	}

	githubStats, err := github.GetGithubStats()
	if err != nil {
		log.Fatalf("Error getting GitHub stats: %v", err)
	}

	tmpl, err := template.ParseFiles("templates/README.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	output, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			log.Fatalf("Error closing output file: %v", err)
		}
	}(output)

	data := TemplateData{
		Posts:     blogPosts,
		ShowPosts: showPosts,
		Github:    githubStats,
		Profile: struct {
			Experience      []profile.Experience
			Education       []profile.Education
			TotalExperience string
		}{
			Experience:      profile.GetExperience(),
			Education:       profile.GetEducation(),
			TotalExperience: profile.CalculateTotalExperience(),
		},
	}
	if err := tmpl.Execute(output, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
