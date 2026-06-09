package render

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/IjichiNijika99/go-gitblog/internal/bangumi"
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
func BuildREADME(issues []*github.Issue, repoName string, bgmData []bangumi.Collection) error {
	var buf bytes.Buffer

	// 1. Header
	buf.WriteString(fmt.Sprintf("## [Gitblog](https://github.com/%s)\n", repoName))
	buf.WriteString("抓住那绵软的飞机云  让山神为我们指引  \n")
	buf.WriteString("不要放开相牵的双手  别让共同的梦境消散  \n")
	buf.WriteString("一人一半的羽翼  互相依偎  飞翔吧  \n")
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

	if len(bgmData) > 0 {
		buf.WriteString("## Anime Schedule (Powered by Bangumi)\n\n")

		// 5.1 构建表头 (海报图片行)，固定高度防止页面跳动
		for _, item := range bgmData {
			// 将 http 强行替换为 https 避免 GitHub 混合内容拦截
			imgUrl := strings.ReplaceAll(item.Subject.Images.Common, "http://", "https://")
			buf.WriteString(fmt.Sprintf("| <img src=\"%s\" width=\"120\" height=\"170\" /> ", imgUrl))
		}
		buf.WriteString("|\n")

		// 5.2 构建 Markdown 表格分割线 (居中对齐)
		for range bgmData {
			buf.WriteString("| :---: ")
		}
		buf.WriteString("|\n")

		// 5.3 构建信息行 (名字、评分、状态、短评)
		for _, item := range bgmData {
			name := item.Subject.NameCN
			if name == "" {
				name = item.Subject.Name // 如果没有中文名，降级使用原名
			}

			// 解析状态字典
			status := "看过"
			if item.Type == 3 {
				status = "在看"
			} else if item.Type == 1 {
				status = "想看"
			} else if item.Type == 4 {
				status = "搁置"
			}

			// 清洗短评：去除管道符和换行符，防止破坏 Markdown 表格
			comment := item.Comment
			comment = strings.ReplaceAll(comment, "|", "｜")
			comment = strings.ReplaceAll(comment, "\n", " ")
			// 限制短评长度，防止表格被无限撑开
			// 注意：这里粗略截取字符串，实际生产可使用 []rune 处理中文字符切片
			runes := []rune(comment)
			if len(runes) > 18 {
				comment = string(runes[:18]) + "..."
			}

			info := fmt.Sprintf("**%s**<br/>⭐ %d/10<br/>状态: %s", name, item.Rate, status)
			if comment != "" {
				info += fmt.Sprintf("<br/>💬 *%s*", comment)
			}
			buf.WriteString(fmt.Sprintf("| %s ", info))
		}
		buf.WriteString("|\n\n")
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
