package executors

import (
	"context"
	"errors"
	"testing"

	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialExecutor_Execute(t *testing.T) {
	tests := []struct {
		name    string
		tasks   []runner.Task
		wantErr bool
		verify  func(t *testing.T, tasks []runner.Task)
	}{
		{
			name: "execute successfully",
			tasks: []runner.Task{
				&mocks.MockTask{NameVal: "task1"},
				&mocks.MockTask{NameVal: "task2"},
			},
			verify: func(t *testing.T, tasks []runner.Task) {
				for _, task := range tasks {
					mt := task.(*mocks.MockTask)
					assert.Equal(t, 1, mt.ExecuteCalled, "task %s should be executed once", mt.Name())
				}
			},
		},
		{
			name: "abort on error",
			tasks: []runner.Task{
				&mocks.MockTask{NameVal: "task1"},
				&mocks.MockTask{NameVal: "task2", ExecuteErr: errors.New("failed")},
				&mocks.MockTask{NameVal: "task3"},
			},
			wantErr: true,
			verify: func(t *testing.T, tasks []runner.Task) {
				t1 := tasks[0].(*mocks.MockTask)
				assert.Equal(t, 1, t1.ExecuteCalled)

				t2 := tasks[1].(*mocks.MockTask)
				assert.Equal(t, 1, t2.ExecuteCalled)

				t3 := tasks[2].(*mocks.MockTask)
				assert.Equal(t, 0, t3.ExecuteCalled, "task3 should not be executed")
			},
		},
		{
			name:  "empty tasks",
			tasks: []runner.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewSerial()
			err := exec.Execute(context.Background(), tt.tasks)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.verify != nil {
				tt.verify(t, tt.tasks)
			}
		})
	}
}
