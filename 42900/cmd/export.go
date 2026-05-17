package cmd

import (
	"fmt"
	"os"
	"passman/entry"
	"passman/importexport"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export entries to file",
	Long:  "Export entries to CSV, JSON, or encrypted JSON file.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ensureStoreInitialized()
		ensureVaultUnlocked(cmd)
	},
}

var exportCsvCmd = &cobra.Command{
	Use:   "csv [file]",
	Short: "Export entries to CSV file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		entryManager := entry.NewManager(vaultManager, store)

		count, err := importexport.ExportToCSV(filePath, entryManager, vaultManager)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting entries: %v\n", err)
			return
		}

		fmt.Printf("Successfully exported %d entries to %s\n", count, filePath)
		fmt.Println("WARNING: This file contains unencrypted passwords!")
	},
}

var exportJsonCmd = &cobra.Command{
	Use:   "json [file]",
	Short: "Export entries to JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		entryManager := entry.NewManager(vaultManager, store)

		count, err := importexport.ExportToJSON(filePath, entryManager, vaultManager)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting entries: %v\n", err)
			return
		}

		fmt.Printf("Successfully exported %d entries to %s\n", count, filePath)
		fmt.Println("WARNING: This file contains unencrypted passwords!")
	},
}

var exportEncryptedCmd = &cobra.Command{
	Use:   "encrypted [file]",
	Short: "Export entries to encrypted JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		password, err := readPassword("Enter export password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		confirmPassword, err := readPassword("Confirm export password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		if password != confirmPassword {
			fmt.Fprintln(os.Stderr, "Passwords do not match")
			return
		}

		if len(password) < 8 {
			fmt.Fprintln(os.Stderr, "Export password must be at least 8 characters")
			return
		}

		entryManager := entry.NewManager(vaultManager, store)

		count, err := importexport.ExportEncrypted(filePath, password, entryManager, vaultManager)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting entries: %v\n", err)
			return
		}

		fmt.Printf("Successfully exported %d entries to %s (encrypted)\n", count, filePath)
	},
}

func init() {
	exportCmd.PersistentFlags().StringP("vault", "v", "", "Vault name to use (defaults to last active vault)")
	exportCmd.AddCommand(exportCsvCmd)
	exportCmd.AddCommand(exportJsonCmd)
	exportCmd.AddCommand(exportEncryptedCmd)
	rootCmd.AddCommand(exportCmd)
}
