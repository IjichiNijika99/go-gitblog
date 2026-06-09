package render

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
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
func BuildREADME(issues []*github.Issue, repoName string, watching []bangumi.Collection, watched []bangumi.Collection) error {
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

	if len(watching) > 0 || len(watched) > 0 {
		buf.WriteString("## Anime Schedule (Powered by Bangumi)\n\n")

		if len(watching) > 0 {
			buf.WriteString("### 在看\n\n")
			renderAnimeGrid(&buf, watching, 3)
		}

		if len(watched) > 0 {
			// 加分割线
			if len(watching) > 0 {
				buf.WriteString("---\n\n")
			}
			buf.WriteString("### 看过\n\n")
			renderAnimeGrid(&buf, watched, 2)
		}
	}

	return os.WriteFile("README.md", buf.Bytes(), 0644)
}

func renderAnimeGrid(buf *bytes.Buffer, data []bangumi.Collection, watch_type int) {
	const colsPerRow = 5

	for i := 0; i < len(data); i += colsPerRow {
		end := i + colsPerRow
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		// 表头
		for _, item := range chunk {
			imgUrl := strings.ReplaceAll(item.Subject.Images.Common, "http://", "https://")
			itemUrl := "https://bgm.tv/subject/" + strconv.Itoa(item.Subject.ID)
			buf.WriteString(fmt.Sprintf("| [<img src=\"%s\" width=\"120\" height=\"170\" title=\"%s\"/>](%s) ", imgUrl, item.Subject.Name, itemUrl))
		}
		buf.WriteString("|\n")

		// 对齐线
		for range chunk {
			buf.WriteString("| :---: ")
		}
		buf.WriteString("|\n")

		// 信息行
		for _, item := range chunk {
			name := item.Subject.NameCN
			if name == "" {
				name = item.Subject.Name
			}

			var info string
			if watch_type == 3 {
				// 在看 展示分数和进度
				epsStr := "?"
				if item.Subject.Eps > 0 {
					epsStr = fmt.Sprintf("%d", item.Subject.Eps)
				}
				info = fmt.Sprintf("**%s**<br/>%.1f<br/>ep. %d/%s", name, item.Subject.Score, item.EpStatus, epsStr)
			} else if watch_type == 2 {
				// 看过 展示评分和吐槽
				comment := item.Comment
				if comment != "" {
					comment = strings.ReplaceAll(comment, "|", "｜")
					comment = strings.ReplaceAll(comment, "\n", " ")
					//runes := []rune(comment)
					//if len(runes) > 18 {
					//	comment = string(runes[:18]) + "<details><summary></summary>" + string(runes[18:]) + "</details>"
					//}
					////buf.WriteString("<details><summary>显示更多</summary>\n\n")
					//info += fmt.Sprintf("<br/>*%s*", comment)
				}

				if item.Rate == 0 {
					info = fmt.Sprintf("**%s**<br/>%s<br/><details><summary></summary>%s</details>", name, " ", comment)
				} else {
					info = fmt.Sprintf("**%s**<br/>%d<br/><details><summary></summary>%s</details>", name, item.Rate, comment)
				}

			}
			buf.WriteString(fmt.Sprintf("| %s ", info))
		}
		// 表格结束后留白
		buf.WriteString("|\n\n")
	}
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
