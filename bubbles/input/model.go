package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// Model displays the content with terminal style.
type Model struct {
	Input textinput.Model

	BlurPrompt   string                 // BlurPrompt controls the prompt to display in blur mode.
	FocusPrompt  string                 // FocusPrompt controls the prompt to display in focus mode.
	Validate     func(str string) error // Validate validsts the value to dislay input with style.
	TextStyle    lipgloss.Style         //TextSyle applys in blur mode.
	ValidStyle   lipgloss.Style         // ValidStyle applys when the input is valid in focus mode.
	InvalidStyle lipgloss.Style         //InvalidStyle applys when the input in invalid in focus mode.

	id uuid.UUID
}

// New creates a new input model.
func New(opts ...func(Model) Model) Model {
	model := Model{
		Input: textinput.New(),
		id:    uuid.New(),
	}
	for _, opt := range opts {
		model = opt(model)
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Input.Focused() {
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

	return m.Input.View()
}

func (m Model) ID() uuid.UUID {
	return m.id
}

func (m Model) Value() string {
	return m.Input.Value()
}

func (m *Model) SetValue(str string) bool {
	if m.Validate != nil {
		if err := m.Validate(str); err != nil {
			return false
		}
	}
	m.Input.SetValue(str)
	return true
}

func (m Model) Focused() bool {
	return m.Input.Focused()
}

func (m *Model) Focus() tea.Cmd {
	return m.Input.Focus()
}

func (m *Model) Blur() tea.Cmd {
	if m.Validate != nil && m.Validate(m.Input.Value()) != nil {
		return nil
	}
	return func() tea.Msg {
		m.Input.Blur()
		return nil
	}
}

func (m *Model) SetWidth(width int) {
	m.Input.Width = width
}

func (m Model) Height() int {
	return 1
}

func (m Model) Width() int {
	return m.Input.Width
}
