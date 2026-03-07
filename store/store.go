package store

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	MaxHistory  = 100
	MaxSnippets = 100
)

// Entry は履歴・スニペット1件を表す
type Entry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Store は履歴ファイルの読み書きを担う
type Store struct {
	path      string
	maxSize   int
	autoPrune bool // trueなら上限超えで古いものを自動削除、falseならエラーを返す
}

// NewHistory はデフォルトパス (~/.icb_history) の履歴Storeを返す
func NewHistory() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{
		path:      filepath.Join(home, ".icb_history"),
		maxSize:   MaxHistory,
		autoPrune: true,
	}, nil
}

// NewSnippets はデフォルトパス (~/.icb_snippets) のスニペットStoreを返す
func NewSnippets() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{
		path:      filepath.Join(home, ".icb_snippets"),
		maxSize:   MaxSnippets,
		autoPrune: false,
	}, nil
}

// NewWithPath は任意のパスの履歴Storeを返す（テスト用）
func NewWithPath(path string) *Store {
	return &Store{path: path, maxSize: MaxHistory, autoPrune: true}
}

// NewSnippetWithPath は任意のパスのスニペットStoreを返す（テスト用）
func NewSnippetWithPath(path string) *Store {
	return &Store{path: path, maxSize: MaxSnippets, autoPrune: false}
}

// Load は履歴ファイルを読み込んで返す
func (s *Store) Load() ([]Entry, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, scanner.Err()
}

// Add はコンテンツを追記し、追加したエントリを返す
// autoPrune=true なら上限超えで古いものを削除、false なら上限超えでエラーを返す
func (s *Store) Add(content string) (Entry, error) {
	entries, err := s.Load()
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{
		ID:        newID(),
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	entries = append(entries, entry)

	if len(entries) > s.maxSize {
		if s.autoPrune {
			entries = entries[len(entries)-s.maxSize:]
		} else {
			return Entry{}, fmt.Errorf("snippet store is full (%d entries), delete some before adding", s.maxSize)
		}
	}

	return entry, s.save(entries)
}

// Delete はIDに一致するエントリを削除する。存在しない場合はno-op。
func (s *Store) Delete(id string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}

	n := 0
	for _, e := range entries {
		if e.ID != id {
			entries[n] = e
			n++
		}
	}
	if n == len(entries) {
		return nil // not found, no-op
	}
	return s.save(entries[:n])
}

func (s *Store) save(entries []Entry) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
