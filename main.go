package main

import (
	"fmt"
	"os"

	"github.com/ArifKobel/creator-tools/services"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	inputFile := "input.mp4"
	taskId, err := services.RandStringRunes(20)
	if err != nil {
		fmt.Println(err)
	}
	// create a folder for the task
	err = os.Mkdir("tasks/"+taskId, 0755)
	if err != nil {
		fmt.Println(err)
	}
	services.ConvertToWAV(inputFile, "tasks/"+taskId+"/input.wav")
	subtitles, err := services.GenerateSubTitles("tasks/" + taskId + "/input.wav")
	if err != nil {
		fmt.Println(err)
	}
	services.RemoveFile("tasks/" + taskId + "/input.wav")
	services.AddSubtitlesToVideo(inputFile, subtitles, "tasks/"+taskId+"/output.mp4", taskId)

}
