package main

import (
	"flag"
	"fmt"
	"kanmitto/internal/config"
	"kanmitto/internal/utils"
	ollama "kanmitto/pkg/ollama"
	"kanmitto/pkg/openai"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

var Version string = "dev"

var configs = map[string]interface{}{
	"service":             "openai",
	"model":               "gpt-4",
	"conventional-commit": false,
}

func main() {
	home, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.config/kanmit", home)
	configFile := fmt.Sprintf("%s/config.json", configDir)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.Create(configFile)

		for k, v := range configs {
			utils.WriteJSONToFile(configFile, k, v)
		}
	}

	logger := utils.NewLogger()

	ok := validateFlags(configFile, logger)

	if !ok {
		return
	}

	if !utils.IsGitRepo() {
		logger.Error("Not a git repository")
		return
	}

	if !utils.HasStagedChanges() {
		logger.Error("No staged changes")
		return
	}

	var configuration config.Config
	utils.ReadJSONFromFile(configFile, &configuration)

	if configuration.Service == "ollama" && configuration.OllamaAPI == "" {
		logger.Error("Ollama's API is empty")
		return
	}

	for {
		diff := utils.GetStagedDiff()
		prompt := utils.ConstructPrompt(diff)

		stopLoading := make(chan struct{})
		go logger.AnimateLoading("Generating commit message...", stopLoading)
		var (
			commitMsg      string
			commitMsgError error
		)

		switch configuration.Service {
		case "openai":
			commitMsg, commitMsgError = openai.Generate(prompt)
		case "ollama":
			commitMsg, commitMsgError = ollama.Generate(prompt)
		default:
			commitMsg, commitMsgError = openai.Generate(prompt)
		}

		stopLoading <- struct{}{}
		logger.ClearLoading()

		if commitMsgError != nil {
			logger.Error(fmt.Sprintf("Error making API request: %s", commitMsgError.Error()))
			return
		}

		if strings.HasPrefix(commitMsg, `"`) {
			commitMsg = strings.Trim(commitMsg, `"`)
		}

		logger.Success(fmt.Sprintf("Suggested commit message: %s", commitMsg))
		logger.Info("Do you want to generate another commit message or commit the suggested one? g(enerate)/(c)ommit)/(q)uit:")

		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		switch char {
		case 'g':
			continue
		case 'c':
			logger.Info("Committing...")
			utils.Commit(commitMsg)
			return
		case 'q':
			logger.Info("Quitting without committing...")
			return
		default:
			logger.Info("Invalid option. Quitting without committing...")
			return
		}
	}
}

func validateFlags(configFile string, logger *utils.Logger) bool {
	conventionalCommit := flag.String("conventional-commit", "", "Use conventional commit style")
	service := flag.String("service", "", "Set API")
	model := flag.String("model", "", "OpenAI API model")
	ollamaApi := flag.String("ollama-api", "", "Set Ollama API e.g: http://localhost:11434")
	listModels := flag.Bool("list-models", false, "List available OpenAI API models")
	listSevices := flag.Bool("list-services", false, "List available AI Services")
	showConfig := flag.Bool("c", false, "Show current configs")
	resetConfig := flag.Bool("reset-config", false, "Resets the configuration to its default")
	version := flag.Bool("v", false, "Show kanmit version")

	flag.Parse()

	if *conventionalCommit != "" {
		utils.WriteJSONToFile(configFile, "conventional-commit", *conventionalCommit == "true")
		logger.Info(fmt.Sprintf("Using conventional commit style: %s", *conventionalCommit))
		return false
	}

	if *service != "" {
		if !strings.Contains("openai ollama", *service) {
			logger.Error("Invalid service")
			return false
		}

		utils.WriteJSONToFile(configFile, "service", *service)
		logger.Info(fmt.Sprintf("Using %s service", *service))
		return false
	}

	if *ollamaApi != "" {
		utils.WriteJSONToFile(configFile, "ollama_api", *ollamaApi)
		logger.Info(fmt.Sprintf("Ollama API set to %s", *ollamaApi))
		return false
	}

	if *model != "" {
		var configuration config.Config
		if configuration.Service == "openai" {
			if _, ok := config.Models["openai"][config.ModelType(*model)]; !ok {
				logger.Error(fmt.Sprintf("Invalid OpenAI model: %s", *model))
				return false
			}
		}

		utils.WriteJSONToFile(configFile, "model", *model)
		logger.Info(fmt.Sprintf("Using model: %s", *model))
		return false
	}

	if *listModels {
		logger.Info("Available models:")
		for provider, models := range config.Models {
			fmt.Println("Provider:", provider)

			for _, modelName := range models {
				fmt.Println(modelName)
			}
			fmt.Println("-------------------")
		}

		return false
	}

	if *listSevices {
		logger.Info("Available services")
		fmt.Println("ollama")
		fmt.Println("openai")
		return false
	}

	if *showConfig {
		var configuration config.Config
		utils.ReadJSONFromFile(configFile, &configuration)

		logger.Info("Current configs:")
		fmt.Printf("Conventional commit: %v\n", configuration.ConventionalCommit)
		fmt.Printf("Service: %s\n", configuration.Service)
		fmt.Printf("Ollama API: %s\n", configuration.OllamaAPI)
		fmt.Printf("Model: %s\n", configuration.Model)

		return false
	}

	if *resetConfig {
		os.Remove(configFile)

		for k, v := range configs {
			utils.WriteJSONToFile(configFile, k, v)
		}

		return false
	}

	if *version {
		logger.Info(fmt.Sprintf("Kanmit version: %s", Version))
		return false
	}

	return true
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
