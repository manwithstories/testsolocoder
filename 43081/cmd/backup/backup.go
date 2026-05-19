package cmd

import (
	"fmt"

	"finance-tracker/internal/backup"
	"finance-tracker/internal/config"

	"github.com/spf13/cobra"
)

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage data backups",
	Long:  `Create, list, and restore database backups`,
}

var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a local backup",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load()
		path, err := backup.CreateLocal(cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Backup created successfully: %s\n", path)
	},
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	Run: func(cmd *cobra.Command, args []string) {
		backups, err := backup.List()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(backups) == 0 {
			fmt.Println("No backups found.")
			return
		}

		fmt.Printf("%-4s %-25s %-12s %s\n", "ID", "Created At", "Size", "Path")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, b := range backups {
			size := fmt.Sprintf("%.2f MB", float64(b.FileSize)/1024/1024)
			fmt.Printf("%-4d %-25s %-12s %s\n", b.ID, b.CreatedAt.Format("2006-01-02 15:04:05"), size, b.FilePath)
		}
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore from a backup",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")

		fmt.Println("WARNING: This will overwrite the current database!")
		fmt.Print("Type 'yes' to continue: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Restore cancelled.")
			return
		}

		cfg, _ := config.Load()
		err := backup.Restore(path, cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Database restored successfully from: %s\n", path)
	},
}

func init() {
	backupRestoreCmd.Flags().StringP("path", "p", "", "Backup file path")
	backupRestoreCmd.MarkFlagRequired("path")

	BackupCmd.AddCommand(backupCreateCmd)
	BackupCmd.AddCommand(backupListCmd)
	BackupCmd.AddCommand(backupRestoreCmd)
}
