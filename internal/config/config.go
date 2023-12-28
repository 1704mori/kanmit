package config

import "github.com/sashabaranov/go-openai"

type Config struct {
	ConventionalCommit bool   `json:"conventional-commit"`
	Model              string `json:"model"`
	Service            string `json:"service"`
	OllamaAPI          string `json:"ollama_api"`
}

type ModelType string

var openaiModels = map[ModelType]string{
	openai.GPT432K0613:             openai.GPT432K0613,
	openai.GPT432K0314:             openai.GPT432K0314,
	openai.GPT432K:                 openai.GPT432K,
	openai.GPT40613:                openai.GPT40613,
	openai.GPT40314:                openai.GPT40314,
	openai.GPT4TurboPreview:        openai.GPT4TurboPreview,
	openai.GPT4VisionPreview:       openai.GPT4VisionPreview,
	openai.GPT4:                    openai.GPT4,
	openai.GPT3Dot5Turbo1106:       openai.GPT3Dot5Turbo1106,
	openai.GPT3Dot5Turbo0613:       openai.GPT3Dot5Turbo0613,
	openai.GPT3Dot5Turbo0301:       openai.GPT3Dot5Turbo0301,
	openai.GPT3Dot5Turbo16K:        openai.GPT3Dot5Turbo16K,
	openai.GPT3Dot5Turbo16K0613:    openai.GPT3Dot5Turbo16K0613,
	openai.GPT3Dot5Turbo:           openai.GPT3Dot5Turbo,
	openai.GPT3Dot5TurboInstruct:   openai.GPT3Dot5TurboInstruct,
	openai.GPT3TextDavinci003:      openai.GPT3TextDavinci003,
	openai.GPT3TextDavinci002:      openai.GPT3TextDavinci002,
	openai.GPT3TextCurie001:        openai.GPT3TextCurie001,
	openai.GPT3TextBabbage001:      openai.GPT3TextBabbage001,
	openai.GPT3TextAda001:          openai.GPT3TextAda001,
	openai.GPT3TextDavinci001:      openai.GPT3TextDavinci001,
	openai.GPT3DavinciInstructBeta: openai.GPT3DavinciInstructBeta,
	openai.GPT3Davinci:             openai.GPT3Davinci,
	openai.GPT3Davinci002:          openai.GPT3Davinci002,
	openai.GPT3CurieInstructBeta:   openai.GPT3CurieInstructBeta,
	openai.GPT3Curie:               openai.GPT3Curie,
	openai.GPT3Curie002:            openai.GPT3Curie002,
	openai.GPT3Ada:                 openai.GPT3Ada,
	openai.GPT3Ada002:              openai.GPT3Ada002,
	openai.GPT3Babbage:             openai.GPT3Babbage,
	openai.GPT3Babbage002:          openai.GPT3Babbage002,
}

var ollamaModels = map[ModelType]string{
	"info": "we don't know your models yet",
}

var Models = map[string]map[ModelType]string{
	"ollama": ollamaModels,
	"openai": openaiModels,
}
