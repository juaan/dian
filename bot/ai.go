package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	config "github.com/dian/etc"
)

type (
	TailvySearchResponse struct {
		Answer            string                 `json:"answer"`
		Query             string                 `json:"query"`
		ResponseTime      float64                `json:"response_time"`
		FollowUpQuestions []string               `json:"follow_up_questions"`
		Images            []string               `json:"images"`
		Results           []TailvyResultResponse `json:"results"`
	}

	TailvyResultResponse struct {
		Title      string `json:"title"`
		Url        string `json:"url"`
		Content    string `json:"content"`
		RawContent string `json:"raw_content"`
		Score      string `json:"score"`
	}

	TailvySearchRequest struct {
		ApiKey            string `json:"api_key"`
		Query             string `json:"query"`
		SearchDepth       string `json:"search_depth"`
		IncludeAnswer     bool   `json:"include_answer"`
		IncludeImages     bool   `json:"include_images"`
		IncludeRawContent bool   `json:"include_raw_content"`
		MaxResults        int    `json:"max_results"`
	}
)

func askTailvy(ctx context.Context, question string) (TailvySearchResponse, error) {
	host := "https://api.tavily.com"
	url := fmt.Sprintf("%s/search", host)
	method := http.MethodPost

	param := TailvySearchRequest{
		ApiKey:        config.Get().TailvyToken,
		Query:         question,
		IncludeAnswer: true,
		SearchDepth:   "basic",
	}
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(param)
	if err != nil {
		return TailvySearchResponse{}, fmt.Errorf("[askTailvy] failed to encode payload, err: %+v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, &payload)
	if err != nil {
		return TailvySearchResponse{}, fmt.Errorf("[askTailvy] failed to build request, err: %+v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return TailvySearchResponse{}, fmt.Errorf("[askTailvy] failed to execute request, err: %+v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TailvySearchResponse{}, fmt.Errorf("[askTailvy] failed to read body, err: %+v", err)
	}

	if res.StatusCode != 200 {
		err := fmt.Errorf("error request to tailvy api")
		return TailvySearchResponse{}, err
	}

	tailvyResp := TailvySearchResponse{}
	err = json.Unmarshal(body, &tailvyResp)
	if err != nil {
		return TailvySearchResponse{}, fmt.Errorf("[askTailvy] failed to unmarshal response body, err: %+v", err)
	}

	return tailvyResp, nil
}
