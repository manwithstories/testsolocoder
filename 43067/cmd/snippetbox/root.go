package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"snippetbox/internal/config"
	"snippetbox/internal/models"
	"snippetbox/internal/search"
	"snippetbox/internal/storage"
	"snippetbox/internal/template"
)

var rootCmd = &cobra.Command{
	Use:   "snippetbox",
	Short: "A command-line code snippet management tool",
	Long:  `SnippetBox is a powerful command-line tool for managing and reusing code snippets with support for multiple vaults, tags, encryption, and template variables.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(vaultCmd)
	rootCmd.AddCommand(snippetCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(tagCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize SnippetBox with a default vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if cfg.DefaultVault != "default" {
			_, err = storage.LoadVault(cfg.DefaultVault)
			if err == nil {
				fmt.Println("SnippetBox is already initialized.")
				return nil
			}
		}

		vaults, err := storage.ListVaults()
		if err != nil {
			return err
		}

		if len(vaults) > 0 {
			cfg.DefaultVault = vaults[0].ID
			if err := config.Save(cfg); err != nil {
				return err
			}
			fmt.Println("SnippetBox is already initialized.")
			fmt.Printf("Default vault set to: %s\n", cfg.DefaultVault)
			return nil
		}

		vault, err := storage.CreateVault("default")
		if err != nil {
			return err
		}

		cfg.DefaultVault = vault.ID
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Println("SnippetBox initialized successfully!")
		fmt.Printf("Default vault created with ID: %s\n", cfg.DefaultVault)
		return nil
	},
}

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage snippet vaults",
}

var vaultCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		vault, err := storage.CreateVault(name)
		if err != nil {
			return err
		}
		fmt.Printf("Vault '%s' created with ID: %s\n", vault.Name, vault.ID)
		return nil
	},
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaults, err := storage.ListVaults()
		if err != nil {
			return err
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found. Create one with 'snippetbox vault create <name>'.")
			return nil
		}

		cfg, _ := config.Load()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSNIPPETS\tDEFAULT")
		for _, v := range vaults {
			isDefault := ""
			if cfg != nil && cfg.DefaultVault == v.ID {
				isDefault = "*"
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", v.ID, v.Name, len(v.Snippets), isDefault)
		}
		w.Flush()
		return nil
	},
}

var vaultDeleteCmd = &cobra.Command{
	Use:   "delete [vault-id]",
	Short: "Delete a vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID := args[0]

		vault, err := storage.LoadVault(vaultID)
		if err != nil {
			return err
		}

		if len(vault.Snippets) > 0 {
			fmt.Printf("Warning: Vault '%s' contains %d snippet(s).\n", vault.Name, len(vault.Snippets))
			fmt.Print("Are you sure you want to delete it? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				fmt.Println("Operation cancelled.")
				return nil
			}
		}

		if err := storage.DeleteVault(vaultID); err != nil {
			return err
		}

		fmt.Printf("Vault '%s' deleted successfully.\n", vault.Name)
		return nil
	},
}

var vaultUseCmd = &cobra.Command{
	Use:   "use [vault-id]",
	Short: "Set a vault as the default",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID := args[0]

		_, err := storage.LoadVault(vaultID)
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		cfg.DefaultVault = vaultID
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("Default vault set to: %s\n", vaultID)
		return nil
	},
}

var snippetCmd = &cobra.Command{
	Use:   "snippet",
	Short: "Manage code snippets",
}

var (
	snippetTitle       string
	snippetContent     string
	snippetLanguage    string
	snippetTags        []string
	snippetDescription string
	snippetEncrypted   bool
	snippetVault       string
)

var snippetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snippet",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if snippetTitle == "" {
			return errors.New("title is required (use --title)")
		}

		if snippetContent == "" {
			editor := cfg.DefaultEditor
			if editor == "" {
				editor = "vim"
			}
			content, err := editWithEditor(editor, "")
			if err != nil {
				return fmt.Errorf("failed to edit content: %w", err)
			}
			snippetContent = content
		}

		snippet := &models.Snippet{
			Title:       snippetTitle,
			Content:     snippetContent,
			Language:    snippetLanguage,
			Tags:        snippetTags,
			Description: snippetDescription,
			Encrypted:   snippetEncrypted,
		}

		if err := storage.AddSnippet(vaultID, snippet); err != nil {
			return err
		}

		fmt.Printf("Snippet '%s' created with ID: %s\n", snippet.Title, snippet.ID)
		return nil
	},
}

var snippetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		snippets, err := storage.ListSnippets(vaultID)
		if err != nil {
			return err
		}

		if len(snippets) == 0 {
			fmt.Println("No snippets found. Create one with 'snippetbox snippet create'.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tLANGUAGE\tTAGS\tENCRYPTED")
		for _, s := range snippets {
			tags := strings.Join(s.Tags, ", ")
			encrypted := "No"
			if s.Encrypted {
				encrypted = "Yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", s.ID, s.Title, s.Language, tags, encrypted)
		}
		w.Flush()
		return nil
	},
}

var snippetShowCmd = &cobra.Command{
	Use:   "show [snippet-id]",
	Short: "Show a snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		snippetID := args[0]
		snippet, err := storage.GetSnippet(vaultID, snippetID, true)
		if err != nil {
			return err
		}

		fmt.Printf("Title: %s\n", snippet.Title)
		fmt.Printf("ID: %s\n", snippet.ID)
		fmt.Printf("Language: %s\n", snippet.Language)
		fmt.Printf("Tags: %s\n", strings.Join(snippet.Tags, ", "))
		if snippet.Description != "" {
			fmt.Printf("Description: %s\n", snippet.Description)
		}
		fmt.Printf("Created: %s\n", snippet.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", snippet.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("\n--- Code ---")
		fmt.Println(snippet.Content)
		return nil
	},
}

var snippetEditCmd = &cobra.Command{
	Use:   "edit [snippet-id]",
	Short: "Edit a snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		snippetID := args[0]
		snippet, err := storage.GetSnippet(vaultID, snippetID, true)
		if err != nil {
			return err
		}

		if snippetTitle != "" {
			snippet.Title = snippetTitle
		}
		if snippetLanguage != "" {
			snippet.Language = snippetLanguage
		}
		if snippetDescription != "" {
			snippet.Description = snippetDescription
		}
		if len(snippetTags) > 0 {
			snippet.Tags = snippetTags
		}

		if cmd.Flags().Changed("content") || snippetContent != "" {
			snippet.Content = snippetContent
		} else {
			editor := cfg.DefaultEditor
			if editor == "" {
				editor = "vim"
			}
			content, err := editWithEditor(editor, snippet.Content)
			if err != nil {
				return fmt.Errorf("failed to edit content: %w", err)
			}
			snippet.Content = content
		}

		if cmd.Flags().Changed("encrypt") {
			snippet.Encrypted = snippetEncrypted
		}

		if err := storage.UpdateSnippet(vaultID, snippet); err != nil {
			return err
		}

		fmt.Printf("Snippet '%s' updated successfully.\n", snippet.Title)
		return nil
	},
}

var snippetDeleteCmd = &cobra.Command{
	Use:   "delete [snippet-id]",
	Short: "Delete a snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		snippetID := args[0]
		snippet, err := storage.GetSnippet(vaultID, snippetID, false)
		if err != nil {
			return err
		}

		fmt.Printf("Are you sure you want to delete snippet '%s'? (y/N): ", snippet.Title)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			fmt.Println("Operation cancelled.")
			return nil
		}

		if err := storage.DeleteSnippet(vaultID, snippetID); err != nil {
			return err
		}

		fmt.Printf("Snippet '%s' deleted successfully.\n", snippet.Title)
		return nil
	},
}

var snippetUseCmd = &cobra.Command{
	Use:   "use [snippet-id]",
	Short: "Use a snippet (process template variables and copy)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		snippetID := args[0]
		snippet, err := storage.GetSnippet(vaultID, snippetID, true)
		if err != nil {
			return err
		}

		content, err := template.ProcessTemplate(snippet.Content)
		if err != nil {
			return err
		}

		fmt.Println("\n--- Generated Code ---")
		fmt.Println(content)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search snippets",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}

		tags, _ := cmd.Flags().GetStringSlice("tag")
		language, _ := cmd.Flags().GetString("lang")
		fields, _ := cmd.Flags().GetStringSlice("field")

		snippets, err := storage.GetAllSnippets(vaultID, true)
		if err != nil {
			return err
		}

		query := models.SearchQuery{
			Keyword:  keyword,
			Tags:     tags,
			Language: language,
			Fields:   fields,
		}

		results := search.SearchSnippets(snippets, query)

		if len(results) == 0 {
			fmt.Println("No matching snippets found.")
			return nil
		}

		fmt.Printf("Found %d matching snippet(s):\n\n", len(results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tLANGUAGE\tTAGS")
		for _, s := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Title, s.Language, strings.Join(s.Tags, ", "))
		}
		w.Flush()
		return nil
	},
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage snippet tags",
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		snippets, err := storage.ListSnippets(vaultID)
		if err != nil {
			return err
		}

		tagCounts := make(map[string]int)
		for _, s := range snippets {
			for _, tag := range s.Tags {
				tagCounts[tag]++
			}
		}

		if len(tagCounts) == 0 {
			fmt.Println("No tags found.")
			return nil
		}

		fmt.Println("Tags:")
		for tag, count := range tagCounts {
			fmt.Printf("  %s (%d)\n", tag, count)
		}
		return nil
	},
}

var tagFilterCmd = &cobra.Command{
	Use:   "filter [tags...]",
	Short: "Filter snippets by tags",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultID, err := getVaultID()
		if err != nil {
			return err
		}

		snippets, err := storage.GetAllSnippets(vaultID, true)
		if err != nil {
			return err
		}

		query := models.SearchQuery{
			Tags: args,
		}

		results := search.SearchSnippets(snippets, query)

		if len(results) == 0 {
			fmt.Println("No snippets found with the specified tags.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tLANGUAGE\tTAGS")
		for _, s := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Title, s.Language, strings.Join(s.Tags, ", "))
		}
		w.Flush()
		return nil
	},
}

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import vaults from a JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		merge, _ := cmd.Flags().GetBool("merge")

		importedIDs, err := storage.ImportVaults(inputPath, merge)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully imported %d vault(s):\n", len(importedIDs))
		for _, id := range importedIDs {
			fmt.Printf("  - %s\n", id)
		}
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export vaults to a JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath := args[0]
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if err := storage.ExportAllVaults(outputPath); err != nil {
				return err
			}
			fmt.Printf("All vaults exported to: %s\n", outputPath)
		} else {
			vaultID, err := getVaultID()
			if err != nil {
				return err
			}

			if err := storage.ExportVault(vaultID, outputPath); err != nil {
				return err
			}
			fmt.Printf("Vault %s exported to: %s\n", vaultID, outputPath)
		}
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		switch key {
		case "default_vault":
			cfg.DefaultVault = value
		case "default_editor":
			cfg.DefaultEditor = value
		case "highlight_theme":
			cfg.HighlightTheme = value
		case "encryption_key":
			cfg.EncryptionKey = value
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

func getVaultID() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	vaultID := snippetVault
	if vaultID == "" {
		vaultID = cfg.DefaultVault
	}

	if vaultID == "" {
		return "", errors.New("no default vault set. Please create a vault first with 'snippetbox vault create <name>'")
	}

	return vaultID, nil
}

func editWithEditor(editor, content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "snippet-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return string(data), nil
}

func init() {
	snippetCreateCmd.Flags().StringVar(&snippetTitle, "title", "", "Snippet title")
	snippetCreateCmd.Flags().StringVar(&snippetContent, "content", "", "Snippet content")
	snippetCreateCmd.Flags().StringVar(&snippetLanguage, "lang", "", "Programming language")
	snippetCreateCmd.Flags().StringSliceVar(&snippetTags, "tag", []string{}, "Tags (can be specified multiple times)")
	snippetCreateCmd.Flags().StringVar(&snippetDescription, "desc", "", "Description")
	snippetCreateCmd.Flags().BoolVar(&snippetEncrypted, "encrypt", false, "Encrypt the snippet")
	snippetCreateCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	snippetListCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	snippetShowCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	snippetEditCmd.Flags().StringVar(&snippetTitle, "title", "", "New title")
	snippetEditCmd.Flags().StringVar(&snippetContent, "content", "", "New content")
	snippetEditCmd.Flags().StringVar(&snippetLanguage, "lang", "", "New language")
	snippetEditCmd.Flags().StringSliceVar(&snippetTags, "tag", []string{}, "New tags")
	snippetEditCmd.Flags().StringVar(&snippetDescription, "desc", "", "New description")
	snippetEditCmd.Flags().BoolVar(&snippetEncrypted, "encrypt", false, "Enable encryption")
	snippetEditCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	snippetDeleteCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	snippetUseCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	snippetCmd.AddCommand(snippetCreateCmd)
	snippetCmd.AddCommand(snippetListCmd)
	snippetCmd.AddCommand(snippetShowCmd)
	snippetCmd.AddCommand(snippetEditCmd)
	snippetCmd.AddCommand(snippetDeleteCmd)
	snippetCmd.AddCommand(snippetUseCmd)

	vaultCmd.AddCommand(vaultCreateCmd)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultDeleteCmd)
	vaultCmd.AddCommand(vaultUseCmd)

	searchCmd.Flags().StringSlice("tag", []string{}, "Filter by tags")
	searchCmd.Flags().String("lang", "", "Filter by language")
	searchCmd.Flags().StringSlice("field", []string{}, "Search fields (title, content, tags, language, description)")
	searchCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagFilterCmd)
	tagFilterCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to use")

	importCmd.Flags().Bool("merge", false, "Merge with existing vaults")

	exportCmd.Flags().Bool("all", false, "Export all vaults")
	exportCmd.Flags().StringVar(&snippetVault, "vault", "", "Vault ID to export")

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}
