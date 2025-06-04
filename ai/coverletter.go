package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/telman03/autolead-ai/parser"
	"github.com/telman03/autolead-ai/scraper"
)

func GenerateCoverLetter(resume *parser.ResumeData, job scraper.Job) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set in environment")
	}

	client := openai.NewClient(apiKey)

	resumeSummary := fmt.Sprintf(
		"My name is %s. I have experience in: %s. My skills include: %s.",
		resume.Name,
		formatList(resume.Experience),
		formatList(resume.Skills),
	)

	prompt := fmt.Sprintf(`
Write a short and professional cover letter (max 200 words) for the following job:

📌 Job Title: %s
🏢 Company: %s
🔖 Tags: %v

My resume:
%s

Use a confident and humble tone. Do not copy the job description.
`, job.Title, job.Company, job.Tags, resumeSummary)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func formatList(list []string) string {
	out := ""
	for _, item := range list {
		if item != "" {
			out += "- " + item + "\n"
		}
	}
	return out
}
