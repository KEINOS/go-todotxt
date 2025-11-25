package segment

import (
	"strings"
	"testing"
	"time"
)

var testFields = []string{"@home", "work", "@office", "task", "@", "a", ""}

func BenchmarkDirectPrefix(b *testing.B) {
	for _, field := range testFields {
		for range b.N {
			_ = len(field) > 1 && field[0] == '@'
		}
	}
}

func BenchmarkHasPrefix(b *testing.B) {
	for _, field := range testFields {
		for range b.N {
			_ = len(field) > 1 && strings.HasPrefix(field, "@")
		}
	}
}

var testDates = []string{"2016-05-20", "2023-12-31", "0000-01-01", "9999-12-31", "invalid", "2016-13-01", "2016-05-32"}

func BenchmarkIsDate(b *testing.B) {
	for _, date := range testDates {
		for range b.N {
			_ = Segment(date).IsDate()
		}
	}
}

func BenchmarkTimeParse(b *testing.B) {
	for _, date := range testDates {
		for range b.N {
			_, err := time.Parse("2006-01-02", date)
			_ = err == nil
		}
	}
}
