package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

func GenerateSubTitles(audioFilePath string) (Response, error) {
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
		return Response{}, err
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error: ", err)
		return Response{}, err
	}
	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("Error: ", err)
		return Response{}, err
	}
	for _, channel := range response.Results.Channels {
		for _, alternative := range channel.Alternatives {
			for _, word := range alternative.Words {
				fmt.Println(word.PunctuatedWord, word.Start, word.End)
			}
		}
	}
	return response, nil
}

func FormatTimeForSubtitles(timeInSeconds float64) string {
	hours := int(timeInSeconds / 3600)
	minutes := int(timeInSeconds / 60)
	seconds := int(timeInSeconds) % 60
	milliseconds := int((timeInSeconds - float64(int(timeInSeconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}
func AddSubtitlesToVideo(videoPath string, subtitles Response, outputPath string) {
	removeExistingFile(outputPath)
	subtitleFileName := createSubtitleFile(subtitles)
	addSubtitlesToVideo(videoPath, subtitleFileName, outputPath)
	removeSubtitleFile(subtitleFileName)
}

func removeExistingFile(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			log.Println("Error: ", err)
		}
	}
}

func createSubtitleFile(subtitles Response) string {
	subtitleFileName, err := RandStringRunes(10)
	if err != nil {
		log.Println("Error: ", err)
	}
	subtitlesFile, err := os.Create(subtitleFileName + ".srt")
	if err != nil {
		log.Println("Error: ", err)
	}
	defer subtitlesFile.Close()

	for _, channel := range subtitles.Results.Channels {
		for _, alternative := range channel.Alternatives {
			for _, word := range alternative.Words {
				_, err := subtitlesFile.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", 1, FormatTimeForSubtitles(word.Start), FormatTimeForSubtitles(word.End), strings.ToUpper(word.Word)))
				if err != nil {
					log.Println("Error: ", err)
				}
			}
		}
	}
	return subtitleFileName
}

func addSubtitlesToVideo(videoPath string, subtitleFileName string, outputPath string) {
	err := Ffmpeg("-i", videoPath, "-vf", "subtitles="+subtitleFileName+".srt:force_style='FontSize=24,Bold=1,Alignment=2,MarginV=80'", "-c:a", "copy", outputPath)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func removeSubtitleFile(subtitleFileName string) {
	err := os.Remove(subtitleFileName + ".srt")
	if err != nil {
		log.Println("Error: ", err)
	}
}
