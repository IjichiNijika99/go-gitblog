package config

import (
	"flag"
	"fmt"
	"os"
)

// Config 保存程序运行所需的所有配置参数
type Config struct {
	GitHubToken string
	RepoName    string
	IssueNumber int    // 可选 默认为0代表处理所有issue
	BackupDir   string // 备份目录 默认为"BACKUP"
}

// Parse 解析命令行参数并返回Config实例
func Parse() (*Config, error) {
	cfg := &Config{
		BackupDir: "BACKUP", // 默认的备份目录
	}
	// 定义命令行参数
	flag.StringVar(&cfg.GitHubToken, "token", "", "GitHub Personal Access Token (Required)")
	flag.StringVar(&cfg.RepoName, "repo", "", "GitHub Repository Name (Required)")
	flag.IntVar(&cfg.IssueNumber, "issue", 0, "Specific issue number to process (Optional)")

	// 重写Usage函数 自定义输出
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0]) // os.Args[0]代表执行的文件名
		flag.PrintDefaults()
	}

	// 执行解析
	flag.Parse()

	if cfg.GitHubToken == "" || cfg.RepoName == "" {
		flag.Usage()
		return nil, fmt.Errorf("token and repo are required")
	}
	return cfg, nil
}
