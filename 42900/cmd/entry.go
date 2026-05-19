package cmd

import (
	"fmt"
	"os"
	"passman/clipboard"
	"passman/entry"
	"passman/generator"
	"passman/storage"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	entryManager *entry.Manager
)

var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Manage password entries",
	Long:  "Add, list, search, update, and delete password entries.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ensureStoreInitialized()
		ensureVaultUnlocked(cmd)
		entryManager = entry.NewManager(vaultManager, store)
	},
}

var entryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password entry",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		tags, _ := cmd.Flags().GetString("tags")
		generate, _ := cmd.Flags().GetBool("generate")

		if name == "" {
			fmt.Fprintln(os.Stderr, "Name is required")
			return
		}

		if password == "" && !generate {
			fmt.Fprintln(os.Stderr, "Password is required or use --generate")
			return
		}

		if generate {
			genOpts := generator.NewDefaultOptions()
			genLength, _ := cmd.Flags().GetInt("length")
			if genLength > 0 {
				genOpts.Length = genLength
			}
			var err error
			password, err = generator.Generate(genOpts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
				return
			}
			fmt.Printf("Generated password: %s\n", password)
		}

		entry, err := entryManager.Add(name, username, password, url, tags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding entry: %v\n", err)
			return
		}

		fmt.Printf("Entry '%s' added successfully (ID: %d)\n", entry.Name, entry.ID)
	},
}

var entryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entries",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := entryManager.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing entries: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No entries found")
			return
		}

		printEntries(entries)
	},
}

var entrySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search entries by name or tag",
	Run: func(cmd *cobra.Command, args []string) {
		nameQuery, _ := cmd.Flags().GetString("name")
		tagQuery, _ := cmd.Flags().GetString("tag")

		if nameQuery == "" && tagQuery == "" {
			fmt.Fprintln(os.Stderr, "Please provide --name or --tag for searching")
			return
		}

		entries, err := entryManager.Search(nameQuery, tagQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching entries: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No matching entries found")
			return
		}

		printEntries(entries)
	},
}

var entryGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get entry details by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %v\n", err)
			return
		}

		showPassword, _ := cmd.Flags().GetBool("show-password")
		copyPassword, _ := cmd.Flags().GetBool("copy")

		entry, decryptedPassword, err := entryManager.Get(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting entry: %v\n", err)
			return
		}

		fmt.Printf("ID: %d\n", entry.ID)
		fmt.Printf("Name: %s\n", entry.Name)
		fmt.Printf("Username: %s\n", entry.Username)
		if showPassword {
			fmt.Printf("Password: %s\n", decryptedPassword)
		} else {
			fmt.Println("Password: ******** (use --show-password to reveal)")
		}
		fmt.Printf("URL: %s\n", entry.URL)
		fmt.Printf("Tags: %s\n", entry.Tags)
		fmt.Printf("Created: %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05"))

		if copyPassword {
			err := clipboard.CopyToClipboard(decryptedPassword, 30*time.Second)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
			}
		}
	},
}

var entryUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update an entry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %v\n", err)
			return
		}

		name, _ := cmd.Flags().GetString("name")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		tags, _ := cmd.Flags().GetString("tags")

		err = entryManager.Update(id, name, username, password, url, tags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating entry: %v\n", err)
			return
		}

		fmt.Println("Entry updated successfully")
	},
}

var entryDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an entry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %v\n", err)
			return
		}

		confirm, err := readConfirmation("Are you sure you want to delete this entry? (y/N): ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading confirmation: %v\n", err)
			return
		}

		if !confirm {
			fmt.Println("Delete cancelled")
			return
		}

		err = entryManager.Delete(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting entry: %v\n", err)
			return
		}

		fmt.Println("Entry deleted successfully")
	},
}

var entryCopyCmd = &cobra.Command{
	Use:   "copy [id]",
	Short: "Copy password to clipboard",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %v\n", err)
			return
		}

		_, decryptedPassword, err := entryManager.Get(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting entry: %v\n", err)
			return
		}

		err = clipboard.CopyToClipboard(decryptedPassword, 30*time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
			return
		}
	},
}

var entryTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags",
	Run: func(cmd *cobra.Command, args []string) {
		tags, err := entryManager.GetAllTags()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting tags: %v\n", err)
			return
		}

		if len(tags) == 0 {
			fmt.Println("No tags found")
			return
		}

		fmt.Println("Tags:")
		for _, tag := range tags {
			fmt.Printf("  - %s\n", tag)
		}
	},
}

func printEntries(entries []*storage.Entry) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tUSERNAME\tURL\tTAGS\tUPDATED")
	for _, e := range entries {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			e.ID, e.Name, e.Username, e.URL, e.Tags, e.UpdatedAt.Format("2006-01-02"))
	}
	w.Flush()
}

func init() {
	entryAddCmd.Flags().StringP("name", "n", "", "Entry name")
	entryAddCmd.Flags().StringP("username", "u", "", "Username")
	entryAddCmd.Flags().StringP("password", "p", "", "Password")
	entryAddCmd.Flags().StringP("url", "w", "", "URL")
	entryAddCmd.Flags().StringP("tags", "t", "", "Tags (comma-separated)")
	entryAddCmd.Flags().Bool("generate", false, "Generate a password")
	entryAddCmd.Flags().Int("length", 16, "Length of generated password")

	entrySearchCmd.Flags().StringP("name", "n", "", "Search by name")
	entrySearchCmd.Flags().StringP("tag", "t", "", "Search by tag")

	entryGetCmd.Flags().Bool("show-password", false, "Show the password in plain text")
	entryGetCmd.Flags().BoolP("copy", "c", false, "Copy password to clipboard")

	entryUpdateCmd.Flags().StringP("name", "n", "", "New entry name")
	entryUpdateCmd.Flags().StringP("username", "u", "", "New username")
	entryUpdateCmd.Flags().StringP("password", "p", "", "New password")
	entryUpdateCmd.Flags().StringP("url", "w", "", "New URL")
	entryUpdateCmd.Flags().StringP("tags", "t", "", "New tags (comma-separated)")

	entryCmd.PersistentFlags().StringP("vault", "v", "", "Vault name to use (defaults to last active vault)")

	entryCmd.AddCommand(entryAddCmd)
	entryCmd.AddCommand(entryListCmd)
	entryCmd.AddCommand(entrySearchCmd)
	entryCmd.AddCommand(entryGetCmd)
	entryCmd.AddCommand(entryUpdateCmd)
	entryCmd.AddCommand(entryDeleteCmd)
	entryCmd.AddCommand(entryCopyCmd)
	entryCmd.AddCommand(entryTagsCmd)
	rootCmd.AddCommand(entryCmd)
}
