package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

func GetStagedDiff() []byte {
	out, _ := exec.Command("git", "diff", "--staged").Output()
	return out
}

func ConstructPrompt(stagedFiles []byte) string {
	files := strings.Split(string(stagedFiles), "\n")
	prompt := `I have made the following changes and I will provide them to you using the command "git diff --staged".
	Summarize these changes in a commit message with a maximum length of 50 characters.
	NB: The commit message should be in the imperative mood, e.g. "Add feature" rather than "Added feature".
			Also, do not end the commit message with a period neither you should add quotes around it.
	`

	for _, file := range files {
		if len(file) == 0 {
			continue
		}

		prompt += file + "\n"
	}

	return prompt + "\n\nCommit message:"
}

func Commit(message string) {
	out, err := exec.Command("git", "commit", "-m", message).Output()
	if err != nil {
		fmt.Printf("Error committing: %s", err.Error())
		return
	}

	fmt.Println(string(out))
}
