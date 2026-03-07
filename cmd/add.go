package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yyYank/icb/store"
	"github.com/yyYank/icb/tui"
)

var inputFn = tui.RunInput

var snippetStoreFn = func() (*store.Store, error) {
	return store.NewSnippets()
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new snippet interactively",
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	text, err := inputFn()
	if err != nil || text == "" {
		return err
	}
	snippetStore, err := snippetStoreFn()
	if err != nil {
		return err
	}
	_, err = snippetStore.Add(text)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "saved as snippet")
	return nil
}
