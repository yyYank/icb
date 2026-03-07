package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const zshScript = `
# icb shell integration (zsh)
_icb_insert() {
  local result
  result=$(icb 2>/dev/tty)
  if [[ -n "$result" ]]; then
    LBUFFER+="$result"
  fi
  zle reset-prompt
}
zle -N _icb_insert
bindkey '^Xi' _icb_insert
`

const bashScript = `
# icb shell integration (bash)
_icb_insert() {
  local result
  result=$(icb 2>/dev/tty)
  if [[ -n "$result" ]]; then
    READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${result}${READLINE_LINE:$READLINE_POINT}"
    READLINE_POINT=$(( READLINE_POINT + ${#result} ))
  fi
}
bind -x '"\C-xi": _icb_insert'
`

var initCmd = &cobra.Command{
	Use:   "init [bash|zsh]",
	Short: "Print shell integration script",
	Long: `Print shell integration script to enable Ctrl+X I keybinding.

Add the following line to your .zshrc or .bashrc:

  eval "$(icb init)"

Then Ctrl+X I will insert the selected clipboard entry at the cursor.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	shell := detectShell(args)
	switch shell {
	case "zsh":
		fmt.Fprint(cmd.OutOrStdout(), zshScript)
	case "bash":
		fmt.Fprint(cmd.OutOrStdout(), bashScript)
	default:
		return fmt.Errorf("unsupported shell: %q (supported: bash, zsh)", shell)
	}
	return nil
}

func detectShell(args []string) string {
	if len(args) > 0 {
		return strings.ToLower(args[0])
	}
	shellPath := os.Getenv("SHELL")
	return filepath.Base(shellPath)
}
