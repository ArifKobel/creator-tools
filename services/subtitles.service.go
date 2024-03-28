package services

import (
	"context"
	"encoding/json"
	"fmt"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/prerecorded/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
)

type Response struct {
	Metadata Metadata `json:"metadata"`
	Results  Results  `json:"results"`
}

type Metadata struct {
	TransactionKey string               `json:"transaction_key"`
	RequestID      string               `json:"request_id"`
	Sha256         string               `json:"sha256"`
	Created        string               `json:"created"`
	Duration       float64              `json:"duration"`
	Channels       int                  `json:"channels"`
	Models         []string             `json:"models"`
	ModelInfo      map[string]ModelInfo `json:"model_info"`
}

type ModelInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
}

type Results struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript string     `json:"transcript"`
	Confidence float64    `json:"confidence"`
	Words      []Word     `json:"words"`
	Paragraphs Paragraphs `json:"paragraphs"`
}

type Word struct {
	Word           string  `json:"word"`
	Start          float64 `json:"start"`
	End            float64 `json:"end"`
	Confidence     float64 `json:"confidence"`
	PunctuatedWord string  `json:"punctuated_word"`
}

type Paragraphs struct {
	Transcript string      `json:"transcript"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Paragraph struct {
	Sentences []Sentence `json:"sentences"`
	NumWords  int        `json:"num_words"`
	Start     float64    `json:"start"`
	End       float64    `json:"end"`
}

type Sentence struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

func GenerateSubTitles(audioFilePath string) ([]Word, error) {
	client.InitWithDefault()

	ctx := context.Background()

	options := interfaces.PreRecordedTranscriptionOptions{
		Model:       "nova-2-general",
		SmartFormat: true,
		Language:    "de",
	}

	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	res, err := dg.FromFile(ctx, audioFilePath, options)
	if err != nil {
		fmt.Println("Error: ", err)
		return []Word{}, err
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error: ", err)
		return []Word{}, err
	}
	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("Error: ", err)
		return []Word{}, err
	}
	var words []Word
	for _, channel := range response.Results.Channels {
		for _, alternative := range channel.Alternatives {
			words = append(words, alternative.Words...)
		}
	}
	return words, nil
}
