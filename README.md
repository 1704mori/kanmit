# Kanmit

Kanmit is a tool that facilitates easy and straightforward commit messages, making the process of committing code changes a little less painful.

## Getting Started

### Build and Install

```bash
git clone https://github.com/1704mori/kanmit.git
cd kanmit

go build -o kanmit
sudo mv kanmit /usr/local/bin
```

### Usage

Set up OpenAI API key:

```bash
export OPENAI_API_KEY=your-api-key # or add this line to your .bashrc or .zshrc
```

Generate a commit message using Kanmit:

```bash
cd /path/to/your/repo
kanmit
```

Follow the on-screen prompts to re-generate a commit message, commit the changes, or exit the tool.

## Contribution

If you would like to contribute to this project or have any suggestions, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License (LICENSE).


## Why "Kanmit" or "Kanmitto"?

The name "Kanmit" is a combination of "Kan" (簡) from the Japanese word "簡単" (Kantan), meaning simple or easy, and "mit" from "commit." It's a silly name that I came up with while I was coding in the middle of the night, but it stuck.
