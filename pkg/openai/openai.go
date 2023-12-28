package openai

import (
	"context"
	"fmt"
	config "kanmitto/internal/config"
	"kanmitto/internal/utils"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Generate(prompt string) (string, error) {
	var configuration config.Config
	utils.ReadJSONFromFile(fmt.Sprintf("%s/.config/kanmit/config.json", os.Getenv("HOME")), &configuration)

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: config.Models["openai"][config.ModelType(configuration.Model)],
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
