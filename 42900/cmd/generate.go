package cmd

import (
	"fmt"
	"os"
	"passman/clipboard"
	"passman/generator"
	"time"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a strong password",
	Long:  "Generate a strong password with customizable options.",
	Run: func(cmd *cobra.Command, args []string) {
		length, _ := cmd.Flags().GetInt("length")
		includeLower, _ := cmd.Flags().GetBool("lower")
		includeUpper, _ := cmd.Flags().GetBool("upper")
		includeDigits, _ := cmd.Flags().GetBool("digits")
		includeSpecial, _ := cmd.Flags().GetBool("special")
		excludeAmbiguous, _ := cmd.Flags().GetBool("no-ambiguous")
		phrase, _ := cmd.Flags().GetBool("phrase")
		words, _ := cmd.Flags().GetInt("words")
		copyGen, _ := cmd.Flags().GetBool("copy")

		if phrase {
			password, err := generator.GeneratePhrase(words)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating passphrase: %v\n", err)
				return
			}
			fmt.Println(password)

			if copyGen {
				_ = clipboard.CopyToClipboard(password, 30*time.Second)
			}
			return
		}

		opts := &generator.Options{
			Length:           length,
			IncludeLower:     includeLower,
			IncludeUpper:     includeUpper,
			IncludeDigits:    includeDigits,
			IncludeSpecial:   includeSpecial,
			ExcludeAmbiguous: excludeAmbiguous,
		}

		password, err := generator.Generate(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
			return
		}

		fmt.Println(password)

		if copyGen {
			_ = clipboard.CopyToClipboard(password, 30*time.Second)
		}
	},
}

func init() {
	generateCmd.Flags().IntP("length", "l", 16, "Password length (8-128)")
	generateCmd.Flags().Bool("lower", true, "Include lowercase letters")
	generateCmd.Flags().Bool("upper", true, "Include uppercase letters")
	generateCmd.Flags().Bool("digits", true, "Include digits")
	generateCmd.Flags().Bool("special", true, "Include special characters")
	generateCmd.Flags().Bool("no-ambiguous", true, "Exclude ambiguous characters (0, O, 1, l, I)")
	generateCmd.Flags().Bool("phrase", false, "Generate a passphrase instead of random characters")
	generateCmd.Flags().Int("words", 4, "Number of words in passphrase (3-10)")
	generateCmd.Flags().BoolP("copy", "c", false, "Copy generated password to clipboard")
	rootCmd.AddCommand(generateCmd)
}
