package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yyYank/icb/store"
	"github.com/yyYank/icb/tui"
)

var rootCmd = &cobra.Command{
	Use:   "icb",
	Short: "Internal/Isolated Clipboard — terminal clipboard history",
	Long: `icb is a standalone clipboard history tool for terminal environments.
Works entirely within your shell — no OS clipboard, no GUI, no dependencies.

  echo "some text" | icb   # store
  icb                      # pick from history → stdout`,
	RunE:          run,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	histStore, err := store.NewHistory()
	if err != nil {
		return err
	}

	// パイプあり → 蓄積モード
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		text := strings.TrimRight(string(content), "\n")
		if text == "" {
			return nil
		}
		return histStore.Add(text)
	}

	// TTY → TUI起動
	snippetStore, err := store.NewSnippets()
	if err != nil {
		return err
	}

	history, err := histStore.Load()
	if err != nil {
		return err
	}

	snippets, err := snippetStore.Load()
	if err != nil {
		return err
	}

	selected, err := tui.Run(history, snippets, histStore, snippetStore)
	if err != nil {
		return err
	}

	if selected != "" {
		fmt.Println(selected)
	}
	return nil
}
