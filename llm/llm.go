package llm

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/openai"
)

var (
	LLM_Model   = os.Getenv("JobSynk_LLM_Model")
	LLM_BaseUrl = os.Getenv("JobSynk_LLM_BaseUrl")
	LLM_Token   = os.Getenv("JobSynk_LLM_Token")
)

// 大模型输出
type LLApplicationEvaluationMOutPut struct {
	Rating   int    `json:"rating"`   // 评分
	Greeting string `json:"greeting"` // 问候语
}

// 用户输入
type JobApplication struct {
	Resume   string `json:"resume"`   // 简历内容
	Company  string `json:"company"`  // 公司名称
	JobTitle string `json:"jobTitle"` // 工作名称
	Jd       string `json:"jd"`       // 工作内容
}

func NewLLM(ctx context.Context, jobApplicationCH chan JobApplication, llApplicationEvaluationMOutPutCH chan LLApplicationEvaluationMOutPut) error {
	llm, err := openai.New(openai.WithModel(LLM_Model), openai.WithBaseURL(LLM_BaseUrl), openai.WithToken(LLM_Token))
	if err != nil {
		return err
	}
	largeContext := Prompt
	for jobApplication := range jobApplicationCH {
		b, _ := json.Marshal(jobApplication)
		// 构建消息
		messages := []llms.MessageContent{
			{
				Role: llms.ChatMessageTypeSystem,
				Parts: []llms.ContentPart{
					// Mark the large context for caching
					llms.WithCacheControl(llms.TextPart(largeContext), anthropic.EphemeralCache()),
				},
			},
			{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextPart(string(b)),
				},
			},
		}
		resp, err := llm.GenerateContent(ctx, messages, llms.WithMaxTokens(200), anthropic.WithPromptCaching())
		if err != nil {
			slog.Error("Failed to generate LLM content", "error", err, "jobApplication", jobApplication)
			continue
		}
		content := resp.Choices[0].Content
		var output LLApplicationEvaluationMOutPut
		// 解析成struct并返回
		err = json.Unmarshal([]byte(content), &output)
		if err != nil {
			slog.Error("Failed to parse LLM response JSON", "error", err, "rawContent", content, "jobApplication", jobApplication)
			continue
		}
		llApplicationEvaluationMOutPutCH <- output
	}
	return nil
}
