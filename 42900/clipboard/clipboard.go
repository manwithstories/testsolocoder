package clipboard

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
)

const DefaultClearDelay = 30 * time.Second
const ClearClipboardCommand = "__clear_clipboard__"

func CopyToClipboard(text string, clearDelay time.Duration) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	if clearDelay <= 0 {
		clearDelay = DefaultClearDelay
	}

	tmpDir := os.TempDir()
	tmpFileName := fmt.Sprintf("passman-clipboard-%d-%s", os.Getpid(), randomString(8))
	tmpFilePath := filepath.Join(tmpDir, tmpFileName)

	if err := os.WriteFile(tmpFilePath, []byte(text), 0600); err != nil {
		fmt.Printf("Password copied to clipboard. WARNING: Failed to create secure temp file, clipboard will NOT be cleared automatically.\n")
		return nil
	}

	execPath, err := os.Executable()
	if err != nil {
		os.Remove(tmpFilePath)
		fmt.Printf("Password copied to clipboard. WARNING: Failed to start auto-clear process, clipboard will NOT be cleared automatically.\n")
		return nil
	}

	cmd := exec.Command(execPath, ClearClipboardCommand, strconv.Itoa(int(clearDelay.Seconds())), tmpFilePath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		os.Remove(tmpFilePath)
		fmt.Printf("Password copied to clipboard. WARNING: Failed to start auto-clear process, clipboard will NOT be cleared automatically.\n")
		return nil
	}

	if err := cmd.Process.Release(); err != nil {
		fmt.Printf("Password copied to clipboard. WARNING: Failed to detach auto-clear process.\n")
		return nil
	}

	fmt.Printf("Password copied to clipboard. Will be cleared in %v.\n", clearDelay)
	return nil
}

func ClearClipboardAfterDelay(seconds int, tempFilePath string) error {
	time.Sleep(time.Duration(seconds) * time.Second)

	text, err := os.ReadFile(tempFilePath)
	os.Remove(tempFilePath)

	if err != nil {
		return err
	}

	expectedText := string(text)

	current, err := clipboard.ReadAll()
	if err != nil {
		return err
	}

	if current == expectedText {
		return clipboard.WriteAll("")
	}

	return nil
}

func randomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "fallback"
	}
	return hex.EncodeToString(bytes)[:n]
}

func ReadFromClipboard() (string, error) {
	return clipboard.ReadAll()
}

func ClearClipboard() error {
	return clipboard.WriteAll("")
}
