package cmd

import (
	"fmt"
	"time"

	"finance-tracker/internal/stats"
	"finance-tracker/internal/transaction"

	"github.com/spf13/cobra"
)

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View statistics and reports",
	Long:  `Generate financial reports and statistics by time range, category, or account`,
}

var statsSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show financial summary",
	Run: func(cmd *cobra.Command, args []string) {
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		accountID, _ := cmd.Flags().GetInt64("account")

		var start, end time.Time
		var err error

		if startStr != "" {
			start, err = transaction.ValidateDate(startStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else {
			start = time.Now().AddDate(0, 0, -30)
		}
		if endStr != "" {
			end, err = transaction.ValidateDate(endStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else {
			end = time.Now()
		}

		summary, err := stats.GetSummary(start, end, accountID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("========================================")
		fmt.Println("FINANCIAL SUMMARY")
		fmt.Println("========================================")
		fmt.Printf("Period: %s to %s\n\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
		fmt.Printf("Total Income:  %15.2f\n", summary.TotalIncome)
		fmt.Printf("Total Expense: %15.2f\n", summary.TotalExpense)
		fmt.Println("----------------------------------------")
		netColor := ""
		if summary.NetProfit >= 0 {
			netColor = "+"
		}
		fmt.Printf("Net Profit:    %s%14.2f\n", netColor, summary.NetProfit)
	},
}

var statsCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Show statistics by category",
	Run: func(cmd *cobra.Command, args []string) {
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		transType, _ := cmd.Flags().GetString("type")
		accountID, _ := cmd.Flags().GetInt64("account")

		var start, end time.Time
		var err error

		if startStr != "" {
			start, err = transaction.ValidateDate(startStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else {
			start = time.Now().AddDate(0, 0, -30)
		}
		if endStr != "" {
			end, err = transaction.ValidateDate(endStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else {
			end = time.Now()
		}

		if transType == "" {
			transType = "expense"
		}

		categoryStats, err := stats.GetByCategory(start, end, transType, accountID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(categoryStats) == 0 {
			fmt.Printf("No %s transactions found in this period.\n", transType)
			return
		}

		fmt.Printf("%s BY CATEGORY (%s to %s)\n", transType, start.Format("2006-01-02"), end.Format("2006-01-02"))
		fmt.Println("----------------------------------------")
		fmt.Printf("%-20s %-12s %-8s\n", "Category", "Amount", "Percent")
		fmt.Println("----------------------------------------")
		for _, s := range categoryStats {
			bar := ""
			barLen := int(s.Percentage / 5)
			for i := 0; i < barLen; i++ {
				bar += "="
			}
			fmt.Printf("%-20s %12.2f %6.1f%%  %s\n", s.CategoryName, s.Amount, s.Percentage, bar)
		}
	},
}

var statsReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a detailed report",
	Run: func(cmd *cobra.Command, args []string) {
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		accountID, _ := cmd.Flags().GetInt64("account")
		output, _ := cmd.Flags().GetString("output")

		var start, end time.Time
		var err error

		if startStr != "" {
			start, err = transaction.ValidateDate(startStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else {
			start = time.Now().AddDate(0, 0, -30)
		}
		if endStr != "" {
			end, err = transaction.ValidateDate(endStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else {
			end = time.Now()
		}

		err = stats.GenerateReport(start, end, accountID, output)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Report generated: %s\n", output)
	},
}

func init() {
	statsSummaryCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	statsSummaryCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
	statsSummaryCmd.Flags().Int64P("account", "a", 0, "Filter by account ID")

	statsCategoryCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	statsCategoryCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
	statsCategoryCmd.Flags().StringP("type", "t", "expense", "Transaction type (income or expense)")
	statsCategoryCmd.Flags().Int64P("account", "a", 0, "Filter by account ID")

	statsReportCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	statsReportCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
	statsReportCmd.Flags().Int64P("account", "a", 0, "Filter by account ID")
	statsReportCmd.Flags().StringP("output", "o", "report.txt", "Output file path (.txt or .csv)")

	StatsCmd.AddCommand(statsSummaryCmd)
	StatsCmd.AddCommand(statsCategoryCmd)
	StatsCmd.AddCommand(statsReportCmd)
}
