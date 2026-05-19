package main

import (
	"os"
	"passman/clipboard"
	"passman/cmd"
	"strconv"
)

func main() {
	if len(os.Args) >= 4 && os.Args[1] == clipboard.ClearClipboardCommand {
		seconds, _ := strconv.Atoi(os.Args[2])
		expectedText := os.Args[3]
		_ = clipboard.ClearClipboardAfterDelay(seconds, expectedText)
		return
	}

	cmd.Execute()
}
