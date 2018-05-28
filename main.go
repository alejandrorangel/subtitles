package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	trasncription := getTranscription("./asrOutput.json")
	newItems := getNewObject(trasncription.Results.Items)
	writeToFile(newItems)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeToFile(items []TranscriptionItem) {
	f, err := os.Create("caption.vtt")
	check(err)
	defer f.Close()

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(fmt.Sprintf("WEBVTT\n\n"))
	for _, item := range items {
		_, err = writer.WriteString(item.toString())
		_, err = writer.WriteString("\n")
	}

	writer.Flush()
}
