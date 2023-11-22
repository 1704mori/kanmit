package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func main() {
	logger := NewLogger()

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
