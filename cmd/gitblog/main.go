package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IjichiNijika99/go-gitblog/internal/backup"
	"github.com/IjichiNijika99/go-gitblog/internal/bangumi"
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

	var bgmData []bangumi.Collection
	if cfg.BangumiUser != "" {
		fmt.Printf("Fetching Bangumi data for user: %s...\n", cfg.BangumiUser)
		// 限制拉取前 5 条数据，正好填满 Markdown 表格的一行
		data, err := bangumi.FetchRecentAnime(cfg.BangumiUser, 5)
		if err != nil {
			// 容错处理：不使用 Fatalf，只打印警告，保证博客继续生成
			fmt.Printf("Warning: Failed to fetch Bangumi data: %v\n", err)
		} else {
			bgmData = data
			fmt.Printf("Successfully fetched %d Bangumi records.\n", len(bgmData))
		}
	}

	fmt.Println("Rendering README.md...")
	err = render.BuildREADME(issues, cfg.RepoName, bgmData)
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
