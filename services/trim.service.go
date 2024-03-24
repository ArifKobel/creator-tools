package services

import (
	"fmt"
	"os"
	"strconv"
)

func Trim(videoPath string, startInSeconds int, endInSeconds int, outputPath string) error {
	err := removeFileIfExists(outputPath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	startTime := SecondsToHMS(startInSeconds)
	duration := SecondsToHMS(endInSeconds - startInSeconds)

	err = Ffmpeg("-ss", startTime, "-i", videoPath, "-t", duration, "-c", "copy", outputPath)
	if err != nil {
		return fmt.Errorf("failed to trim video: %w", err)
	}

	return nil
}

func removeFileIfExists(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return err
		}
	}
	return nil
}

func pad(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

func SecondsToHMS(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%s:%s:%s", pad(hours), pad(minutes), pad(seconds))
}
