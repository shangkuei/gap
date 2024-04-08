package text

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// Model displays the content with terminal style.
// Model is thread-safe after running Init().
type Model struct {
	Style   lipgloss.Style
	Content any

	id uuid.UUID
	mu sync.RWMutex
}

func (m *Model) Init() tea.Cmd {
	m.id = uuid.New()
	return nil
}

func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Style.Inline(true).Render(fmt.Sprint(m.Content))
}

func (m *Model) ID() uuid.UUID {
	return m.id
}

func (m *Model) SetContent(content any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Content = content
}

func (m *Model) Value() string {
	return fmt.Sprint(m.Content)
}

func (m *Model) String() string {
	return m.View()
}

func (m *Model) Height() int {
	return 1
}
