package parse

import (
	"testing"
)

func BenchmarkFromTaskString(b *testing.B) {
	// goos: darwin
	// goarch: arm64
	// pkg: github.com/KEINOS/go-todotxt/todo/internal/parse
	// cpu: Apple M4
	// 4,084,708 ops/sec, 295.8 ns/op, 368 B/op, 3 allocs/op
	for _, testStr := range dataValidVarious {
		b.ResetTimer()

		for range b.N {
			_, _ = FromTaskString(testStr.input)
		}
	}
}

func BenchmarkParsed_reuse_parsed_segments(b *testing.B) {
	testString := "x (A) 2016-05-20 2016-05-18 Thank Mom for the meatballs " +
		"+dinner @phone due:2016-05-25 # comment"

	parsed, err := FromTaskString(testString)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	// goos: darwin
	// goarch: arm64
	// pkg: github.com/KEINOS/go-todotxt/todo/internal/parse
	// cpu: Apple M4
	// 12,493,983 ops/sec, 96.13 ns/op, 32 B/op, 2 allocs/op
	for range b.N {
		_ = parsed.HasInlineComment()
		_ = parsed.IsCommentLine()
		_ = parsed.IsDone()
		_ = parsed.Priority()
		_ = parsed.DateCompleted()
		_ = parsed.DateCreated()
		_ = parsed.Contexts()
		_ = parsed.Projects()
		_ = parsed.Description()
		_ = parsed.KeyValues()
	}
}

func BenchmarkParsed_cache(b *testing.B) {
	testString := "#x (A) 2016-05-20 2016-04-30 measure space for " +
		"+chapelShelving @chapel due:2016-05-30 location:mainOffice"

	parsed, err := FromTaskString(testString)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for range b.N {
		_ = parsed.IsDone()
		_ = parsed.IsCommentLine()
		_ = parsed.consumeHeadParts()
	}
}

func BenchmarkParsed_parse_and_new(b *testing.B) {
	b.Run("create new", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			for _, testStr := range dataValidVarious {
				_, _ = FromTaskString(testStr.input)
			}
		}
	})

	b.Run("parse existing", func(b *testing.B) {
		parsed := new(Parsed)

		b.ResetTimer()

		for range b.N {
			for _, testStr := range dataValidVarious {
				// set original text and update/re-parse
				_ = parsed.Parse(testStr.input)
			}
		}
	})
}
