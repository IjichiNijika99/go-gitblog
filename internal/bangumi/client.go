package bangumi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Collection Bangumi API 返回的单条收藏数据结构
type Collection struct {
	Rate     int    `json:"rate"`      // 打分 (1-10)
	Type     int    `json:"type"`      // 收藏状态: 1:想看 2:看过 3:在看 4:搁置 5:抛弃
	Comment  string `json:"comment"`   // 吐槽
	UpdateAt string `json:"update_at"` // 更新时间
	EpStatus int    `json:"ep_status"` // 当前看到的集数
	Subject  struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		NameCN string  `json:"name_cn"`
		Score  float64 `json:"score"`
		Eps    int     `json:"eps"`
		Images struct {
			Large  string `json:"large"`
			Common string `json:"common"`
		} `json:"images"`
	} `json:"subject"`
}

// APIResponse 最外层的 JSON 结构
type APIResponse struct {
	Data []Collection `json:"data"`
}

// FetchRecentAnime 获取用户最近标记的动画状态
// collectionType: 2为"看过"  3为"在看"
func FetchRecentAnime(username string, collection_type int, limit int) ([]Collection, error) {
	// subject_type=2 Anime
	url := fmt.Sprintf("https://api.bgm.tv/v0/users/%s/collections?subject_type=2&type=%d&limit=%d", username, collection_type, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Bangumi API 强制要求提供自定义的 User-Agent 否则会返回 403
	userAgent := fmt.Sprintf("%s/go-gitblog (https://github.com/%s/go-gitblog)", username, username)
	req.Header.Set("User-Agent", userAgent)

	// 10s Timeout
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bangumi api returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	return apiResp.Data, nil
}
