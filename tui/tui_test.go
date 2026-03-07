package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyYank/icb/store"
)

// WindowSizeMsgでmodelのwidthとheightが更新される
func TestModel_WindowSizeUpdate(t *testing.T) {
	m := newModel(nil, nil, nil, nil)
	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	got := result.(model)
	if got.width != 120 {
		t.Errorf("want width 120, got %d", got.width)
	}
	if got.height != 40 {
		t.Errorf("want height 40, got %d", got.height)
	}
}

// 複数行のアイテムが選択されていてwidthが十分なとき、Viewにスプリットセパレータが含まれる
func TestView_SplitViewForMultilineItem(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2\nline3"},
	}
	m := newModel(entries, nil, nil, nil)
	m.width = 120
	m.height = 40

	view := m.View()
	if !strings.Contains(view, "│") {
		t.Errorf("want │ separator in split view, got:\n%s", view)
	}
	// プレビューペインにフル内容が含まれる
	if !strings.Contains(view, "line2") {
		t.Errorf("want preview to contain 'line2', got:\n%s", view)
	}
	if !strings.Contains(view, "line3") {
		t.Errorf("want preview to contain 'line3', got:\n%s", view)
	}
}

// 1行のアイテムのみのときはスプリットビューを表示しない
func TestView_NoSplitViewForSinglelineItems(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "single line only"},
	}
	m := newModel(entries, nil, nil, nil)
	m.width = 120
	m.height = 40

	view := m.View()
	if strings.Contains(view, "│") {
		t.Errorf("want no │ separator for single-line item, got:\n%s", view)
	}
}

// widthが80未満のときはスプリットビューを表示しない
func TestView_NoSplitViewWhenNarrow(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2"},
	}
	m := newModel(entries, nil, nil, nil)
	m.width = 40
	m.height = 40

	view := m.View()
	if strings.Contains(view, "│") {
		t.Errorf("want no │ separator when narrow, got:\n%s", view)
	}
}
