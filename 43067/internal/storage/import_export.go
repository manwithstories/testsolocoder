package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"snippetbox/internal/logger"
	"snippetbox/internal/models"
)

func ExportVault(vaultID, outputPath string) error {
	vault, err := LoadVault(vaultID)
	if err != nil {
		return err
	}

	exportData := models.ExportData{
		Version: "1.0",
		Vaults:  []models.Vault{*vault},
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	logger.Info("Exported vault %s to %s", vaultID, outputPath)
	return nil
}

func ExportAllVaults(outputPath string) error {
	vaults, err := ListVaults()
	if err != nil {
		return err
	}

	exportData := models.ExportData{
		Version: "1.0",
		Vaults:  vaults,
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	logger.Info("Exported %d vaults to %s", len(vaults), outputPath)
	return nil
}

func ImportVaults(inputPath string, merge bool) ([]string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	var importData models.ExportData
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	var importedVaults []string

	for _, vault := range importData.Vaults {
		existingVault, err := LoadVault(vault.ID)
		if err == nil && merge {
			existingSnippetIDs := make(map[string]bool)
			for _, s := range existingVault.Snippets {
				existingSnippetIDs[s.ID] = true
			}

			for _, s := range vault.Snippets {
				if !existingSnippetIDs[s.ID] {
					existingVault.Snippets = append(existingVault.Snippets, s)
				}
			}

			if err := SaveVault(existingVault); err != nil {
				return importedVaults, fmt.Errorf("failed to save merged vault %s: %w", vault.ID, err)
			}
			importedVaults = append(importedVaults, vault.ID)
		} else {
			if err := SaveVault(&vault); err != nil {
				return importedVaults, fmt.Errorf("failed to save imported vault %s: %w", vault.ID, err)
			}
			importedVaults = append(importedVaults, vault.ID)
		}
	}

	logger.Info("Imported %d vaults from %s", len(importedVaults), inputPath)
	return importedVaults, nil
}
