package cmd

import (
	"fmt"
	"strconv"
	"time"

	"finance-tracker/internal/account"

	"github.com/spf13/cobra"
)

var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	Long:  `Create, list, and manage your accounts (cash, bank, Alipay, etc.)`,
}

var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		accountType, _ := cmd.Flags().GetString("type")
		currency, _ := cmd.Flags().GetString("currency")
		balance, _ := cmd.Flags().GetFloat64("balance")

		acc, err := account.Create(name, accountType, currency, balance)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Account created successfully:\n")
		fmt.Printf("  ID:       %d\n", acc.ID)
		fmt.Printf("  Name:     %s\n", acc.Name)
		fmt.Printf("  Type:     %s\n", acc.Type)
		fmt.Printf("  Balance:  %.2f %s\n", acc.Balance, acc.Currency)
	},
}

var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := account.List()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts found.")
			return
		}

		fmt.Printf("%-4s %-15s %-10s %-15s %s\n", "ID", "Name", "Type", "Balance", "Currency")
		fmt.Println("------------------------------------------------------------")
		for _, a := range accounts {
			fmt.Printf("%-4d %-15s %-10s %15.2f %s\n", a.ID, a.Name, a.Type, a.Balance, a.Currency)
		}
	},
}

var accountTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer money between accounts",
	Run: func(cmd *cobra.Command, args []string) {
		fromID, _ := cmd.Flags().GetInt64("from")
		toID, _ := cmd.Flags().GetInt64("to")
		amount, _ := cmd.Flags().GetFloat64("amount")
		description, _ := cmd.Flags().GetString("desc")
		dateStr, _ := cmd.Flags().GetString("date")

		date := time.Now()
		if dateStr != "" {
			var err error
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				fmt.Printf("Error: invalid date format, use YYYY-MM-DD\n")
				return
			}
		}

		err := account.Transfer(fromID, toID, amount, description, date)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Transfer successful: %.2f from account %d to account %d\n", amount, fromID, toID)
	},
}

var accountDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an account",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: account ID required")
			return
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("Error: invalid account ID\n")
			return
		}

		err = account.Delete(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Account %d deleted successfully\n", id)
	},
}

func init() {
	accountCreateCmd.Flags().StringP("name", "n", "", "Account name")
	accountCreateCmd.Flags().StringP("type", "t", "cash", "Account type (cash, bank, alipay, wechat, credit, other)")
	accountCreateCmd.Flags().StringP("currency", "c", "CNY", "Currency code")
	accountCreateCmd.Flags().Float64P("balance", "b", 0, "Initial balance")
	accountCreateCmd.MarkFlagRequired("name")

	accountTransferCmd.Flags().Int64P("from", "f", 0, "Source account ID")
	accountTransferCmd.Flags().Int64P("to", "t", 0, "Target account ID")
	accountTransferCmd.Flags().Float64P("amount", "a", 0, "Transfer amount")
	accountTransferCmd.Flags().StringP("desc", "d", "", "Description")
	accountTransferCmd.Flags().String("date", "", "Transfer date (YYYY-MM-DD)")
	accountTransferCmd.MarkFlagRequired("from")
	accountTransferCmd.MarkFlagRequired("to")
	accountTransferCmd.MarkFlagRequired("amount")

	AccountCmd.AddCommand(accountCreateCmd)
	AccountCmd.AddCommand(accountListCmd)
	AccountCmd.AddCommand(accountTransferCmd)
	AccountCmd.AddCommand(accountDeleteCmd)
}
