package crawler

import (
	"UniCrawler/cmd/util"
	"fmt"
	"testing"
)

func TestCrawler(t *testing.T) {
	fmt.Println("Starting...")
	util.InitDatabase()
	DefaultInit()
	Start()
}

func TestStart(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start()
		})
	}
}
