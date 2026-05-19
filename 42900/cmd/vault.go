package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage password vaults",
	Long:  "Create, open, switch, and delete password vaults.",
}

var vaultCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		masterPassword, err := readPassword("Enter master password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		confirmPassword, err := readPassword("Confirm master password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		if masterPassword != confirmPassword {
			fmt.Fprintln(os.Stderr, "Passwords do not match")
			return
		}

		if len(masterPassword) < 8 {
			fmt.Fprintln(os.Stderr, "Master password must be at least 8 characters")
			return
		}

		_, err = vaultManager.CreateVault(name, masterPassword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating vault: %v\n", err)
			return
		}

		fmt.Printf("Vault '%s' created successfully!\n", name)
	},
}

var vaultOpenCmd = &cobra.Command{
	Use:   "open [name]",
	Short: "Open a vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		masterPassword, err := readPassword("Enter master password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		err = vaultManager.OpenVault(name, masterPassword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening vault: %v\n", err)
			return
		}

		fmt.Printf("Vault '%s' opened successfully!\n", name)
	},
}

var vaultCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close the current vault",
	Run: func(cmd *cobra.Command, args []string) {
		if vaultManager.GetCurrentVault() == nil {
			fmt.Println("No vault is open")
			return
		}

		vaultManager.CloseVault()
		fmt.Println("Vault closed")
	},
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vaults",
	Run: func(cmd *cobra.Command, args []string) {
		vaults, err := vaultManager.ListVaults()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing vaults: %v\n", err)
			return
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found")
			return
		}

		currentVault := vaultManager.GetCurrentVault()

		fmt.Println("Vaults:")
		for _, v := range vaults {
			prefix := "  "
			if currentVault != nil && currentVault.ID == v.ID {
				prefix = "* "
			}
			fmt.Printf("%s%s (created: %s)\n", prefix, v.Name, v.CreatedAt.Format("2006-01-02"))
		}
	},
}

var vaultDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		confirm, err := readConfirmation(fmt.Sprintf("Are you sure you want to delete vault '%s'? This cannot be undone. (y/N): ", name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading confirmation: %v\n", err)
			return
		}

		if !confirm {
			fmt.Println("Delete cancelled")
			return
		}

		err = vaultManager.DeleteVault(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting vault: %v\n", err)
			return
		}

		fmt.Printf("Vault '%s' deleted successfully\n", name)
	},
}

var vaultChangePwdCmd = &cobra.Command{
	Use:   "change-password",
	Short: "Change the master password for the current vault",
	Run: func(cmd *cobra.Command, args []string) {
		ensureStoreInitialized()
		ensureVaultUnlocked(cmd)

		oldPassword, err := readPassword("Enter current master password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		newPassword, err := readPassword("Enter new master password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		confirmPassword, err := readPassword("Confirm new master password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			return
		}

		if newPassword != confirmPassword {
			fmt.Fprintln(os.Stderr, "Passwords do not match")
			return
		}

		if len(newPassword) < 8 {
			fmt.Fprintln(os.Stderr, "New master password must be at least 8 characters")
			return
		}

		err = vaultManager.ChangeMasterPassword(oldPassword, newPassword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error changing master password: %v\n", err)
			return
		}

		fmt.Println("Master password changed successfully! All entries have been re-encrypted.")
	},
}

var vaultStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current vault status",
	Run: func(cmd *cobra.Command, args []string) {
		currentVault := vaultManager.GetCurrentVault()
		if currentVault == nil {
			fmt.Println("No vault is open")

			lastVault, err := vaultManager.GetLastActiveVault()
			if err == nil && lastVault != "" {
				fmt.Printf("Last active vault: %s\n", lastVault)
			}
			return
		}

		fmt.Printf("Current vault: %s\n", currentVault.Name)
		fmt.Printf("Created: %s\n", currentVault.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Last updated: %s\n", currentVault.UpdatedAt.Format("2006-01-02 15:04:05"))
	},
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

func readConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

func ensureVaultUnlocked(cmd *cobra.Command) {
	if vaultManager.GetCurrentVault() != nil {
		return
	}

	lastVault, err := vaultManager.GetLastActiveVault()
	if err != nil || lastVault == "" {
		fmt.Fprintln(os.Stderr, "No vault is open. Use 'passman vault open <name>' first.")
		os.Exit(1)
	}

	vaultName, _ := cmd.Flags().GetString("vault")
	if vaultName == "" {
		vaultName = lastVault
	}

	masterPassword, err := readPassword(fmt.Sprintf("Enter master password for vault '%s': ", vaultName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	err = vaultManager.OpenVault(vaultName, masterPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening vault: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	vaultChangePwdCmd.Flags().StringP("vault", "v", "", "Vault name to change password for (defaults to last active vault)")
	vaultCmd.AddCommand(vaultCreateCmd)
	vaultCmd.AddCommand(vaultOpenCmd)
	vaultCmd.AddCommand(vaultCloseCmd)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultDeleteCmd)
	vaultCmd.AddCommand(vaultChangePwdCmd)
	vaultCmd.AddCommand(vaultStatusCmd)
	rootCmd.AddCommand(vaultCmd)
}
