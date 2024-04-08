package bubbles

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// SequenceMsg is a message that represents a sequence of tea.Cmd.
type SequenceMsg []tea.Cmd

// Sequence returns a tea.Cmd that represents a sequence of tea.Cmd.
func Sequence(cmds ...tea.Cmd) tea.Cmd {
	var validCmds []tea.Cmd
	for _, c := range cmds {
		if c == nil {
			continue
		}
		validCmds = append(validCmds, c)
	}

	switch len(validCmds) {
	case 0:
		return nil
	case 1:
		return validCmds[0]
	default:
		return func() tea.Msg {
			return SequenceMsg(validCmds)
		}
	}

}

// BatchMsg is a message that represents a batch of tea.Cmd.
type BatchMsg []tea.Cmd

// Batch returns a tea.Cmd that represents a batch of tea.Cmd.
func Batch(cmds ...tea.Cmd) tea.Cmd {
	var validCmds []tea.Cmd
	for _, c := range cmds {
		if c == nil {
			continue
		}
		validCmds = append(validCmds, c)
	}

	switch len(validCmds) {
	case 0:
		return nil
	case 1:
		return validCmds[0]
	default:
		return func() tea.Msg {
			return BatchMsg(validCmds)
		}
	}
}

// FrameMsg is a message that force to refresh a frame.
type FrameMsg struct{}

// Filter filters out SequenceMsg and BatchMsg and transforms them to official Sequence and Batch
// tea.Msg. This brings support to nested models.
func Filter(m tea.Model, msg tea.Msg) tea.Msg {
	switch msg := msg.(type) {
	case SequenceMsg:
		var validCmds []tea.Cmd
		for _, c := range msg {
			if c == nil {
				continue
			}
			validCmds = append(validCmds, c)
		}
		if len(validCmds) == 0 {
			return nil
		}
		return tea.Sequence(validCmds...)()
	case BatchMsg:
		var validCmds []tea.Cmd
		for _, c := range msg {
			if c == nil {
				continue
			}
			validCmds = append(validCmds, c)
		}
		if len(validCmds) == 0 {
			return nil
		}
		return tea.Batch(validCmds...)()
	}
	return msg
}

type IDMsg struct{}

// IDModel is a model that has an UUID.
type IDModel interface {
	ID() uuid.UUID

	tea.Model
}

// NestedMsg is a message that represents a nested tea.Msg.
type NestedMsg struct {
	ID  uuid.UUID
	Msg tea.Msg
}

// Nest nests a tea.Cmd to a specific model.
func Nest(model IDModel, cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	return func() tea.Msg {
		msg := cmd()
		switch msg := msg.(type) {
		case SequenceMsg:
			var cmds SequenceMsg
			for _, cmd := range msg {
				cmds = append(cmds, Nest(model, cmd))
			}
			return cmds
		case BatchMsg:
			var cmds BatchMsg
			for _, cmd := range msg {
				cmds = append(cmds, Nest(model, cmd))
			}
			return cmds
		case tea.BatchMsg:
			var cmds BatchMsg
			for _, cmd := range msg {
				cmds = append(cmds, Nest(model, cmd))
			}
			return cmds
		}
		return NestedMsg{
			ID:  model.ID(),
			Msg: msg,
		}
	}
}

// NestedModel is a model that has nested models.
type NestedModel interface {
	IDModels() []IDModel
	UpdateNestedMsg(NestedMsg) (NestedModel, tea.Cmd)

	tea.Model
}

// InitNested initializes nested models.
func InitNested(m NestedModel) tea.Cmd {
	var batch []tea.Cmd
	models := m.IDModels()
	if len(models) == 0 {
		return nil
	}
	for _, model := range models {
		batch = append(batch, Nest(model, model.Init()))
	}
	return Batch(batch...)
}

// UpdateNested passed the dedicated tea.Msg to nested model.
func UpdateNestedModel(m NestedModel, msg tea.Msg) (NestedModel, tea.Cmd) {
	switch msg := msg.(type) {
	case NestedMsg:
		return m.UpdateNestedMsg(msg)
	default:
		return m, nil
	}
}

// TickModel is a model that updates on a tick.
type TickModel interface {
	tea.Model

	FrameDuration() time.Duration
}

// TickMsg is a message that represents a tick.
type TickMsg time.Time

// Tick returns a tea.Cmd that represents a tick.
func Tick(m TickModel) tea.Cmd {
	return tea.Tick(m.FrameDuration(), func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// TickUpdate updates the model with a tick.
func UpdateTickModel(m TickModel, msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case TickMsg:
		return Tick(m)
	}
	return nil
}

// CmdModel is a model that has a channel of tea.Cmd.
type CmdModel interface {
	CmdChan() <-chan tea.Cmd

	tea.Model
}

// Cmd returns a tea.Cmd that represents a command from the channel.
func Cmd(m CmdModel) tea.Cmd {
	return func() tea.Msg {
		cmd := <-m.CmdChan()
		return cmd()
	}
}

// FocusModel is a model that can be focused.
type FocusModel interface {
	IDModel

	Focused() bool
	Focus() tea.Cmd
	Blur() tea.Cmd
}
