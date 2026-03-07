package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yyYank/icb/store"
)

func newTempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	return store.NewWithPath(filepath.Join(dir, "test_history"))
}

// 空の履歴ファイルが存在しないとき、Loadはnilを返す
func TestLoad_NoFile(t *testing.T) {
	s := newTempStore(t)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("want nil error, got %v", err)
	}
	if entries != nil {
		t.Fatalf("want nil entries, got %v", entries)
	}
}

// AddしたエントリがあとでLoadできる
func TestAdd_AndLoad(t *testing.T) {
	s := newTempStore(t)

	if err := s.Add("hello world"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Content != "hello world" {
		t.Errorf("want 'hello world', got %q", entries[0].Content)
	}
	if entries[0].ID == "" {
		t.Error("want non-empty ID")
	}
	if entries[0].CreatedAt.IsZero() {
		t.Error("want non-zero CreatedAt")
	}
}

// 複数エントリを順番通りに保持する
func TestAdd_MultipleEntries(t *testing.T) {
	s := newTempStore(t)

	texts := []string{"first", "second", "third"}
	for _, text := range texts {
		if err := s.Add(text); err != nil {
			t.Fatalf("Add(%q) failed: %v", text, err)
		}
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
	for i, want := range texts {
		if entries[i].Content != want {
			t.Errorf("entries[%d]: want %q, got %q", i, want, entries[i].Content)
		}
	}
}

// 上限(1000件)を超えたとき古いエントリが削除される
func TestAdd_PrunesOldEntries(t *testing.T) {
	s := newTempStore(t)

	for i := 0; i < 1001; i++ {
		if err := s.Add("entry"); err != nil {
			t.Fatalf("Add failed at %d: %v", i, err)
		}
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1000 {
		t.Errorf("want 1000 entries, got %d", len(entries))
	}
}

// 改行を含むコンテンツが正しく保存・復元される
func TestAdd_MultilineContent(t *testing.T) {
	s := newTempStore(t)

	content := "line1\nline2\nline3"
	if err := s.Add(content); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Content != content {
		t.Errorf("want %q, got %q", content, entries[0].Content)
	}
}

// 不正な行があっても他のエントリは読み込める
func TestLoad_SkipsInvalidLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history")

	// 正常なJSON行と不正な行を混在させて書き込む
	data := `{"id":"1","content":"valid","created_at":"2026-01-01T00:00:00Z"}
not-json-at-all
{"id":"2","content":"also valid","created_at":"2026-01-01T00:00:00Z"}
`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	s := store.NewWithPath(path)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
}
