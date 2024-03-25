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

type SubtitleStyle struct {
	StyleName       string
	FontName        string
	FontSize        string
	PrimaryColor    string
	SecondaryColor  string
	OutlineColor    string
	BackgroundColor string
	IsBold          string
	IsItalic        string
	IsUnderlined    string
	IsStrikethrough string
	ScaleX          string
	ScaleY          string
	LetterSpacing   string
	RotationAngle   string
	BorderStyle     string
	HasOutline      string
	ShadowDepth     string
	TextAlignment   string
	LeftMargin      string
	RightMargin     string
	VerticalMargin  string
	TextEncoding    string
}

func (s *SubtitleStyle) Format() string {
	return fmt.Sprintf("Style: %s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
		s.StyleName, s.FontName, s.FontSize, s.PrimaryColor, s.SecondaryColor, s.OutlineColor,
		s.BackgroundColor, s.IsBold, s.IsItalic, s.IsUnderlined, s.IsStrikethrough, s.ScaleX,
		s.ScaleY, s.LetterSpacing, s.RotationAngle, s.BorderStyle, s.HasOutline, s.ShadowDepth,
		s.TextAlignment, s.LeftMargin, s.RightMargin, s.VerticalMargin, s.TextEncoding)
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

func FormatTimeForSubtitles(timeInSeconds float64) string {
	hours := int(timeInSeconds / 3600)
	minutes := int(timeInSeconds / 60)
	seconds := int(timeInSeconds) % 60
	milliseconds := int((timeInSeconds - float64(int(timeInSeconds))) * 100)
	return fmt.Sprintf("%02d:%02d:%02d.%02d", hours, minutes, seconds, milliseconds)
}
func AddSubtitlesToVideo(videoPath string, subtitles []Word, outputPath string, taskId string) {
	removeExistingFile(outputPath)
	subtitleFileName := createSubtitleFile(subtitles, taskId)
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

func createSubtitleFile(subtitles []Word, taskId string) string {
	subtitleFileName, err := RandStringRunes(10)
	if err != nil {
		log.Println("Error: ", err)
	}
	subtitlesFile, err := os.Create("tasks/" + taskId + "/" + subtitleFileName + ".ass")
	if err != nil {
		log.Println("Error: ", err)
	}
	defer subtitlesFile.Close()
	subtitlesFile.WriteString("[Script Info]\nTitle: Deepgram Subtitles\n")
	subtitlesFile.WriteString("ScriptType: v4.00\nWrapStyle: 0\nScaledBorderAndShadow: yes\nYCbCr Matrix: None\n\n")
	subtitlesFile.WriteString("[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n")
	subtitleStyle := SubtitleStyle{
		StyleName:       "Default",
		FontName:        "Arial",
		FontSize:        "16",
		PrimaryColor:    "&H00FFFFFF",
		SecondaryColor:  "&H000000FF",
		OutlineColor:    "&H00000000",
		BackgroundColor: "&H00000000",
		IsBold:          "1",
		IsItalic:        "0",
		IsUnderlined:    "0",
		IsStrikethrough: "0",
		ScaleX:          "100",
		ScaleY:          "100",
		LetterSpacing:   "0",
		RotationAngle:   "0",
		BorderStyle:     "1",
		HasOutline:      "1",
		ShadowDepth:     "0",
		TextAlignment:   "2",
		LeftMargin:      "10",
		RightMargin:     "10",
		VerticalMargin:  "80",
		TextEncoding:    "1",
	}
	subtitlesFile.WriteString(fmt.Sprintf("%s\n", subtitleStyle.Format()))
	subtitlesFile.WriteString("[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n")
	for _, word := range subtitles {
		subtitlesFile.WriteString(fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n", FormatTimeForSubtitles(word.Start), FormatTimeForSubtitles(word.End), strings.ToUpper(word.Word)))
		if err != nil {
			log.Println("Error: ", err)
		}
	}
	return "tasks/" + taskId + "/" + subtitleFileName
}

func addSubtitlesToVideo(videoPath string, subtitleFileName string, outputPath string) {
	err := Ffmpeg("-i", videoPath, "-vf", "subtitles="+subtitleFileName+".ass:force_style='FontSize=12,Bold=1,Alignment=2,MarginV=80'", "-c:a", "copy", outputPath)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func removeSubtitleFile(subtitleFileName string) {
	err := os.Remove(subtitleFileName + ".ass")
	if err != nil {
		log.Println("Error: ", err)
	}
}
