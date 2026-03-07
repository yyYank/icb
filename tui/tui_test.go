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

// pキーでshowPreviewがトグルされる
func TestModel_PKeyTogglesPreview(t *testing.T) {
	m := newModel(nil, nil, nil, nil)
	if m.showPreview {
		t.Fatal("want showPreview=false initially")
	}

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = result.(model)
	if !m.showPreview {
		t.Error("want showPreview=true after first p")
	}

	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = result.(model)
	if m.showPreview {
		t.Error("want showPreview=false after second p")
	}
}

// showPreview=trueかつwidth>=80のとき、スプリットセパレータが表示される
func TestView_SplitViewWhenPreviewEnabled(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2\nline3"},
	}
	m := newModel(entries, nil, nil, nil)
	m.width = 120
	m.height = 40
	m.showPreview = true

	view := m.View()
	if !strings.Contains(view, "│") {
		t.Errorf("want │ separator in split view, got:\n%s", view)
	}
	if !strings.Contains(view, "line2") {
		t.Errorf("want preview to contain 'line2', got:\n%s", view)
	}
	if !strings.Contains(view, "line3") {
		t.Errorf("want preview to contain 'line3', got:\n%s", view)
	}
}

// showPreview=falseのときはスプリットビューを表示しない
func TestView_NoSplitViewWhenPreviewDisabled(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2\nline3"},
	}
	m := newModel(entries, nil, nil, nil)
	m.width = 120
	m.height = 40
	// showPreview はデフォルト false

	view := m.View()
	if strings.Contains(view, "│") {
		t.Errorf("want no │ separator when preview disabled, got:\n%s", view)
	}
}

// widthが80未満のときはshowPreview=trueでもスプリットビューを表示しない
func TestView_NoSplitViewWhenNarrow(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2"},
	}
	m := newModel(entries, nil, nil, nil)
	m.width = 40
	m.height = 40
	m.showPreview = true

	view := m.View()
	if strings.Contains(view, "│") {
		t.Errorf("want no │ separator when narrow, got:\n%s", view)
	}
}
