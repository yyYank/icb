package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	value    string
	canceled bool
}

func (m inputModel) Init() tea.Cmd {
	return nil
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.canceled = true
			return m, tea.Quit
		case tea.KeyCtrlS:
			return m, tea.Quit
		case tea.KeyEnter:
			m.value += "\n"
		case tea.KeyBackspace:
			if len(m.value) > 0 {
				runes := []rune(m.value)
				m.value = string(runes[:len(runes)-1])
			}
		default:
			if len(msg.Runes) > 0 {
				m.value += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m inputModel) View() string {
	return fmt.Sprintf("Enter snippet (Ctrl+S to save, Ctrl+C to cancel):\n%s▌\n", m.value)
}

// RunInput はテキスト入力TUIを起動し、入力されたテキストを返す
func RunInput() (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open /dev/tty: %w", err)
	}
	defer tty.Close()

	m := inputModel{}
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))

	result, err := p.Run()
	if err != nil {
		return "", err
	}

	final := result.(inputModel)
	if final.canceled {
		return "", nil
	}
	return strings.TrimRight(final.value, "\n"), nil
}
