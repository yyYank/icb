package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/yyYank/icb/store"
)

// icb addでテキストを入力するとスニペットとして保存され "saved as snippet" が出力される
func TestRunAdd_SavesSnippet(t *testing.T) {
	origInputFn := inputFn
	defer func() { inputFn = origInputFn }()
	inputFn = func() (string, error) {
		return "test snippet", nil
	}

	f, err := os.CreateTemp("", "icb_snippets_test_*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	origSnippetStoreFn := snippetStoreFn
	defer func() { snippetStoreFn = origSnippetStoreFn }()
	snippetStoreFn = func() (*store.Store, error) {
		return store.NewSnippetWithPath(f.Name()), nil
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"add"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(buf.String(), "saved as snippet") {
		t.Errorf("want 'saved as snippet', got: %q", buf.String())
	}

	// スニペットが実際に保存されたことを確認
	s := store.NewSnippetWithPath(f.Name())
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Content != "test snippet" {
		t.Errorf("want content 'test snippet', got %q", entries[0].Content)
	}
}

// 入力が空のときは何も保存しない
func TestRunAdd_EmptyInput(t *testing.T) {
	origInputFn := inputFn
	defer func() { inputFn = origInputFn }()
	inputFn = func() (string, error) {
		return "", nil
	}

	f, err := os.CreateTemp("", "icb_snippets_test_*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	origSnippetStoreFn := snippetStoreFn
	defer func() { snippetStoreFn = origSnippetStoreFn }()
	snippetStoreFn = func() (*store.Store, error) {
		return store.NewSnippetWithPath(f.Name()), nil
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"add"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// 空入力時は何も保存されない
	s := store.NewSnippetWithPath(f.Name())
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("want 0 entries, got %d", len(entries))
	}
}
