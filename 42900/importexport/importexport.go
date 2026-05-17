package importexport

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"passman/crypto"
	"passman/entry"
	"passman/vault"
	"time"
)

type CSVEntry struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Tags     string `json:"tags"`
}

type EncryptedExport struct {
	Version   int       `json:"version"`
	Salt      []byte    `json:"salt"`
	Nonce     []byte    `json:"nonce"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

func ImportFromCSV(filePath string, entryManager *entry.Manager, vaultManager *vault.Manager) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	nameIdx := -1
	usernameIdx := -1
	passwordIdx := -1
	urlIdx := -1
	tagsIdx := -1

	for i, h := range headers {
		switch h {
		case "name", "Name", "title", "Title":
			nameIdx = i
		case "username", "Username", "user", "User":
			usernameIdx = i
		case "password", "Password":
			passwordIdx = i
		case "url", "URL", "website", "Website":
			urlIdx = i
		case "tags", "Tags", "label", "Label":
			tagsIdx = i
		}
	}

	if nameIdx == -1 || passwordIdx == -1 {
		return 0, errors.New("CSV must contain 'name' and 'password' columns")
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("failed to read CSV record: %w", err)
		}

		csvEntry := CSVEntry{}
		if nameIdx >= 0 && nameIdx < len(record) {
			csvEntry.Name = record[nameIdx]
		}
		if usernameIdx >= 0 && usernameIdx < len(record) {
			csvEntry.Username = record[usernameIdx]
		}
		if passwordIdx >= 0 && passwordIdx < len(record) {
			csvEntry.Password = record[passwordIdx]
		}
		if urlIdx >= 0 && urlIdx < len(record) {
			csvEntry.URL = record[urlIdx]
		}
		if tagsIdx >= 0 && tagsIdx < len(record) {
			csvEntry.Tags = record[tagsIdx]
		}

		if csvEntry.Name == "" || csvEntry.Password == "" {
			continue
		}

		_, err = entryManager.Add(csvEntry.Name, csvEntry.Username, csvEntry.Password, csvEntry.URL, csvEntry.Tags)
		if err != nil {
			return count, fmt.Errorf("failed to add entry '%s': %w", csvEntry.Name, err)
		}
		count++
	}

	return count, nil
}

func ExportToCSV(filePath string, entryManager *entry.Manager, vaultManager *vault.Manager) (int, error) {
	entries, err := entryManager.List()
	if err != nil {
		return 0, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"name", "username", "password", "url", "tags"}
	if err := writer.Write(headers); err != nil {
		return 0, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, e := range entries {
		decryptedPassword, err := vaultManager.DecryptPassword(e.Password)
		if err != nil {
			return 0, fmt.Errorf("failed to decrypt password for entry '%s': %w", e.Name, err)
		}

		record := []string{
			e.Name,
			e.Username,
			decryptedPassword,
			e.URL,
			e.Tags,
		}
		if err := writer.Write(record); err != nil {
			return 0, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return len(entries), nil
}

func ExportEncrypted(filePath string, password string, entryManager *entry.Manager, vaultManager *vault.Manager) (int, error) {
	entries, err := entryManager.List()
	if err != nil {
		return 0, err
	}

	exportEntries := make([]CSVEntry, 0, len(entries))
	for _, e := range entries {
		decryptedPassword, err := vaultManager.DecryptPassword(e.Password)
		if err != nil {
			return 0, fmt.Errorf("failed to decrypt password for entry '%s': %w", e.Name, err)
		}

		exportEntries = append(exportEntries, CSVEntry{
			Name:     e.Name,
			Username: e.Username,
			Password: decryptedPassword,
			URL:      e.URL,
			Tags:     e.Tags,
		})
	}

	jsonData, err := json.Marshal(exportEntries)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal entries: %w", err)
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return 0, err
	}

	key := crypto.DeriveKey(password, salt)

	encryptedData, err := crypto.Encrypt(jsonData, key)
	if err != nil {
		return 0, err
	}

	export := EncryptedExport{
		Version:   1,
		Salt:      salt,
		Nonce:     encryptedData.Nonce,
		Data:      encryptedData.Encode(),
		CreatedAt: time.Now(),
	}

	exportJSON, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(filePath, exportJSON, 0600); err != nil {
		return 0, fmt.Errorf("failed to write file: %w", err)
	}

	return len(entries), nil
}

func ImportEncrypted(filePath string, password string, entryManager *entry.Manager) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	var export EncryptedExport
	if err := json.Unmarshal(data, &export); err != nil {
		return 0, fmt.Errorf("failed to parse export file: %w", err)
	}

	if export.Version != 1 {
		return 0, fmt.Errorf("unsupported export version: %d", export.Version)
	}

	key := crypto.DeriveKey(password, export.Salt)

	encData, err := crypto.DecodeEncodedData(export.Data)
	if err != nil {
		return 0, err
	}

	plaintext, err := crypto.Decrypt(encData, key)
	if err != nil {
		return 0, errors.New("invalid password or corrupted data")
	}

	var entries []CSVEntry
	if err := json.Unmarshal(plaintext, &entries); err != nil {
		return 0, fmt.Errorf("failed to parse decrypted entries: %w", err)
	}

	count := 0
	for _, e := range entries {
		_, err = entryManager.Add(e.Name, e.Username, e.Password, e.URL, e.Tags)
		if err != nil {
			return count, fmt.Errorf("failed to add entry '%s': %w", e.Name, err)
		}
		count++
	}

	return count, nil
}

func ImportFromBrowserCSV(filePath string, entryManager *entry.Manager, vaultManager *vault.Manager) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	nameIdx, hasName := headerMap["name"]
	urlIdx, hasURL := headerMap["url"]
	usernameIdx, hasUsername := headerMap["username"]
	passwordIdx, hasPassword := headerMap["password"]

	if !hasName || !hasPassword {
		return 0, errors.New("CSV must contain 'name' and 'password' columns")
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("failed to read CSV record: %w", err)
		}

		csvEntry := CSVEntry{
			Name:     getField(record, nameIdx),
			Password: getField(record, passwordIdx),
		}

		if hasUsername {
			csvEntry.Username = getField(record, usernameIdx)
		}
		if hasURL {
			csvEntry.URL = getField(record, urlIdx)
		}

		if csvEntry.Name == "" || csvEntry.Password == "" {
			continue
		}

		_, err = entryManager.Add(csvEntry.Name, csvEntry.Username, csvEntry.Password, csvEntry.URL, csvEntry.Tags)
		if err != nil {
			return count, fmt.Errorf("failed to add entry '%s': %w", csvEntry.Name, err)
		}
		count++
	}

	return count, nil
}

func getField(record []string, idx int) string {
	if idx >= 0 && idx < len(record) {
		return record[idx]
	}
	return ""
}

func ExportToJSON(filePath string, entryManager *entry.Manager, vaultManager *vault.Manager) (int, error) {
	entries, err := entryManager.List()
	if err != nil {
		return 0, err
	}

	type ExportEntry struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		URL       string `json:"url"`
		Tags      string `json:"tags"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	exportEntries := make([]ExportEntry, 0, len(entries))
	for _, e := range entries {
		decryptedPassword, err := vaultManager.DecryptPassword(e.Password)
		if err != nil {
			return 0, fmt.Errorf("failed to decrypt password for entry '%s': %w", e.Name, err)
		}

		exportEntries = append(exportEntries, ExportEntry{
			ID:        e.ID,
			Name:      e.Name,
			Username:  e.Username,
			Password:  decryptedPassword,
			URL:       e.URL,
			Tags:      e.Tags,
			CreatedAt: e.CreatedAt.Format(time.RFC3339),
			UpdatedAt: e.UpdatedAt.Format(time.RFC3339),
		})
	}

	jsonData, err := json.MarshalIndent(exportEntries, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("failed to marshal entries: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0600); err != nil {
		return 0, fmt.Errorf("failed to write file: %w", err)
	}

	return len(entries), nil
}
