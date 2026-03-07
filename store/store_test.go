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

func newTempSnippetStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	return store.NewSnippetWithPath(filepath.Join(dir, "test_snippets"))
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

// 履歴の上限(100件)を超えたとき古いエントリが削除される
func TestAdd_HistoryPrunesAt100(t *testing.T) {
	s := newTempStore(t)

	for i := 0; i < 101; i++ {
		if err := s.Add("entry"); err != nil {
			t.Fatalf("Add failed at %d: %v", i, err)
		}
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 100 {
		t.Errorf("want 100 entries, got %d", len(entries))
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

// DeleteはIDに一致するエントリを削除する
func TestDelete_RemovesEntry(t *testing.T) {
	s := newTempStore(t)

	if err := s.Add("keep"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := s.Add("delete me"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries, _ := s.Load()
	deleteID := entries[1].ID

	if err := s.Delete(deleteID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	after, err := s.Load()
	if err != nil {
		t.Fatalf("Load after delete failed: %v", err)
	}
	if len(after) != 1 {
		t.Fatalf("want 1 entry after delete, got %d", len(after))
	}
	if after[0].Content != "keep" {
		t.Errorf("want 'keep', got %q", after[0].Content)
	}
}

// 存在しないIDを削除してもエラーにならない
func TestDelete_NonExistentID(t *testing.T) {
	s := newTempStore(t)
	if err := s.Add("entry"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := s.Delete("no-such-id"); err != nil {
		t.Errorf("want no error for non-existent ID, got %v", err)
	}

	entries, _ := s.Load()
	if len(entries) != 1 {
		t.Errorf("want 1 entry unchanged, got %d", len(entries))
	}
}

// スニペットストアは上限(100件)に達すると自動削除せずエラーを返す
func TestSnippetAdd_ErrorWhenFull(t *testing.T) {
	s := newTempSnippetStore(t)

	for i := 0; i < 100; i++ {
		if err := s.Add("snippet"); err != nil {
			t.Fatalf("Add failed at %d: %v", i, err)
		}
	}

	// 101件目はエラーになるべき
	err := s.Add("one more")
	if err == nil {
		t.Error("want error when snippet store is full, got nil")
	}
}

// スニペットストアは上限未満なら追加できる
func TestSnippetAdd_WithinLimit(t *testing.T) {
	s := newTempSnippetStore(t)

	if err := s.Add("my snippet"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
}
