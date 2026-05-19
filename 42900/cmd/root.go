package cmd

import (
	"fmt"
	"os"
	"passman/storage"
	"passman/vault"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	store       *storage.Store
	vaultManager *vault.Manager
	dbPath      string
)

func initStore() {
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		dbPath = filepath.Join(homeDir, ".passman", "passman.db")
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating database directory: %v\n", err)
		os.Exit(1)
	}

	var err error
	store, err = storage.NewStore(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}

	vaultManager = vault.NewManager(store)
}

func ensureStoreInitialized() {
	if store == nil {
		initStore()
	}
}

var rootCmd = &cobra.Command{
	Use:   "passman",
	Short: "A secure command-line password manager",
	Long: `passman is a secure command-line password manager that helps you store,
generate, and manage your passwords securely using AES-256-GCM encryption.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initStore()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if store != nil {
			store.Close()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "Path to the database file (default: ~/.passman/passman.db)")
}
