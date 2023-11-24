package main

import (
	"flag"
	"fmt"
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
	}

	logger := NewLogger()

	conventionalCommit := flag.String("conventional-commit", "", "Use conventional commit style")
	model := flag.String("model", "", "OpenAI API model")
	listModels := flag.Bool("list-models", false, "List available OpenAI API models")

	flag.Parse()

	if *listModels {
		var keys []string
		for k := range MODELS {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		logger.Info("Available models:")
		for _, k := range keys {
			fmt.Println(k)
		}

		return
	}

	if *conventionalCommit != "" {
		WriteJSONToFile(configFile, "conventional-commit", *conventionalCommit == "true")
		logger.Info(fmt.Sprintf("Using conventional commit style: %s", *conventionalCommit))
		return
	}

	if *model != "" {
		if _, ok := MODELS[*model]; !ok {
			logger.Error(fmt.Sprintf("Invalid model: %s", *model))
			return
		}

		WriteJSONToFile(configFile, "model", *model)
		logger.Info(fmt.Sprintf("Using model: %s", *model))
		return
	}

	if !IsGitRepo() {
		logger.Error("Not a git repository")
		return
	}

	if !HasStagedChanges() {
		logger.Error("No staged changes")
		return
	}

	for {
		diff := GetStagedDiff()
		prompt := ConstructPrompt(diff)

		stopLoading := make(chan struct{})
		go logger.AnimateLoading("Generating commit message...", stopLoading)
		commitMsg, commitMsgError := CallOpenaiAPI(prompt)
		stopLoading <- struct{}{}
		logger.ClearLoading()

		logger.Success(fmt.Sprintf("Suggested commit message: %s", commitMsg))
		logger.Info("Do you want to generate another commit message or commit the suggested one? g(enerate)/(c)ommit)/(q)uit:")

		if commitMsgError != nil {
			logger.Error(fmt.Sprintf("Error making API request: %s", commitMsgError.Error()))
			return
		}

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
			Commit(commitMsg)
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
