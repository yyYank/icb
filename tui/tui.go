package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyYank/icb/store"
)

const maxVisible = 10

type model struct {
	all      []store.Entry // 全エントリ（新しい順）
	filtered []store.Entry // 検索後エントリ
	cursor   int
	query    string
	selected string
}

func newModel(entries []store.Entry) model {
	// 新しい順に並べる
	reversed := make([]store.Entry, len(entries))
	for i, e := range entries {
		reversed[len(entries)-1-i] = e
	}
	return model{
		all:      reversed,
		filtered: reversed,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor].Content
			}
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "backspace":
			if len(m.query) > 0 {
				// rune単位で削除
				runes := []rune(m.query)
				m.query = string(runes[:len(runes)-1])
				m.cursor = 0
				m.filtered = filterEntries(m.all, m.query)
			}
		default:
			if len(msg.Runes) > 0 {
				m.query += string(msg.Runes)
				m.cursor = 0
				m.filtered = filterEntries(m.all, m.query)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	fmt.Fprintf(&b, "> %s\n", m.query)
	b.WriteString(strings.Repeat("─", 40) + "\n")

	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		line := m.filtered[i].Content
		// 複数行コンテンツは最初の行だけ表示
		if idx := strings.Index(line, "\n"); idx >= 0 {
			line = line[:idx] + " ..."
		}
		if i == m.cursor {
			fmt.Fprintf(&b, "▶ %s\n", line)
		} else {
			fmt.Fprintf(&b, "  %s\n", line)
		}
	}

	b.WriteString(strings.Repeat("─", 40) + "\n")
	fmt.Fprintf(&b, "%d/%d  ↑↓で移動  Enterで選択  Ctrl+Cでキャンセル\n", len(m.filtered), len(m.all))

	return b.String()
}

// Run はTUIを起動し、選択されたコンテンツを返す
func Run(entries []store.Entry) (string, error) {
	if len(entries) == 0 {
		return "", nil
	}

	m := newModel(entries)
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	result, err := p.Run()
	if err != nil {
		return "", err
	}

	return result.(model).selected, nil
}

func filterEntries(entries []store.Entry, query string) []store.Entry {
	if query == "" {
		return entries
	}
	q := strings.ToLower(query)
	var result []store.Entry
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Content), q) {
			result = append(result, e)
		}
	}
	return result
}
