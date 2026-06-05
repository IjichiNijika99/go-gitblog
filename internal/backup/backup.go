package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	gitclient "github.com/IjichiNijika99/go-gitblog/internal/github"
	"github.com/google/go-github/v60/github"
)

// Sync 负责将增量的Issues备份到本地md文件中
func Sync(ctx context.Context, client *gitclient.Client, issues []*github.Issue, backupDir string, targetIssueNum int) error {
	// 1. 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup dir: %w", err)
	}

	// 2. 读取本地已存在的备份 避免重复网络请求
	existingFiles, _ := os.ReadDir(backupDir)
	backedUpMap := make(map[int]bool)
	for _, file := range existingFiles {
		// 只处理md文件
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		// 文件名格式为 "001_Title.md"
		parts := strings.Split(file.Name(), "_")
		if len(parts) > 0 {
			if id, err := strconv.Atoi(parts[0]); err == nil {
				backedUpMap[id] = true
			}
		}
	}

	// 3. 遍历所有文章并按需备份
	for _, issue := range issues {
		num := issue.GetNumber()

		// 如果是Actions指定要求更新的文章(targetIssueNum) 或者是本地还没备份过的新文章
		if targetIssueNum == num || !backedUpMap[num] {
			err := saveSingleIssue(ctx, client, issue, backupDir)
			if err != nil {
				fmt.Printf("Warning: failed to backup issue %d: %v\n", num, err)
			} else {
				fmt.Printf("Successfully backed up Issue #%d\n", num)
			}
		}
	}
	return nil

}

func saveSingleIssue(ctx context.Context, client *gitclient.Client, issue *github.Issue, backupDir string) error {
	title := issue.GetTitle()
	// 替换掉标题中不能作为文件名的非法字符和空格
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, " ", ".")

	fileName := fmt.Sprintf("%d_%s.md", issue.GetNumber(), title)
	filePath := filepath.Join(backupDir, fileName)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入标题超链接和正文
	f.WriteString(fmt.Sprintf("# [%s](%s)\n\n", issue.GetTitle(), issue.GetHTMLURL()))
	f.WriteString(issue.GetBody())

	// 如果有评论
	if issue.GetComments() > 0 {
		comments, err := client.FetchComments(ctx, issue.GetNumber())
		if err != nil {
			return err // 遇到网络错误抛出，避免生成残缺的文件
		}
		for _, c := range comments {
			f.WriteString("\n\n---\n\n")
			f.WriteString(c.GetBody())
		}
	}

	return nil
}
