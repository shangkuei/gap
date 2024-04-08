package text

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// Model displays the content with terminal style.
// Model is thread-safe after running Init().
type Model struct {
	Style lipgloss.Style

	content any
	width   int

	id uuid.UUID
}

func New(opts ...func(Model) Model) Model {
	model := Model{
		Style: lipgloss.NewStyle(),
	}
	for _, opt := range opts {
		model = opt(model)
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	style := m.Style.Inline(true)
	if m.width > 0 {
		style = style.Width(m.width)
	}
	return style.Render(fmt.Sprint(m.content))
}

func (m Model) ID() uuid.UUID {
	return m.id
}

func (m *Model) SetContent(content any) {
	m.content = content
}

func (m Model) Value() any {
	return m.content
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
