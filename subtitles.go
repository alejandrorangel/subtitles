package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const TRESHOLD float64 = 0.2

// TranscribeObject : structure obtain from AWS Transcribe Service
type TranscribeObject struct {
	JobName   string `json:"jobName"`
	AccountID string `json:"accountId"`
	Status    string `json:"status"`
	Results   Result `json:"results"`
}

// Result : structure obteain from AWS Transcribe Service
type Result struct {
	Items []TranscriptionItem `json:"items"`
}

// TranscriptionItem : structure obteain from AWS Transcribe Service
type TranscriptionItem struct {
	StartTime    float64        `json:"start_time,string"`
	EndTime      float64        `json:"end_time,string"`
	Alertantives []Aleternative `json:"alternatives"`
	Type         string         `json:"type"`
}

// Aleternative : structure obteain from AWS Transcribe Service
type Aleternative struct {
	Confidance string `json:"confidence"`
	Content    string `json:"content"`
}

func (p TranscribeObject) toString() string {
	return fmt.Sprintf("%s - %s - %s", p.JobName, p.AccountID, p.Status)
}

func (p Aleternative) toString() string {
	return fmt.Sprintf("%s", p.Content)
}

func (p TranscriptionItem) toString() string {
	return fmt.Sprintf("00:%.3f --> 00:%.3f \n%s\n", p.StartTime, p.EndTime, p.Alertantives[0].toString())
}

func getTranscription() TranscribeObject {
	raw, err := ioutil.ReadFile("./asrOutput.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	var c TranscribeObject
	json.Unmarshal(raw, &c)
	return c
}

func getBuffer(current TranscriptionItem, next TranscriptionItem, textBuffer bytes.Buffer) bytes.Buffer {
	if next.Type == "punctuation" {
		textBuffer.WriteString(fmt.Sprintf("%s ", next.Alertantives[0].toString()))
	} else {
		if (next.StartTime - current.EndTime) < TRESHOLD {
			textBuffer.WriteString(fmt.Sprintf("%s ", next.Alertantives[0].toString()))
		}
	}
	return textBuffer
}

func getNewObject(items []TranscriptionItem) []TranscriptionItem {
	var newItems []TranscriptionItem

	var textBuffer bytes.Buffer
	for i := 0; i < len(items); i++ {
		item := items[i]
		var newItem TranscriptionItem
		if item.Type == "punctuation" {
			textBuffer.WriteString(item.Alertantives[0].toString())
		} else {
			newItem.StartTime = item.EndTime
			previusWithTime := item
			textBuffer.WriteString(fmt.Sprintf("%s ", item.Alertantives[0].toString()))
			for k := i + 1; k < len(items); k++ {
				currentItem := items[k]
				if currentItem.Type == "punctuation" {
					textBuffer.WriteString(fmt.Sprintf("%s", currentItem.Alertantives[0].toString()))
					if k == len(items)-1 {
						var newAlternatives Aleternative
						newItem.EndTime = previusWithTime.EndTime
						newAlternatives.Content = textBuffer.String()
						newItem.Alertantives = append(newItem.Alertantives, newAlternatives)
						newItems = append(newItems, newItem)
						textBuffer.Reset()
						i = k - 1
						break
					} else if currentItem.Alertantives[0].toString() == "." {
						var newAlternatives Aleternative
						newItem.EndTime = previusWithTime.EndTime
						newAlternatives.Content = textBuffer.String()
						newItem.Alertantives = append(newItem.Alertantives, newAlternatives)
						newItems = append(newItems, newItem)
						textBuffer.Reset()
						i = k
						break
					}
				} else {
					if (currentItem.StartTime - previusWithTime.EndTime) < TRESHOLD {
						textBuffer.WriteString(fmt.Sprintf("%s ", currentItem.Alertantives[0].toString()))
						previusWithTime = currentItem
					} else {
						var newAlternatives Aleternative
						newItem.EndTime = currentItem.StartTime
						newAlternatives.Content = textBuffer.String()
						newItem.Alertantives = append(newItem.Alertantives, newAlternatives)
						newItems = append(newItems, newItem)
						textBuffer.Reset()
						i = k - 1
						break
					}
				}

			}
		}
	}
	return newItems
}

func main() {

	trasncription := getTranscription()

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

	w := bufio.NewWriter(f)
	_, err = w.WriteString(fmt.Sprintf("WEBVTT\n\n"))
	for _, item := range items {
		_, err = w.WriteString(item.toString())
		_, err = w.WriteString("\n")
	}

	w.Flush()
}
