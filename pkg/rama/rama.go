package rama

import (
	"errors"
	"fmt"
	config "kanmitto/internal/config"
	"kanmitto/internal/http"
	"kanmitto/internal/utils"
	"os"
	"time"
)

const RAMA_API = "https://rama.upixels.dev/api"

type RamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
	Error     string    `json:"error"`
}

func Generate(prompt string) (string, error) {
	var configuration config.Config
	utils.ReadJSONFromFile(fmt.Sprintf("%s/.config/kanmit/config.json", os.Getenv("HOME")), &configuration)

	response, err := http.MakeRequest[RamaResponse](http.RequestOptions{
		Method: "POST",
		URL:    RAMA_API + "/generate",
		Body: map[string]interface{}{
			"model":  config.Models[config.ModelType(configuration.Model)],
			"prompt": prompt,
			"format": "json",
			"stream": false,
		},
	})

	if err != nil {
		return "", err
	}

	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.Response, nil
}
