package main

import (
	"flag"
	"fmt"
	"kanmitto/internal/config"
	"kanmitto/internal/utils"
	"kanmitto/pkg/openai"
	"kanmitto/pkg/rama"
	"os"
	"sort"
	"strings"

	"github.com/eiannone/keyboard"
)

func main() {
	home, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.config/kanmit", home)
	configFile := fmt.Sprintf("%s/config.json", configDir)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.Create(configFile)
		configs := map[string]interface{}{
			"service":             "openai",
			"model":               "gpt-4",
			"conventional-commit": false,
		}

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
			commitMsg, commitMsgError = rama.Generate(prompt)
		default:
			commitMsg, commitMsgError = openai.Generate(prompt)
		}

		stopLoading <- struct{}{}
		logger.ClearLoading()

		if commitMsgError != nil {
			logger.Error(fmt.Sprintf("Error making API request: %s", commitMsgError.Error()))
			return
		}

		if strings.HasPrefix(`"`, commitMsg) {
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
	ollamaApi := flag.String("ollama-api", "", "Set Ollama API e.g: http://localhost:11434")
	model := flag.String("model", "", "OpenAI API model")
	listModels := flag.Bool("models", false, "List available OpenAI API models")
	listSevices := flag.Bool("services", false, "List available APIs")
	showConfigs := flag.Bool("c", false, "Show current configs")

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
		if _, ok := config.Models[config.ModelType(*model)]; !ok {
			logger.Error(fmt.Sprintf("Invalid model: %s", *model))
			return false
		}

		utils.WriteJSONToFile(configFile, "model", *model)
		logger.Info(fmt.Sprintf("Using model: %s", *model))
		return false
	}

	if *listModels {
		var keys []string
		for k := range config.Models {
			keys = append(keys, string(k))
		}

		sort.Strings(keys)

		logger.Info("Available models:")
		for _, k := range keys {
			fmt.Println(k)
		}

		return false
	}

	if *listSevices {
		logger.Info("Available services")
		fmt.Println("ollama")
		fmt.Println("openai")
		return false
	}

	if *showConfigs {
		var configuration config.Config
		utils.ReadJSONFromFile(configFile, &configuration)

		logger.Info("Current configs:")
		fmt.Printf("Conventional commit: %v\n", configuration.ConventionalCommit)
		fmt.Printf("Service: %s\n", configuration.Service)
		fmt.Printf("Ollama API: %s\n", configuration.OllamaAPI)
		fmt.Printf("Model: %s\n", configuration.Model)

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
