package task

import (
	"fmt"
	"testing"

	"github.com/KEINOS/go-todotxt/todo/internal/todoparse"
)

// Various valid task strings for benchmarking.
var dataValidTasks = []string{
	"buy milk",
	"(A) buy milk @store +shopping",
	"x (B) 2024-01-15 2024-01-01 complex task @office @home +work +personal due:2024-12-25 note:urgent # comment",
	"(Z) simple task",
	"x (A) 2024-01-01 completed with priority",
	"meeting @conference +bigProject due:2024-06-30 status:pending priority:high",
}

func BenchmarkNew(b *testing.B) {
	for index, taskStr := range dataValidTasks {
		// smoke test
		_, err := New(taskStr)
		if err != nil {
			b.Fatal(err)
		}

		title := fmt.Sprintf("Test #%d: %s", index+1, taskStr)

		b.ResetTimer()
		b.Run(title, func(b *testing.B) {
			for range b.N {
				_, _ = New(taskStr)
			}
		})
	}
}

// Benchmark to compare task.New vs todoparse.FromTaskString to see the overhead difference.
func Benchmark_task_vs_parse(b *testing.B) {
	for index, taskStr := range dataValidTasks {
		title := fmt.Sprintf("Task #%d", index+1)

		b.Run("[todoparse.FromTaskString]: "+title, func(b *testing.B) {
			for range b.N {
				_, _ = todoparse.FromTaskString(taskStr)
			}
		})

		b.Run("[task.New]: "+title, func(b *testing.B) {
			for range b.N {
				_, _ = New(taskStr)
			}
		})
	}
}
