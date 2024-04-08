package input

import (
	"sync"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// Model displays the content with terminal style.
// Model is thread-safe after running Init().
type Model struct {
	Input textinput.Model

	BlurPrompt   string                 // BlurPrompt controls the prompt to display in blur mode.
	FocusPrompt  string                 // FocusPrompt controls the prompt to display in focus mode.
	Validate     func(str string) error // Validate validsts the value to display input with style.
	TextStyle    lipgloss.Style         //TextSyle applies in blur mode.
	ValidStyle   lipgloss.Style         // ValidStyle applies when the input is valid in focus mode.
	InvalidStyle lipgloss.Style         //InvalidStyle applies when the input in invalid in focus mode.

	id uuid.UUID
	mu sync.RWMutex
}

func (m *Model) Init() tea.Cmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.id = uuid.New()
	return nil
}

func (m *Model) Update(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var style lipgloss.Style
	if m.Input.Focused() && m.Validate != nil {
		m.Input.Prompt = m.FocusPrompt
		if m.Validate == nil || m.Validate(m.Input.Value()) == nil {
			m.Input.TextStyle = m.ValidStyle
			m.Input.Cursor.TextStyle = m.ValidStyle
		} else {
			m.Input.TextStyle = m.InvalidStyle
			m.Input.Cursor.TextStyle = m.InvalidStyle
		}
	} else {
		m.Input.Prompt = m.BlurPrompt
		m.Input.TextStyle = m.TextStyle
		m.Input.Cursor.TextStyle = m.TextStyle
	}

	return style.Render(m.Input.View())
}

func (m *Model) ID() uuid.UUID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.id
}

func (m *Model) Value() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.Input.Value()
}

func (m *Model) SetValue(str string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Validate != nil {
		if err := m.Validate(str); err != nil {
			return false
		}
	}
	m.Input.SetValue(str)
	return true
}

func (m *Model) Focused() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.Input.Focused()
}

func (m *Model) Focus() (bool, tea.Cmd) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Input.Cursor.Blink = true
	return true, m.Input.Focus()
}

func (m *Model) Blur() (bool, tea.Cmd) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Validate != nil {
		if err := m.Validate(m.Input.Value()); err == nil {
			m.Input.Cursor.Blink = false
			m.Input.Blur()
			return true, nil
		} else {
			return false, nil
		}
	}
	m.Input.Blur()
	return true, nil
}

func (m *Model) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.View()
}

func (m *Model) SetWidth(width int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Input.Width = width
}

func (m *Model) Height() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return 1
}

func (m *Model) Width() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Input.Width
}
