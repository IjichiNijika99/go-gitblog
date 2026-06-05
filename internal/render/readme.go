package render

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v60/github"
)

// IgnoreLabels 需要忽略的标签
var IgnoreLabels = map[string]bool{
	"Top":     true,
	"TODO":    true,
	"Friends": true,
	"About":   true,
	"Things":  true,
}

const AnchorNumber int = 5 // 超过5篇折叠

// BuildREADME 将issues转为README.md并写入
func BuildREADME(issues []*github.Issue, repoName string) error {
	var buf bytes.Buffer

	// 1. Header
	buf.WriteString(fmt.Sprintf("## [Gitblog](https://github.com/%s)\n", repoName))
	buf.WriteString("抓住那绵软的飞机云  让山神为我们指引\n")
	buf.WriteString("不要放开相牵的双手  别让共同的梦境消散\n")
	buf.WriteString("一人一半的羽翼  互相依偎  飞翔吧\n")
	buf.WriteString(fmt.Sprintf("[RSS Feed (TODO)](https://raw.githubusercontent.com/%s/master/feed.xml)\n\n", repoName))

	// 2. Top Articles (Top Issues)
	buf.WriteString("## 置顶文章\n")
	for _, issue := range issues {
		if hasLabel(issue, "Top") {
			buf.WriteString(formatIssue(issue))
		}
	}
	buf.WriteString("\n")

	// 3. Recent Articles (Recent Issues)
	buf.WriteString("## 最近更新\n")
	recentCount := 0
	for _, issue := range issues {
		buf.WriteString(formatIssue(issue))
		recentCount++
		if recentCount >= 5 {
			break
		}
	}
	buf.WriteString("\n")

	// 4. Articles By Labels
	labelMap := make(map[string][]*github.Issue)
	var labelNames []string

	// 按标签归类
	for _, issue := range issues {
		for _, label := range issue.Labels {
			name := label.GetName()
			// 跳过忽略的标签
			if IgnoreLabels[name] {
				continue
			}
			// 第一次遇到这个标签 记录它的名字用于排序
			if _, exists := labelMap[name]; !exists {
				labelNames = append(labelNames, name)
			}
			labelMap[name] = append(labelMap[name], issue)

		}
	}

	sort.Strings(labelNames)

	for _, labelName := range labelNames {
		buf.WriteString("## " + labelName + "\n\n")
		labelIssues := labelMap[labelName]

		for index, issue := range labelIssues {
			// 超过阈值使用 <details> 折叠
			if index == AnchorNumber {
				buf.WriteString("<details><summary>显示更多</summary>\n\n")
			}
			buf.WriteString(formatIssue(issue))
		}
		if len(labelIssues) > AnchorNumber {
			buf.WriteString("</details>\n\n")
		} else {
			buf.WriteString("\n")
		}
	}

	return os.WriteFile("README.md", buf.Bytes(), 0644)
}

// hasLabel 判断Issue是否包含某个特定的标签
func hasLabel(issue *github.Issue, targetLabel string) bool {
	for _, label := range issue.Labels {
		// 忽略大小写进行比较
		if strings.EqualFold(label.GetName(), targetLabel) {
			return true
		}
	}
	return false
}

// formatIssue 格式化单篇Issue为Markdown列表项
func formatIssue(issue *github.Issue) string {
	// 获取时间 格式yyyy-mm-dd
	timeStr := issue.GetCreatedAt().Format("2006-01-02")
	return fmt.Sprintf("- [%s](%s)--%s\n", issue.GetTitle(), issue.GetHTMLURL(), timeStr)
}
