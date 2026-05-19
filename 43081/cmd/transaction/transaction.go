package cmd

import (
	"fmt"
	"strconv"
	"time"

	"finance-tracker/internal/budget"
	"finance-tracker/internal/transaction"

	"github.com/spf13/cobra"
)

var TransactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Manage transactions",
	Long:  `Add, list, delete transactions and import from CSV`,
}

var transactionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new transaction",
	Run: func(cmd *cobra.Command, args []string) {
		transType, _ := cmd.Flags().GetString("type")
		amount, _ := cmd.Flags().GetFloat64("amount")
		categoryID, _ := cmd.Flags().GetInt64("category")
		accountID, _ := cmd.Flags().GetInt64("account")
		description, _ := cmd.Flags().GetString("desc")
		dateStr, _ := cmd.Flags().GetString("date")

		date := time.Now()
		if dateStr != "" {
			var err error
			date, err = transaction.ValidateDate(dateStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		t, err := transaction.Add(transType, amount, categoryID, accountID, description, date)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Transaction added successfully:\n")
		fmt.Printf("  ID:        %d\n", t.ID)
		fmt.Printf("  Type:      %s\n", t.Type)
		fmt.Printf("  Amount:    %.2f\n", t.Amount)
		fmt.Printf("  Category:  %s\n", t.CategoryName)
		fmt.Printf("  Account:   %s\n", t.AccountName)
		fmt.Printf("  Date:      %s\n", t.Date.Format("2006-01-02"))

		month := fmt.Sprintf("%02d", date.Month())
		year := date.Year()
		warnings := budget.GetWarnings(month, year)
		for _, w := range warnings {
			fmt.Println(w)
		}
	},
}

var transactionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions",
	Run: func(cmd *cobra.Command, args []string) {
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		categoryID, _ := cmd.Flags().GetInt64("category")
		accountID, _ := cmd.Flags().GetInt64("account")
		transType, _ := cmd.Flags().GetString("type")

		var start, end time.Time
		var err error

		if startStr != "" {
			start, err = transaction.ValidateDate(startStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
		if endStr != "" {
			end, err = transaction.ValidateDate(endStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		transactions, err := transaction.List(start, end, categoryID, accountID, transType)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(transactions) == 0 {
			fmt.Println("No transactions found.")
			return
		}

		fmt.Printf("%-4s %-10s %-7s %-10s %-15s %-15s %s\n", "ID", "Date", "Type", "Amount", "Category", "Account", "Description")
		fmt.Println("-------------------------------------------------------------------------------------------")
		for _, t := range transactions {
			fmt.Printf("%-4d %-10s %-7s %10.2f %-15s %-15s %s\n",
				t.ID, t.Date.Format("2006-01-02"), t.Type, t.Amount, t.CategoryName, t.AccountName, t.Description)
		}
	},
}

var transactionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a transaction",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: transaction ID required")
			return
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("Error: invalid transaction ID\n")
			return
		}

		err = transaction.Delete(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Transaction %d deleted successfully\n", id)
	},
}

var transactionImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import transactions from CSV",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")

		count, err := transaction.ImportCSV(filePath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Successfully imported %d transactions\n", count)
	},
}

func init() {
	transactionAddCmd.Flags().StringP("type", "t", "", "Transaction type (income or expense)")
	transactionAddCmd.Flags().Float64P("amount", "a", 0, "Transaction amount")
	transactionAddCmd.Flags().Int64P("category", "c", 0, "Category ID")
	transactionAddCmd.Flags().Int64P("account", "A", 0, "Account ID")
	transactionAddCmd.Flags().StringP("desc", "d", "", "Description")
	transactionAddCmd.Flags().String("date", "", "Transaction date (YYYY-MM-DD)")
	transactionAddCmd.MarkFlagRequired("type")
	transactionAddCmd.MarkFlagRequired("amount")
	transactionAddCmd.MarkFlagRequired("category")
	transactionAddCmd.MarkFlagRequired("account")

	transactionListCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	transactionListCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
	transactionListCmd.Flags().Int64P("category", "c", 0, "Filter by category ID")
	transactionListCmd.Flags().Int64P("account", "a", 0, "Filter by account ID")
	transactionListCmd.Flags().StringP("type", "t", "", "Filter by type (income or expense)")

	transactionImportCmd.Flags().StringP("file", "f", "", "CSV file path")
	transactionImportCmd.MarkFlagRequired("file")

	TransactionCmd.AddCommand(transactionAddCmd)
	TransactionCmd.AddCommand(transactionListCmd)
	TransactionCmd.AddCommand(transactionDeleteCmd)
	TransactionCmd.AddCommand(transactionImportCmd)
}
