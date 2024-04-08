package input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/muesli/reflow/ansi"
	"github.com/shangkuei/gap/bubbles"
)

// Model displays the content with terminal style.
type Model struct {
	BlurPrompt   string                 // BlurPrompt controls the prompt to display in blur mode.
	FocusPrompt  string                 // FocusPrompt controls the prompt to display in focus mode.
	Validate     func(str string) error // Validate validsts the value to display input with style.
	TextStyle    lipgloss.Style         // TextSyle applies in blur mode.
	ValidStyle   lipgloss.Style         // ValidStyle applies when the input is valid in focus mode.
	InvalidStyle lipgloss.Style         //InvalidStyle applies when the input in invalid in focus mode.

	id     uuid.UUID
	input  textinput.Model
	cached string
	width  int
}

// New creates a new input model.
func New(opts ...func(Model) Model) Model {
	model := Model{
		id:    uuid.New(),
		input: textinput.New(),
	}
	model.input.KeyMap = textinput.KeyMap{
		CharacterForward: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", ""),
		),
		CharacterBackward: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", ""),
		),
		DeleteCharacterBackward: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("BS", ""),
		),
		DeleteCharacterForward: key.NewBinding(
			key.WithKeys("delete"),
			key.WithHelp("DEL", ""),
		),
		LineStart: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("HOME", ""),
		),
		LineEnd: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("END", ""),
		),
		Paste: key.NewBinding(
			key.WithKeys("ctrl+v"),
			key.WithHelp("ctrl+v", "貼上"),
		),
		AcceptSuggestion: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "套用"),
		),
		NextSuggestion: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", ""),
		),
		PrevSuggestion: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", ""),
		),
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
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.input.Focused() {
		m.input.Prompt = m.FocusPrompt
		if m.Validate == nil || m.Validate(m.input.Value()) == nil {
			m.input.TextStyle = m.ValidStyle
			m.input.Cursor.TextStyle = m.ValidStyle
		} else {
			m.input.TextStyle = m.InvalidStyle
			m.input.Cursor.TextStyle = m.InvalidStyle
		}
	} else {
		m.input.Prompt = m.BlurPrompt
		m.input.TextStyle = m.TextStyle
		m.input.Cursor.TextStyle = m.TextStyle
	}

	// the textinput v0.18.0 doesn't count prompt length.
	m.input.Width = m.width - ansi.PrintableRuneWidth(m.input.Prompt) - 1
	return m.input.View()
}

func (m Model) ID() uuid.UUID {
	return m.id
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m *Model) SetValue(str string) bool {
	if m.Validate != nil {
		if err := m.Validate(str); err != nil {
			return false
		}
	}
	m.cached = str
	m.cached = str
	m.input.SetValue(str)
	return true
}

func (m *Model) EnableSuggestion(value bool) {
	m.input.ShowSuggestions = value
}

func (m *Model) SetSuggestions(suggestions []string) {
	m.input.SetSuggestions(suggestions)
}

func (m Model) Focused() bool {
	return m.input.Focused()
}

func (m *Model) Focus() tea.Cmd {
	m.input.Focus()
	return nil
}

func (m *Model) Blur() tea.Cmd {
	if m.Validate != nil && m.Validate(m.input.Value()) != nil {
		if m.cached == "" {
			return nil
		}
		m.input.SetValue(m.cached)
	}
	m.input.Blur()
	return nil
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) Height() int {
	return 1
}

func (m Model) Width() int {
	return m.width
}

func (m Model) KeyBindings() (bindings []key.Binding) {
	bindings = append(bindings,
		m.input.KeyMap.CharacterBackward,
		m.input.KeyMap.CharacterForward,
	)
	if m.input.ShowSuggestions {
		bindings = append(bindings,
			m.input.KeyMap.NextSuggestion,
			m.input.KeyMap.PrevSuggestion,
			m.input.KeyMap.AcceptSuggestion,
		)
	}
	bindings = append(bindings,
		m.input.KeyMap.Paste,
		m.input.KeyMap.DeleteCharacterBackward,
		m.input.KeyMap.DeleteCharacterForward,
		m.input.KeyMap.LineStart,
		m.input.KeyMap.LineEnd,
	)
	return bindings
}

var _ bubbles.FocusModel = (*Model)(nil)
