package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/yyYank/icb/store"
)

// icb addでテキストを入力するとヒストリーに保存され "added to history" が出力される
func TestRunAdd_AddsToHistory(t *testing.T) {
	origInputFn := inputFn
	defer func() { inputFn = origInputFn }()
	inputFn = func() (string, error) {
		return "test entry", nil
	}

	f, err := os.CreateTemp("", "icb_history_test_*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	origHistoryStoreFn := historyStoreFn
	defer func() { historyStoreFn = origHistoryStoreFn }()
	historyStoreFn = func() (*store.Store, error) {
		return store.NewWithPath(f.Name()), nil
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"add"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(buf.String(), "added to history") {
		t.Errorf("want 'added to history', got: %q", buf.String())
	}

	// ヒストリーに実際に保存されたことを確認
	s := store.NewWithPath(f.Name())
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Content != "test entry" {
		t.Errorf("want content 'test entry', got %q", entries[0].Content)
	}
}

// 入力が空のときは何も保存しない
func TestRunAdd_EmptyInput(t *testing.T) {
	origInputFn := inputFn
	defer func() { inputFn = origInputFn }()
	inputFn = func() (string, error) {
		return "", nil
	}

	f, err := os.CreateTemp("", "icb_history_test_*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	origHistoryStoreFn := historyStoreFn
	defer func() { historyStoreFn = origHistoryStoreFn }()
	historyStoreFn = func() (*store.Store, error) {
		return store.NewWithPath(f.Name()), nil
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"add"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// 空入力時は何も保存されない
	s := store.NewWithPath(f.Name())
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("want 0 entries, got %d", len(entries))
	}
}
