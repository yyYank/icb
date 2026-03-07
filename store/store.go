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

const maxEntries = 1000

// Entry は履歴の1件を表す
type Entry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Store は履歴ファイルの読み書きを担う
type Store struct {
	path string
}

// New はデフォルトパス (~/.icb_history) のStoreを返す
func New() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(home, ".icb_history")}, nil
}

// NewWithPath は任意のパスのStoreを返す（テスト用）
func NewWithPath(path string) *Store {
	return &Store{path: path}
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

// Add はコンテンツを履歴に追記する（上限超えたら古いものを削除）
func (s *Store) Add(content string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}

	entry := Entry{
		ID:        newID(),
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	entries = append(entries, entry)

	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	return s.save(entries)
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
