package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var PROMPTS = map[string]string{
	"simple": `I have made the following changes and I will provide them to you using the command "git diff --staged".
	Summarize these changes in a commit message with a maximum length of 50 characters.
	NB: The commit message should be in the imperative mood, e.g. "Add feature" rather than "Added feature".
			Also, do not end the commit message with a period neither you should add quotes around it.
	`,
	"conventional": `I have made the following changes and I will provide them to you using the command "git diff --staged".
	Summarize these changes in a commit message using the conventional commit style with a maximum length of 50 characters.

	Conventional commit style:
	<type>[optional scope]: <description>
	E.g. "feat: add new feature" or "fix: fix bug" or with scope "feat(parser): add new parser feature"

	NB: The commit message should be in the imperative mood, e.g. "feat: add feature" rather than "feat: added feature".
			Also, do not end the commit message with a period neither you should add quotes around it.
	`,
}

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

type Config struct {
	ConventionalCommit bool   `json:"conventional-commit"`
	Model              string `json:"model"`
}

func ConstructPrompt(stagedFiles []byte) string {
	var config Config
	ReadJSONFromFile(fmt.Sprintf("%s/.config/kanmit/config.json", os.Getenv("HOME")), &config)

	files := strings.Split(string(stagedFiles), "\n")
	prompt := PROMPTS["simple"]

	if config.ConventionalCommit {
		prompt = PROMPTS["conventional"]
	}

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
