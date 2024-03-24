package services

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func Ffmpeg(args ...string) error {
	if _, err := os.Stat("./ffmpeg"); os.IsNotExist(err) {
		log.Println("Error: ffmpeg not found")
		return err
	}
	// Specify the codec
	args = append(args, "-c:v", "libx264")
	// Add preset for faster encoding
	args = append(args, "-preset", "ultrafast")
	// Use hardware acceleration if available
	args = append(args, "-hwaccel", "auto")
	cmd := exec.Command("./ffmpeg", args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()
	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			log.Println("Error: ", err)
		}
	}()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
