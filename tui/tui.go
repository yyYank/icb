package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyYank/icb/store"
)

const maxVisible = 10

// Source はエントリの種別を表す
type Source int

const (
	SourceHistory Source = iota
	SourceSnippet
)

// Item はTUI表示用エントリ（種別付き）
type Item struct {
	Entry  store.Entry
	Source Source
}

type model struct {
	all          []Item
	filtered     []Item
	cursor       int
	query        string
	selected     string
	historyStore *store.Store
	snippetStore *store.Store
	statusMsg   string
	width       int
	height      int
	showPreview bool
}

func newModel(history, snippets []store.Entry, histStore, snippetStore *store.Store) model {
	var all []Item
	// スニペットを先頭に（新しい順）
	for i := len(snippets) - 1; i >= 0; i-- {
		all = append(all, Item{Entry: snippets[i], Source: SourceSnippet})
	}
	// 履歴を後に（新しい順）
	for i := len(history) - 1; i >= 0; i-- {
		all = append(all, Item{Entry: history[i], Source: SourceHistory})
	}
	return model{
		all:          all,
		filtered:     all,
		historyStore: histStore,
		snippetStore: snippetStore,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		m.statusMsg = "" // キー操作でステータスをクリア

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor].Entry.Content
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
		case "d":
			if len(m.filtered) > 0 {
				item := m.filtered[m.cursor]
				var err error
				if item.Source == SourceSnippet {
					err = m.snippetStore.Delete(item.Entry.ID)
				} else {
					err = m.historyStore.Delete(item.Entry.ID)
				}
				if err != nil {
					m.statusMsg = "error: " + err.Error()
				} else {
					m.all = removeItem(m.all, item.Entry.ID)
					m.filtered = filterItems(m.all, m.query)
					if m.cursor >= len(m.filtered) && m.cursor > 0 {
						m.cursor--
					}
					m.statusMsg = "deleted"
				}
			}
		case "p":
			m.showPreview = !m.showPreview
		case "s":
			if len(m.filtered) > 0 {
				item := m.filtered[m.cursor]
				if item.Source == SourceSnippet {
					m.statusMsg = "already a snippet"
				} else {
					newEntry, err := m.snippetStore.Add(item.Entry.Content)
					if err != nil {
						m.statusMsg = "error: " + err.Error()
					} else {
						// 先頭にスニペット行を追加して即反映
						newItem := Item{Entry: newEntry, Source: SourceSnippet}
						m.all = append([]Item{newItem}, m.all...)
						m.filtered = filterItems(m.all, m.query)
						m.statusMsg = "saved as snippet"
					}
				}
			}
		case "backspace":
			if len(m.query) > 0 {
				runes := []rune(m.query)
				m.query = string(runes[:len(runes)-1])
				m.cursor = 0
				m.filtered = filterItems(m.all, m.query)
			}
		default:
			if len(msg.Runes) > 0 {
				m.query += string(msg.Runes)
				m.cursor = 0
				m.filtered = filterItems(m.all, m.query)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.showPreview && m.width >= 80 {
		return m.splitView()
	}
	return m.singleView()
}

func (m model) currentContent() string {
	if len(m.filtered) == 0 {
		return ""
	}
	return m.filtered[m.cursor].Entry.Content
}

func (m model) singleView() string {
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
		item := m.filtered[i]
		line := item.Entry.Content
		if idx := strings.Index(line, "\n"); idx >= 0 {
			line = line[:idx] + " ..."
		}

		cursor := " "
		if i == m.cursor {
			cursor = "▶"
		}
		star := " "
		if item.Source == SourceSnippet {
			star = "★"
		}
		fmt.Fprintf(&b, "%s%s %s\n", cursor, star, line)
	}

	b.WriteString(strings.Repeat("─", 40) + "\n")

	if m.statusMsg != "" {
		fmt.Fprintf(&b, "%s\n", m.statusMsg)
	} else {
		fmt.Fprintf(&b, "%d/%d  d:delete  s:snippet  p:preview  Enter:select  Ctrl+C:cancel\n",
			len(m.filtered), len(m.all))
	}

	return b.String()
}

func (m model) splitView() string {
	leftWidth := m.width / 2

	leftLines := m.buildLeftLines(leftWidth)
	rightLines := strings.Split(m.currentContent(), "\n")

	var b strings.Builder
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}
	for i := 0; i < maxLines; i++ {
		l := ""
		if i < len(leftLines) {
			l = leftLines[i]
		}
		r := ""
		if i < len(rightLines) {
			r = rightLines[i]
		}
		fmt.Fprintf(&b, "%s│%s\n", padRight(l, leftWidth), r)
	}

	if m.statusMsg != "" {
		fmt.Fprintf(&b, "%s\n", m.statusMsg)
	} else {
		fmt.Fprintf(&b, "%d/%d  d:delete  s:snippet  p:preview  Enter:select  Ctrl+C:cancel\n",
			len(m.filtered), len(m.all))
	}

	return b.String()
}

func (m model) buildLeftLines(width int) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("> %s", m.query))
	lines = append(lines, strings.Repeat("─", width))

	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		item := m.filtered[i]
		line := item.Entry.Content
		if idx := strings.Index(line, "\n"); idx >= 0 {
			line = line[:idx] + " ..."
		}
		cursor := " "
		if i == m.cursor {
			cursor = "▶"
		}
		star := " "
		if item.Source == SourceSnippet {
			star = "★"
		}
		lines = append(lines, fmt.Sprintf("%s%s %s", cursor, star, line))
	}

	lines = append(lines, strings.Repeat("─", width))
	return lines
}

func padRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-len(runes))
}

// Run はTUIを起動し、選択されたコンテンツを返す
func Run(history, snippets []store.Entry, histStore, snippetStore *store.Store) (string, error) {
	if len(history) == 0 && len(snippets) == 0 {
		return "", nil
	}

	// シェルウィジェットから呼ばれると stdin/stdout がリダイレクトされているため、
	// /dev/tty を直接開いて入出力を端末に固定する
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open /dev/tty: %w", err)
	}
	defer tty.Close()

	m := newModel(history, snippets, histStore, snippetStore)
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))

	result, err := p.Run()
	if err != nil {
		return "", err
	}

	return result.(model).selected, nil
}

func filterItems(items []Item, query string) []Item {
	if query == "" {
		return items
	}
	q := strings.ToLower(query)
	var result []Item
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Entry.Content), q) {
			result = append(result, item)
		}
	}
	return result
}

func removeItem(items []Item, id string) []Item {
	result := make([]Item, 0, len(items))
	for _, item := range items {
		if item.Entry.ID != id {
			result = append(result, item)
		}
	}
	return result
}
