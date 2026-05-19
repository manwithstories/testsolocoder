package cmd

import (
	"fmt"
	"os"

	accountcmd "finance-tracker/cmd/account"
	backupcmd "finance-tracker/cmd/backup"
	budgetcmd "finance-tracker/cmd/budget"
	categorycmd "finance-tracker/cmd/category"
	statscmd "finance-tracker/cmd/stats"
	transactioncmd "finance-tracker/cmd/transaction"
	"finance-tracker/internal/backup"
	"finance-tracker/internal/config"
	"finance-tracker/internal/logger"
	"finance-tracker/internal/storage"

	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	rootCmd    = &cobra.Command{
		Use:   "finance",
		Short: "Personal finance tracker CLI tool",
		Long:  `A personal finance tracking tool built with Go and Cobra for managing accounts, transactions, budgets, and generating financial reports.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			cfg, err = config.Load()
			if err != nil {
				fmt.Printf("Failed to load config: %v\n", err)
				os.Exit(1)
			}

			if err := logger.Init(cfg); err != nil {
				fmt.Printf("Failed to initialize logger: %v\n", err)
				os.Exit(1)
			}

			if err := storage.Init(cfg); err != nil {
				fmt.Printf("Failed to initialize database: %v\n", err)
				os.Exit(1)
			}

			backup.CheckAndAutoBackup(cfg)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			storage.Close()
			logger.Close()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(accountcmd.AccountCmd)
	rootCmd.AddCommand(categorycmd.CategoryCmd)
	rootCmd.AddCommand(transactioncmd.TransactionCmd)
	rootCmd.AddCommand(budgetcmd.BudgetCmd)
	rootCmd.AddCommand(statscmd.StatsCmd)
	rootCmd.AddCommand(backupcmd.BackupCmd)
}
