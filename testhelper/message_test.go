package helper

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestMessage(t *testing.T) {
	oldTrace := Trace
	Trace = true
	defer func() {
		Trace = oldTrace
	}()

	type args struct {
		t     TestingT
		msg   string
		lines []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy path",
			args: func() args {
				mock := gomock.NewController(t)
				mockT := NewMockTestingT(mock)
				mockT.EXPECT().Helper().Return()
				mockT.EXPECT().Name().Return("TestMessage")
				return args{
					t:   mockT,
					msg: "happy path",
				}
			}(),
			want: `
Message: happy path
Case: TestMessage
Error Trace:
	/Users/schen07/Devs/shangkuei/gap/testhelper/message_test.go:47 testhelper.TestMessage.func3`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Message(tt.args.t, tt.args.msg, tt.args.lines...); got != tt.want {
				t.Errorf("Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
