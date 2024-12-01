package advent

import (
	"log"
	"os"
	"strings"
)

func Read(fileName string) []string {
	body, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	var lines []string

	for _, line := range strings.Split(string(body), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}
