package cmd

import (
	"fmt"
	"strconv"

	"finance-tracker/internal/category"
	"finance-tracker/internal/models"

	"github.com/spf13/cobra"
)

var CategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage categories",
	Long:  `Create, list, and manage income/expense categories with support for hierarchical categories`,
}

var categoryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new category",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		categoryType, _ := cmd.Flags().GetString("type")
		parentID, _ := cmd.Flags().GetInt64("parent")

		var parentPtr *int64
		if parentID > 0 {
			parentPtr = &parentID
		}

		cat, err := category.Create(name, categoryType, parentPtr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Category created successfully:\n")
		fmt.Printf("  ID:   %d\n", cat.ID)
		fmt.Printf("  Name: %s\n", cat.Name)
		fmt.Printf("  Type: %s\n", cat.Type)
		if cat.ParentID != nil {
			fmt.Printf("  Parent ID: %d\n", *cat.ParentID)
		}
	},
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categories",
	Run: func(cmd *cobra.Command, args []string) {
		categoryType, _ := cmd.Flags().GetString("type")
		tree, _ := cmd.Flags().GetBool("tree")

		if tree {
			cats, err := category.ListTree(categoryType)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			printCategoryTree(cats, 0)
		} else {
			cats, err := category.List(categoryType)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if len(cats) == 0 {
				fmt.Println("No categories found.")
				return
			}

			fmt.Printf("%-4s %-20s %-10s %s\n", "ID", "Name", "Type", "Parent")
			fmt.Println("------------------------------------------------")
			for _, c := range cats {
				parent := "-"
				if c.ParentID != nil {
					parent = strconv.FormatInt(*c.ParentID, 10)
				}
				fmt.Printf("%-4d %-20s %-10s %s\n", c.ID, c.Name, c.Type, parent)
			}
		}
	},
}

var categoryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a category",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: category ID required")
			return
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("Error: invalid category ID\n")
			return
		}

		err = category.Delete(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Category %d deleted successfully\n", id)
	},
}

func printCategoryTree(cats []models.Category, level int) {
	for _, c := range cats {
		indent := ""
		for i := 0; i < level; i++ {
			indent += "  "
		}
		fmt.Printf("%s%-4d %s (%s)\n", indent, c.ID, c.Name, c.Type)
		printCategoryTree(c.Children, level+1)
	}
}

func init() {
	categoryCreateCmd.Flags().StringP("name", "n", "", "Category name")
	categoryCreateCmd.Flags().StringP("type", "t", "expense", "Category type (income or expense)")
	categoryCreateCmd.Flags().Int64P("parent", "p", 0, "Parent category ID")
	categoryCreateCmd.MarkFlagRequired("name")

	categoryListCmd.Flags().StringP("type", "t", "", "Filter by type (income or expense)")
	categoryListCmd.Flags().Bool("tree", false, "Display as tree structure")

	CategoryCmd.AddCommand(categoryCreateCmd)
	CategoryCmd.AddCommand(categoryListCmd)
	CategoryCmd.AddCommand(categoryDeleteCmd)
}
