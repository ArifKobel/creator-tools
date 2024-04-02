package services

import (
	"crypto/rand"
	"os"
	"time"

	ffprobe "github.com/vansante/go-ffprobe"
)

const letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func RandStringRunes(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = letterRunes[int(b[i])%len(letterRunes)]
	}
	return string(b), nil
}

func Mp4Duration(path string) (float64, error) {
	data, err := ffprobe.GetProbeData(path, 1200000*time.Millisecond)
	if err != nil {
		return 0, err
	}
	duration := data.Format.Duration().Seconds()
	return duration, nil
}

func GenerateThumbnail(path string, output string) error {
	err := Ffmpeg("-i", path, "-ss", "00:00:01.000", "-vframes", "1", output)
	if err != nil {
		return err
	}
	return nil
}
