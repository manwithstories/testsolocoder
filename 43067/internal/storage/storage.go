package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"snippetbox/internal/config"
	"snippetbox/internal/crypto"
	"snippetbox/internal/logger"
	"snippetbox/internal/models"
)

var ErrVaultNotFound = errors.New("vault not found")
var ErrSnippetNotFound = errors.New("snippet not found")
var ErrDuplicateSnippetID = errors.New("duplicate snippet ID")
var ErrEmptyVault = errors.New("vault is empty")

func getDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	return filepath.Join(homeDir, ".snippetbox", "vaults"), nil
}

func getVaultPath(vaultID string) (string, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, vaultID+".json"), nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func ListVaults() ([]models.Vault, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return []models.Vault{}, nil
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read vaults directory: %w", err)
	}

	var vaults []models.Vault
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			vaultID := file.Name()[:len(file.Name())-5]
			vault, err := LoadVault(vaultID)
			if err == nil {
				vaults = append(vaults, *vault)
			}
		}
	}

	logger.Info("Listed %d vaults", len(vaults))
	return vaults, nil
}

func LoadVault(vaultID string) (*models.Vault, error) {
	vaultPath, err := getVaultPath(vaultID)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return nil, ErrVaultNotFound
	}

	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	var vault models.Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, fmt.Errorf("failed to parse vault: %w", err)
	}

	logger.Debug("Loaded vault: %s", vaultID)
	return &vault, nil
}

func SaveVault(vault *models.Vault) error {
	vaultPath, err := getVaultPath(vault.ID)
	if err != nil {
		return err
	}

	dataDir := filepath.Dir(vaultPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create vaults directory: %w", err)
	}

	vault.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	if err := os.WriteFile(vaultPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	logger.Info("Saved vault: %s", vault.ID)
	return nil
}

func CreateVault(name string) (*models.Vault, error) {
	vault := &models.Vault{
		ID:        generateID(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Snippets:  []models.Snippet{},
	}

	if err := SaveVault(vault); err != nil {
		return nil, err
	}

	logger.Info("Created vault: %s (%s)", vault.Name, vault.ID)
	return vault, nil
}

func DeleteVault(vaultID string) error {
	vaultPath, err := getVaultPath(vaultID)
	if err != nil {
		return err
	}

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return ErrVaultNotFound
	}

	if err := os.Remove(vaultPath); err != nil {
		return fmt.Errorf("failed to delete vault file: %w", err)
	}

	cfg, err := config.Load()
	if err == nil && cfg.DefaultVault == vaultID {
		remainingVaults, listErr := ListVaults()
		if listErr == nil && len(remainingVaults) > 0 {
			cfg.DefaultVault = remainingVaults[0].ID
			logger.Info("Default vault changed to: %s", cfg.DefaultVault)
		} else {
			cfg.DefaultVault = ""
			logger.Info("No vaults remaining, default vault cleared")
		}
		config.Save(cfg)
	}

	logger.Info("Deleted vault: %s", vaultID)
	return nil
}

func AddSnippet(vaultID string, snippet *models.Snippet) error {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return err
	}

	for _, s := range vault.Snippets {
		if s.ID == snippet.ID {
			return ErrDuplicateSnippetID
		}
	}

	snippet.ID = generateID()
	snippet.CreatedAt = time.Now()
	snippet.UpdatedAt = time.Now()

	if snippet.Encrypted {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config for encryption: %w", err)
		}
		if cfg.EncryptionKey == "" {
			return errors.New("encryption key not set in config")
		}
		encryptedContent, err := crypto.Encrypt(snippet.Content, cfg.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt snippet: %w", err)
		}
		snippet.Content = encryptedContent
	}

	vault.Snippets = append(vault.Snippets, *snippet)
	return SaveVault(vault)
}

func GetSnippet(vaultID, snippetID string, decrypt bool) (*models.Snippet, error) {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return nil, err
	}

	for _, s := range vault.Snippets {
		if s.ID == snippetID {
			snippet := &s
			if decrypt && snippet.Encrypted {
				cfg, err := config.Load()
				if err != nil {
					return nil, fmt.Errorf("failed to load config for decryption: %w", err)
				}
				if cfg.EncryptionKey == "" {
					return nil, errors.New("encryption key not set in config")
				}
				decryptedContent, err := crypto.Decrypt(snippet.Content, cfg.EncryptionKey)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt snippet: %w", err)
				}
				snippet.Content = decryptedContent
			}
			return snippet, nil
		}
	}

	return nil, ErrSnippetNotFound
}

func UpdateSnippet(vaultID string, snippet *models.Snippet) error {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return err
	}

	for i, s := range vault.Snippets {
		if s.ID == snippet.ID {
			snippet.UpdatedAt = time.Now()
			snippet.CreatedAt = s.CreatedAt

			if snippet.Encrypted {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("failed to load config for encryption: %w", err)
				}
				if cfg.EncryptionKey == "" {
					return errors.New("encryption key not set in config")
				}
				encryptedContent, err := crypto.Encrypt(snippet.Content, cfg.EncryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt snippet: %w", err)
				}
				snippet.Content = encryptedContent
			}

			vault.Snippets[i] = *snippet
			return SaveVault(vault)
		}
	}

	return ErrSnippetNotFound
}

func DeleteSnippet(vaultID, snippetID string) error {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return err
	}

	for i, s := range vault.Snippets {
		if s.ID == snippetID {
			vault.Snippets = append(vault.Snippets[:i], vault.Snippets[i+1:]...)
			return SaveVault(vault)
		}
	}

	return ErrSnippetNotFound
}

func ListSnippets(vaultID string) ([]models.Snippet, error) {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return nil, err
	}

	return vault.Snippets, nil
}

func GetAllSnippets(vaultID string, decrypt bool) ([]models.Snippet, error) {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return nil, err
	}

	if len(vault.Snippets) == 0 {
		return []models.Snippet{}, nil
	}

	if !decrypt {
		return vault.Snippets, nil
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var decryptedSnippets []models.Snippet
	for _, s := range vault.Snippets {
		if s.Encrypted && cfg.EncryptionKey != "" {
			decryptedContent, err := crypto.Decrypt(s.Content, cfg.EncryptionKey)
			if err == nil {
				s.Content = decryptedContent
			}
		}
		decryptedSnippets = append(decryptedSnippets, s)
	}

	return decryptedSnippets, nil
}
