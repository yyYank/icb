package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yyYank/icb/store"
	"github.com/yyYank/icb/tui"
)

var inputFn = tui.RunInput

var historyStoreFn = func() (*store.Store, error) {
	return store.NewHistory()
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new entry to history interactively",
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
	histStore, err := historyStoreFn()
	if err != nil {
		return err
	}
	_, err = histStore.Add(text)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "added to history")
	return nil
}
