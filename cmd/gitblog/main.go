package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IjichiNijika99/go-gitblog/internal/backup"
	"github.com/IjichiNijika99/go-gitblog/internal/config"
	gitclient "github.com/IjichiNijika99/go-gitblog/internal/github"
	"github.com/IjichiNijika99/go-gitblog/internal/render"
)

func main() {
	// 解析配置
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("Init failed: %v", err)
	}

	fmt.Println("=== Init success ===")
	fmt.Println("=== Starting GitBlog Build ===")

	// 初始化Github客户端
	ctx := context.Background()
	client, err := gitclient.NewClient(ctx, cfg.GitHubToken, cfg.RepoName)
	if err != nil {
		log.Fatalf("Failed to create Github client: %v", err)
	}

	fmt.Println("Fetching issues from Github...")
	issues, err := client.FetchAllIssues(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch issues: %v", err)
	}

	fmt.Printf("Successfully fetched %d issues authored by you.\n", len(issues))

	fmt.Println("Rendering README.md...")
	err = render.BuildREADME(issues, cfg.RepoName)
	if err != nil {
		log.Fatalf("Failed to render README.md: %v", err)
	}
	fmt.Println("README.md generated successfully!")

	fmt.Println("Starting issue backup process...")
	err = backup.Sync(ctx, client, issues, cfg.BackupDir, cfg.IssueNumber)
	if err != nil {
		log.Fatalf("Failed to backup issues: %v", err)
	}
	fmt.Println("Backup process completed!")
	fmt.Println("==============================")

	//// 简单打印前 3 篇文章的标题看看效果
	//limit := 3
	//if len(issues) < 3 {
	//	limit = len(issues)
	//}
	//fmt.Println("\nRecent articles:")
	//for i := 0; i < limit; i++ {
	//	fmt.Printf("- [%d] %s (Labels: %d)\n",
	//		issues[i].GetNumber(),
	//		issues[i].GetTitle(),
	//		len(issues[i].Labels))
	//}
}
