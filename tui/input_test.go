package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Enterキーで改行が追加される
func TestInputModel_EnterAddsNewline(t *testing.T) {
	m := inputModel{}
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyEnter})
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("world")})

	if m.value != "hello\nworld" {
		t.Errorf("want 'hello\\nworld', got %q", m.value)
	}
}

// Ctrl+Sで確定してquitする
func TestInputModel_CtrlSConfirms(t *testing.T) {
	m := inputModel{}
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("line1")})

	var cmd tea.Cmd
	m, cmd = update(m, tea.KeyMsg{Type: tea.KeyCtrlS})

	if cmd == nil {
		t.Error("want tea.Quit cmd, got nil")
	}
	if m.canceled {
		t.Error("want canceled=false on ctrl+s")
	}
}

// Ctrl+Cでキャンセルする
func TestInputModel_CtrlCCancels(t *testing.T) {
	m := inputModel{}
	var cmd tea.Cmd
	m, cmd = update(m, tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Error("want tea.Quit cmd, got nil")
	}
	if !m.canceled {
		t.Error("want canceled=true on ctrl+c")
	}
}

// 複数行のテキストが正しく保持される
func TestInputModel_MultilineValue(t *testing.T) {
	m := inputModel{}
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("line1")})
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyEnter})
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("line2")})
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyEnter})
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("line3")})

	if m.value != "line1\nline2\nline3" {
		t.Errorf("want 'line1\\nline2\\nline3', got %q", m.value)
	}
}

// ヘルパー: Update の戻り値を inputModel にキャストする
func update(m inputModel, msg tea.Msg) (inputModel, tea.Cmd) {
	result, cmd := m.Update(msg)
	return result.(inputModel), cmd
}
