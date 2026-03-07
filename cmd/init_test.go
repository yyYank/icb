package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// zsh引数を渡すとzsh用スクリプトが出力される
func TestRunInit_Zsh(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"init", "zsh"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "zle -N _icb_insert") {
		t.Errorf("want zsh widget registration, got:\n%s", got)
	}
	if !strings.Contains(got, "bindkey '^Xi'") {
		t.Errorf("want Ctrl+X I binding for zsh, got:\n%s", got)
	}
}

// bash引数を渡すとbash用スクリプトが出力される
func TestRunInit_Bash(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"init", "bash"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `bind -x`) {
		t.Errorf("want bash bind -x, got:\n%s", got)
	}
	if !strings.Contains(got, `\C-xi`) {
		t.Errorf("want Ctrl+X I binding for bash, got:\n%s", got)
	}
}

// detectShell: 引数があればそれを使う
func TestDetectShell_FromArgs(t *testing.T) {
	got := detectShell([]string{"zsh"})
	if got != "zsh" {
		t.Errorf("want 'zsh', got %q", got)
	}
}

// detectShell: 引数がなければSHELL環境変数から取る
func TestDetectShell_FromEnv(t *testing.T) {
	t.Setenv("SHELL", "/bin/zsh")
	got := detectShell(nil)
	if got != "zsh" {
		t.Errorf("want 'zsh', got %q", got)
	}
}
