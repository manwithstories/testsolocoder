package cmd

import (
	"fmt"
	"os"
	"passman/entry"
	"passman/importexport"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import entries from file",
	Long:  "Import entries from CSV or encrypted JSON file.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ensureStoreInitialized()
		ensureVaultUnlocked(cmd)
	},
}

var importCsvCmd = &cobra.Command{
	Use:   "csv [file]",
	Short: "Import entries from CSV file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		browser, _ := cmd.Flags().GetBool("browser")

		entryManager := entry.NewManager(vaultManager, store)

		var count int
		var err error

		if browser {
			count, err = importexport.ImportFromBrowserCSV(filePath, entryManager, vaultManager)
		} else {
			count, err = importexport.ImportFromCSV(filePath, entryManager, vaultManager)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing entries: %v\n", err)
			return
		}

		fmt.Printf("Successfully imported %d entries\n", count)
	},
}

var importEncryptedCmd = &cobra.Command{
	Use:   "encrypted [file]",
	Short: "Import entries from encrypted JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		password, err := readPassword("Enter export password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		entryManager := entry.NewManager(vaultManager, store)

		count, err := importexport.ImportEncrypted(filePath, password, entryManager)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing entries: %v\n", err)
			return
		}

		fmt.Printf("Successfully imported %d entries\n", count)
	},
}

func init() {
	importCsvCmd.Flags().Bool("browser", false, "Import from browser-exported CSV format")
	importCmd.PersistentFlags().StringP("vault", "v", "", "Vault name to use (defaults to last active vault)")
	importCmd.AddCommand(importCsvCmd)
	importCmd.AddCommand(importEncryptedCmd)
	rootCmd.AddCommand(importCmd)
}
