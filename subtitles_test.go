package main

import (
	"testing"
)

func BenchmarkGetTranscription(b *testing.B) {
	getTranscription("./asrOutput.json")
}

func BenchmarkGetObject(B *testing.B) {
	trasncription := getTranscription("./asrOutput.json")
	getNewObject(trasncription.Results.Items)
}
