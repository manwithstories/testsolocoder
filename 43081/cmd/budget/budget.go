package cmd

import (
	"fmt"

	"finance-tracker/internal/budget"

	"github.com/spf13/cobra"
)

var BudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage budgets",
	Long:  `Set, list, and track budgets for categories`,
}

var budgetSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a budget for a category",
	Run: func(cmd *cobra.Command, args []string) {
		categoryID, _ := cmd.Flags().GetInt64("category")
		amount, _ := cmd.Flags().GetFloat64("amount")
		month, _ := cmd.Flags().GetString("month")
		year, _ := cmd.Flags().GetInt("year")

		b, err := budget.Set(categoryID, amount, month, year)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Budget set successfully:\n")
		fmt.Printf("  Category: %d\n", b.CategoryID)
		fmt.Printf("  Period:   %s %d\n", b.Month, b.Year)
		fmt.Printf("  Amount:   %.2f\n", b.Amount)
	},
}

var budgetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all budgets",
	Run: func(cmd *cobra.Command, args []string) {
		month, _ := cmd.Flags().GetString("month")
		year, _ := cmd.Flags().GetInt("year")

		budgets, err := budget.List(month, year)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(budgets) == 0 {
			fmt.Println("No budgets found.")
			return
		}

		fmt.Printf("%-4s %-20s %-10s %-12s %-12s %-12s\n", "ID", "Category", "Period", "Budget", "Spent", "Remaining")
		fmt.Println("---------------------------------------------------------------------------------")
		for _, b := range budgets {
			period := fmt.Sprintf("%s %d", b.Month, b.Year)
			fmt.Printf("%-4d %-20s %-10s %12.2f %12.2f %12.2f",
				b.ID, b.CategoryName, period, b.Amount, b.Spent, b.Remaining)
			if b.Remaining < 0 {
				fmt.Print("  [OVERSPENT!]")
			}
			fmt.Println()
		}
	},
}

var budgetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a budget",
	Run: func(cmd *cobra.Command, args []string) {
		categoryID, _ := cmd.Flags().GetInt64("category")
		month, _ := cmd.Flags().GetString("month")
		year, _ := cmd.Flags().GetInt("year")

		err := budget.Delete(categoryID, month, year)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Budget deleted successfully\n")
	},
}

var budgetCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for overspent budgets",
	Run: func(cmd *cobra.Command, args []string) {
		month, _ := cmd.Flags().GetString("month")
		year, _ := cmd.Flags().GetInt("year")

		warnings := budget.GetWarnings(month, year)
		if len(warnings) == 0 {
			fmt.Println("All budgets are within limits!")
			return
		}

		for _, w := range warnings {
			fmt.Println(w)
		}
	},
}

func init() {
	budgetSetCmd.Flags().Int64P("category", "c", 0, "Category ID")
	budgetSetCmd.Flags().Float64P("amount", "a", 0, "Budget amount")
	budgetSetCmd.Flags().StringP("month", "m", "", "Month (01-12)")
	budgetSetCmd.Flags().IntP("year", "y", 0, "Year")
	budgetSetCmd.MarkFlagRequired("category")
	budgetSetCmd.MarkFlagRequired("amount")
	budgetSetCmd.MarkFlagRequired("month")
	budgetSetCmd.MarkFlagRequired("year")

	budgetListCmd.Flags().StringP("month", "m", "", "Filter by month (01-12)")
	budgetListCmd.Flags().IntP("year", "y", 0, "Filter by year")

	budgetDeleteCmd.Flags().Int64P("category", "c", 0, "Category ID")
	budgetDeleteCmd.Flags().StringP("month", "m", "", "Month (01-12)")
	budgetDeleteCmd.Flags().IntP("year", "y", 0, "Year")
	budgetDeleteCmd.MarkFlagRequired("category")
	budgetDeleteCmd.MarkFlagRequired("month")
	budgetDeleteCmd.MarkFlagRequired("year")

	budgetCheckCmd.Flags().StringP("month", "m", "", "Month (01-12)")
	budgetCheckCmd.Flags().IntP("year", "y", 0, "Year")

	BudgetCmd.AddCommand(budgetSetCmd)
	BudgetCmd.AddCommand(budgetListCmd)
	BudgetCmd.AddCommand(budgetDeleteCmd)
	BudgetCmd.AddCommand(budgetCheckCmd)
}
