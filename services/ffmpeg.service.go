package services

import (
	"os/exec"
)

func Ffmpeg(args ...string) error {
	// Specify the codec
	args = append(args, "-c:v", "libx264")
	// Add preset for faster encoding
	args = append(args, "-preset", "ultrafast")
	// Use hardware acceleration if available
	args = append(args, "-hwaccel", "auto")
	cmd := exec.Command("ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
