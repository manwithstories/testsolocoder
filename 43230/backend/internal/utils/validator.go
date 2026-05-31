package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IsSTLFile(reader io.Reader) (bool, error) {
	buf := make([]byte, 80)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	if n < 5 {
		return false, nil
	}

	header := strings.TrimSpace(string(buf[:min(n, 80)]))
	if strings.HasPrefix(header, "solid") {
		return true, nil
	}

	if n >= 84 {
		return true, nil
	}

	return false, nil
}

func IsOBJFile(reader io.Reader) (bool, error) {
	buf := make([]byte, 100)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	content := string(buf[:n])
	if strings.Contains(content, "v ") || strings.Contains(content, "vt ") ||
		strings.Contains(content, "vn ") || strings.Contains(content, "f ") {
		return true, nil
	}

	return false, nil
}

func ValidateModelFile(filename string, reader io.Reader) (string, error) {
	ext := strings.ToLower(strings.TrimPrefix(getFileExt(filename), "."))

	switch ext {
	case "stl":
		valid, err := IsSTLFile(reader)
		if err != nil {
			return "", fmt.Errorf("error validating STL file: %w", err)
		}
		if !valid {
			return "", fmt.Errorf("invalid STL file format")
		}
		return "stl", nil
	case "obj":
		valid, err := IsOBJFile(reader)
		if err != nil {
			return "", fmt.Errorf("error validating OBJ file: %w", err)
		}
		if !valid {
			return "", fmt.Errorf("invalid OBJ file format")
		}
		return "obj", nil
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
}

func getFileExt(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

func CalculateSTLVolume(reader io.Reader) (float64, error) {
	return 0.0, fmt.Errorf("volume calculation requires 3D processing library")
}

func EstimatePrintTime(volume float64, quality string) float64 {
	speedFactor := map[string]float64{
		"draft":    1.5,
		"standard": 1.0,
		"high":     0.7,
		"ultra":    0.5,
	}

	baseHours := volume * 0.001
	factor := speedFactor[quality]
	if factor == 0 {
		factor = 1.0
	}

	return baseHours / factor
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GenerateToken() string {
	h := md5.Sum([]byte(uuid.New().String() + time.Now().String()))
	return hex.EncodeToString(h[:])
}

func ToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
