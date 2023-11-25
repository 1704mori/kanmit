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
			"api":                 "openai",
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

		switch configuration.Api {
		case "openai":
			commitMsg, commitMsgError = openai.Generate(prompt)
		case "rama":
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
	api := flag.String("api", "", "Set API")
	model := flag.String("model", "", "OpenAI API model")
	listModels := flag.Bool("list-models", false, "List available OpenAI API models")
	listApis := flag.Bool("list-apis", false, "List available APIs")
	showConfigs := flag.Bool("c", false, "Show current configs")

	flag.Parse()

	if *conventionalCommit != "" {
		utils.WriteJSONToFile(configFile, "conventional-commit", *conventionalCommit == "true")
		logger.Info(fmt.Sprintf("Using conventional commit style: %s", *conventionalCommit))
		return false
	}

	if *api != "" {
		utils.WriteJSONToFile(configFile, "api", *api)
		logger.Info(fmt.Sprintf("Using %s API", *api))
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

	if *listApis {
		logger.Info("Available APIs")
		fmt.Println("rama")
		fmt.Println("openai")
	}

	if *showConfigs {
		var configuration config.Config
		utils.ReadJSONFromFile(configFile, &configuration)

		logger.Info("Current configs:")
		fmt.Printf("Conventional commit: %v\n", configuration.ConventionalCommit)
		fmt.Printf("API: %s\n", configuration.Api)
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
