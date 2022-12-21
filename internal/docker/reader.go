package docker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// output from docker client has to be read, otherwise we return before the process finishes
func printOnOneLine(reader io.Reader) {

	scanner := bufio.NewScanner(reader)
	var prevLineLength int
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		fmt.Printf("\r%s", strings.Repeat(" ", prevLineLength))
		fmt.Printf("\r%s", text)
		prevLineLength = len(text)
	}
	if prevLineLength != 0 {
		fmt.Println()
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("cannot read output: %v", err)
	}
}
