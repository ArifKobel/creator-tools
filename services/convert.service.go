package services

import (
	"log"
	"os"
)

func convert(videoPath string, outputPath string, args ...string) {
	// check if output file exists
	if _, err := os.Stat(outputPath); err == nil {
		// remove output file if exists
		if err := os.Remove(outputPath); err != nil {
			log.Println("Error: ", err)
		}
	}
	// run ffmpeg command
	err := Ffmpeg(args...)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func ConvertToMp3(videoPath string, outputPath string) {
	convert(videoPath, outputPath, "-i", videoPath, "-vn", "-acodec", "libmp3lame", outputPath)
}

func ConvertToMp4(videoPath string, outputPath string) {
	convert(videoPath, outputPath, "-i", videoPath, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-strict", "experimental", outputPath)
}

func ConvertToGif(videoPath string, outputPath string) {
	convert(videoPath, outputPath, "-i", videoPath, "-vf", "fps=10,scale=320:-1:flags=lanczos", "-c:v", "gif", outputPath)
}

func ConvertToWAV(videoPath string, outputPath string) {
	convert(videoPath, outputPath, "-i", videoPath, "-vn", "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "2", outputPath)
}
