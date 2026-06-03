package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	owner  string
	repo   string
}

// NewClient 初始化一个带认证信息的Github客户端
func NewClient(ctx context.Context, token, fullrepo string) (*Client, error) {
	parts := strings.Split(fullrepo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repo name: %s", fullrepo)
	}
	// 创建TokenSource认证器
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		})
	tc := oauth2.NewClient(ctx, ts) // 构造带认证器的HTTP客户端
	client := github.NewClient(tc)  // 构造Github API客户端

	return &Client{
		client: client,
		owner:  parts[0],
		repo:   parts[1],
	}, nil
}

// FetchAllIssues 抓取当前仓库下owner的所有Issue
func (c *Client) FetchAllIssues(ctx context.Context) ([]*github.Issue, error) {
	var allIssues []*github.Issue

	// 设置查询参数
	opts := &github.IssueListByRepoOptions{
		State:   "all",
		Creator: c.owner, // 只获取仓库owner的Issue
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		issues, resp, err := c.client.Issues.ListByRepo(ctx, c.owner, c.repo, opts)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch issues: %w", err)
		}

		for _, issue := range issues {
			if !issue.IsPullRequest() {
				allIssues = append(allIssues, issue)
			}
		}

		// NextPage为0说明已经到了最后一页 退出循环
		if resp.NextPage == 0 {
			break
		}

		// 指向下一页继续抓取
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}
