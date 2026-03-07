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

// pキーでpreviewModeがtrueになる
func TestModel_PKeyEntersPreviewMode(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2\nline3"},
	}
	m := newModel(entries, nil, nil, nil)

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = result.(model)
	if !m.previewMode {
		t.Error("want previewMode=true after p")
	}
}

// previewMode中に何かキーを押すとリストに戻る
func TestModel_AnyKeyExitsPreviewMode(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2"},
	}
	m := newModel(entries, nil, nil, nil)
	m.previewMode = true

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = result.(model)
	if m.previewMode {
		t.Error("want previewMode=false after any key in preview mode")
	}
}

// previewMode中のViewは選択アイテムのフル内容を含む
func TestView_PreviewModeShowsFullContent(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2\nline3"},
	}
	m := newModel(entries, nil, nil, nil)
	m.previewMode = true

	view := m.View()
	if !strings.Contains(view, "line1") {
		t.Errorf("want 'line1' in preview, got:\n%s", view)
	}
	if !strings.Contains(view, "line2") {
		t.Errorf("want 'line2' in preview, got:\n%s", view)
	}
	if !strings.Contains(view, "line3") {
		t.Errorf("want 'line3' in preview, got:\n%s", view)
	}
}

// previewMode中はCtrl+Cでquitする
func TestModel_CtrlCQuitsFromPreviewMode(t *testing.T) {
	entries := []store.Entry{
		{ID: "1", Content: "line1\nline2"},
	}
	m := newModel(entries, nil, nil, nil)
	m.previewMode = true

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("want tea.Quit on ctrl+c in preview mode")
	}
}
