package main

import (
	"fmt"
	"log"

	"github.com/IjichiNijika99/go-gitblog/internal/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("Init failed: %v", err)
	}

	fmt.Println("Init success")
	fmt.Printf("Repo Name: %s\n", cfg.RepoName)
	fmt.Printf("Backup Dir: %s\n", cfg.BackupDir)
	if cfg.IssueNumber != 0 {
		fmt.Printf("Issue number: %d\n", cfg.IssueNumber)
	} else {
		fmt.Println("Issue number not specified")
	}
	fmt.Printf("Token Length: %d\n", len(cfg.GitHubToken))
}
