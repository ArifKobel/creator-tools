package services

import (
	"log"
	"os"
	"strconv"
)

func Concat(videoPaths []string, outputPath string) {
	// check if output file exists
	if _, err := os.Stat(outputPath); err == nil {
		// remove output file if exists
		if err := os.Remove(outputPath); err != nil {
			log.Println("Error: ", err)
		}
	}
	// create a list of input files
	var inputFiles []string
	for _, videoPath := range videoPaths {
		inputFiles = append(inputFiles, "-i", videoPath)
	}
	// run ffmpeg command
	err := Ffmpeg(append(inputFiles, "-filter_complex", "concat=n="+strconv.Itoa(len(videoPaths))+":v=1:a=1", "-f", "mp4", outputPath)...)
	if err != nil {
		log.Println("Error: ", err)
	}
}
