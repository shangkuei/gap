package bubbles

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	helper "github.com/shangkuei/gap/testhelper"
)

func init() {
	helper.Trace = false
}

var compareTeaCmd = cmp.Comparer(func(a, b tea.Cmd) bool {
	if reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer() {
		return true
	}
	return cmp.Equal(a(), b())
})

func TestSequence(t *testing.T) {
	t.Parallel()

	var msg struct{}
	cmd := func() tea.Msg { return msg }
	tests := []struct {
		name string
		args []tea.Cmd
		want tea.Msg
	}{
		{
			name: "happy path",
			args: []tea.Cmd{cmd, cmd},
			want: SequenceMsg([]tea.Cmd{cmd, cmd}),
		},
		{
			name: "empty",
			args: []tea.Cmd{nil},
		},
		{
			name: "one cmd",
			args: []tea.Cmd{cmd},
			want: msg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sequence(tt.args...)
			if tt.want == nil {
				if diff, ok := helper.Equal(got, tea.Cmd(nil)); !ok {
					t.Error(helper.Message(t, "expected nil tea.Cmd", diff))
				}
			} else {
				if _, ok := helper.Equal(got, tea.Cmd(nil)); ok {
					t.Error(helper.Message(t, "unexpected nil tea.Cmd"))
				}
				if diff, ok := helper.Equal(got(), tt.want, compareTeaCmd); !ok {
					t.Error(helper.Message(t, "unexpected tea.Cmd", diff))
				}
			}
		})
	}
}

func TestBatch(t *testing.T) {
	t.Parallel()

	var msg struct{}
	cmd := func() tea.Msg { return msg }
	tests := []struct {
		name string
		args []tea.Cmd
		want tea.Msg
	}{
		{
			name: "happy path",
			args: []tea.Cmd{cmd, cmd},
			want: BatchMsg([]tea.Cmd{cmd, cmd}),
		},
		{
			name: "empty",
			args: []tea.Cmd{nil},
		},
		{
			name: "one cmd",
			args: []tea.Cmd{cmd},
			want: msg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Batch(tt.args...)
			if tt.want == nil {
				if diff, ok := helper.Equal(got, tea.Cmd(nil)); !ok {
					t.Error(helper.Message(t, "expected nil tea.Cmd", diff))
				}
			} else {
				if _, ok := helper.Equal(got, tea.Cmd(nil)); ok {
					t.Error(helper.Message(t, "unexpected nil tea.Cmd"))
				}
				if diff, ok := helper.Equal(got(), tt.want, compareTeaCmd); !ok {
					t.Error(helper.Message(t, "unexpected tea.Cmd", diff))
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	cmd1, cmd2 := func() tea.Msg { return nil }, func() tea.Msg { return nil }
	msg := tea.Msg(nil)
	type args struct {
		model tea.Model
		msg   tea.Msg
	}
	tests := []struct {
		name string
		args args
		want tea.Msg
	}{
		{
			name: "sequence",
			args: args{msg: SequenceMsg([]tea.Cmd{cmd1, cmd2})},
			want: tea.Sequence(cmd1, cmd2)(),
		},
		{
			name: "empty sequence",
			args: args{msg: SequenceMsg([]tea.Cmd{nil})},
		},
		{
			name: "batch",
			args: args{msg: BatchMsg([]tea.Cmd{cmd1, cmd2})},
			want: tea.Batch(cmd1, cmd2)(),
		},
		{
			name: "empty batch",
			args: args{msg: BatchMsg([]tea.Cmd{nil})},
		},
		{
			name: "default",
			args: args{msg: msg},
			want: msg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.args.model, tt.args.msg)
			if tt.want == nil {
				if diff, ok := helper.Equal(got, tea.Msg(nil)); !ok {
					t.Error(helper.Message(t, "expected nil tea.Msg", diff))
				}
			} else {
				if _, ok := helper.Equal(got, tea.Cmd(nil)); ok {
					t.Error(helper.Message(t, "unexpected nil tea.Msg"))
				}
				if diff, ok := helper.Equal(got, tt.want, compareTeaCmd); !ok {
					t.Error(helper.Message(t, "unexpected tea.Msg", diff))
				}
			}
		})
	}
}

type fieldMsg struct {
	field string
}

type idModel struct {
	id    uuid.UUID
	field string
}

func (m idModel) ID() uuid.UUID {
	return m.id
}

func (m idModel) Init() tea.Cmd {
	return func() tea.Msg { return nil }
}

func (m idModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fieldMsg:
		m.field = msg.field
	}
	return m, nil
}

func (m idModel) View() string {
	return ""
}

func TestNest(t *testing.T) {
	t.Parallel()

	var msg struct{}
	cmd1, cmd2 := func() tea.Msg { return msg }, func() tea.Msg { return msg }
	model := &idModel{id: uuid.New()}
	type args struct {
		model IDModel
		cmd   tea.Cmd
	}
	tests := []struct {
		name string
		args args
		want tea.Msg
	}{
		{
			name: "happy path",
			args: args{model: model, cmd: cmd1},
			want: NestedMsg{ID: model.ID(), Msg: msg},
		},
		{
			name: "empty cmd",
			args: args{model: model},
		},
		{
			name: "sequence",
			args: args{model: model, cmd: Sequence(cmd1, cmd2)},
			want: SequenceMsg{Nest(model, cmd1), Nest(model, cmd2)},
		},
		{
			name: "batch",
			args: args{model: model, cmd: Batch(cmd1, cmd2)},
			want: BatchMsg{Nest(model, cmd1), Nest(model, cmd2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Nest(tt.args.model, tt.args.cmd)
			if tt.want == nil {
				if diff, ok := helper.Equal(got, tea.Cmd(nil)); !ok {
					t.Error(helper.Message(t, "expected nil tea.Cmd", diff))
				}
			} else {
				if _, ok := helper.Equal(got, tea.Cmd(nil)); ok {
					t.Error(helper.Message(t, "unexpected nil tea.Cmd"))
				}
				if diff, ok := helper.Equal(got(), tt.want, compareTeaCmd); !ok {
					t.Error(helper.Message(t, "unexpected tea.Cmd", diff))
				}
			}
		})
	}
}

type nestedModel struct {
	models []IDModel
}

func (m nestedModel) IDModels() []IDModel {
	return m.models
}

func (m nestedModel) UpdateNestedMsg(msg NestedMsg) (NestedModel, tea.Cmd) {
	index := slices.IndexFunc(m.models, func(i IDModel) bool {
		return i.ID() == msg.ID
	})
	if index == -1 {
		return m, nil
	}
	model, cmd := m.models[index].Update(msg.Msg)
	m.models = slices.Replace(m.models, index, index+1, model.(IDModel))
	return m, cmd
}

func (m nestedModel) Init() tea.Cmd {
	return InitNested(m)
}

func (m nestedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return UpdateNestedModel(m, msg)
}

func (m nestedModel) View() string {
	return ""
}

func TestNestedModel(t *testing.T) {
	t.Parallel()

	id := idModel{id: uuid.New()}

	tests := []struct {
		name      string
		model     NestedModel
		wantInit  tea.Cmd
		wantModel NestedModel
	}{
		{
			name:      "happy path",
			model:     nestedModel{models: []IDModel{id}},
			wantInit:  Batch(Nest(id, id.Init())),
			wantModel: nestedModel{models: []IDModel{idModel{id: id.id, field: "field"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.model.Init()
			if diff, ok := helper.Equal(cmd, tt.wantInit, compareTeaCmd); !ok {
				t.Error(helper.Message(t, "expected tea.Cmd", diff))
			}

			model, cmd := tt.model.Update(Nest(id, func() tea.Msg { return fieldMsg{field: "field"} })())
			if diff, ok := helper.Equal(cmd, tea.Cmd(nil)); !ok {
				t.Error(helper.Message(t, "expected nil tea.Cmd", diff))
			}
			model, ok := model.(NestedModel)
			if !ok {
				t.Error(helper.Message(t, "expected NestedModel", fmt.Sprintf("got: %T", model)))
			}
			if diff, ok := helper.Equal(model, tt.wantModel, cmp.AllowUnexported(nestedModel{}), cmp.AllowUnexported(idModel{})); !ok {
				t.Error(helper.Message(t, "unexpected NestedModel", diff))
			}
		})
	}
}

type tickModel struct {
}

func (m tickModel) FrameDuration() time.Duration {
	return time.Second
}

func (m tickModel) Init() tea.Cmd {
	return nil
}

func (m tickModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, Tick(m)
}

func (m tickModel) View() string {
	return ""
}

func TestUpdateTickModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		model TickModel
	}{
		{
			name:  "happy path",
			model: tickModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := UpdateTickModel(tt.model, TickMsg(time.Now()))
			if _, ok := helper.Equal(cmd, tea.Cmd(nil)); ok {
				t.Error(helper.Message(t, "unexpected nil tea.Cmd"))
			}
		})
	}
}

type cmdMsg struct{}

type cmdModel struct {
	cmd   chan tea.Cmd
	field string
}

func (m cmdModel) CmdChan() <-chan tea.Cmd {
	return m.cmd
}

func (m cmdModel) Init() tea.Cmd {
	return nil
}

func (m cmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case cmdMsg:
		m.field = "ok"
	}
	return m, nil
}

func (m cmdModel) View() string {
	return ""
}

func TestCmd(t *testing.T) {
	t.Parallel()

	ch := make(chan tea.Cmd)
	tests := []struct {
		name  string
		model CmdModel
		want  CmdModel
	}{
		{
			name:  "happy path",
			model: cmdModel{cmd: ch},
			want:  cmdModel{cmd: ch, field: "ok"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				ch <- func() tea.Msg { return cmdMsg{} }
			}()
			model, cmd := tt.model.Update(Cmd(tt.model)())
			if diff, ok := helper.Equal(cmd, tea.Cmd(nil)); !ok {
				t.Error(helper.Message(t, "expected nil tea.Cmd", diff))
			}
			if diff, ok := helper.Equal(model, tt.want, cmp.AllowUnexported(cmdModel{})); !ok {
				t.Error(helper.Message(t, "unexpected CmdModel", diff))
			}
		})
	}
}
