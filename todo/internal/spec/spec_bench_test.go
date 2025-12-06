package spec

import "testing"

// ============================================================================
//  Implemtation examples for IsPrintable()
// ============================================================================

// isPrintableCurrentImpl is the current implementation using two comparisons.
func isPrintableCurrentImpl(m byte) bool {
	return m >= 0x20 && m <= 0x7E
}

// isPrintableSingleCompare uses subtraction to reduce to single comparison.
// The trick: (m - 0x20) will underflow for m < 0x20, making it > 0x5E.
func isPrintableSingleCompare(m byte) bool {
	return (m - 0x20) <= 0x5E // 0x5E = 0x7E - 0x20 = 94
}

// isPrintableBitMask uses bit manipulation approach.
// For printable ASCII (0x20-0x7E), we check if the value minus 0x20
// fits within 95 values (0-94).
func isPrintableBitMask(m byte) bool {
	// Same as single compare but written differently for clarity
	return m-0x20 < 95
}

// isPrintableLookupTable uses a precomputed lookup table.
// Trade memory for speed - 256 bytes for O(1) lookup.
var printableTable = func() [256]bool {
	var table [256]bool
	for i := byte(0x20); i <= 0x7E; i++ {
		table[i] = true
	}

	return table
}()

func isPrintableLookupTable(m byte) bool {
	return printableTable[m]
}

// ----------------------------------------------------------------------------
//  Benchmarks
// ----------------------------------------------------------------------------

// BenchmarkIsPrintable_Current benchmarks the current two-comparison approach.
func BenchmarkIsPrintable_Current(b *testing.B) {
	// Test with various byte values to simulate real usage
	testValues := []byte{0x00, 0x1F, 0x20, 0x40, 0x7E, 0x7F, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, val := range testValues {
			_ = isPrintableCurrentImpl(val)
		}
	}
}

// BenchmarkIsPrintable_SingleCompare benchmarks the single comparison approach.
func BenchmarkIsPrintable_SingleCompare(b *testing.B) {
	testValues := []byte{0x00, 0x1F, 0x20, 0x40, 0x7E, 0x7F, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, val := range testValues {
			_ = isPrintableSingleCompare(val)
		}
	}
}

// BenchmarkIsPrintable_BitMask benchmarks the bit mask approach.
func BenchmarkIsPrintable_BitMask(b *testing.B) {
	testValues := []byte{0x00, 0x1F, 0x20, 0x40, 0x7E, 0x7F, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, val := range testValues {
			_ = isPrintableBitMask(val)
		}
	}
}

// BenchmarkIsPrintable_LookupTable benchmarks the lookup table approach.
func BenchmarkIsPrintable_LookupTable(b *testing.B) {
	testValues := []byte{0x00, 0x1F, 0x20, 0x40, 0x7E, 0x7F, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, val := range testValues {
			_ = isPrintableLookupTable(val)
		}
	}
}

// ----------------------------------------------------------------------------
//  Benchmarks for Mark methods
// ----------------------------------------------------------------------------

// BenchmarkMark_Byte benchmarks the Byte() method (includes validation).
func BenchmarkMark_Byte(b *testing.B) {
	marks := []Mark{0x00, 0x20, 0x40, 0x7E, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, mark := range marks {
			_ = mark.Byte()
		}
	}
}

// BenchmarkMark_Rune benchmarks the Rune() method (includes validation).
func BenchmarkMark_Rune(b *testing.B) {
	marks := []Mark{0x00, 0x20, 0x40, 0x7E, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, mark := range marks {
			_ = mark.Rune()
		}
	}
}

// BenchmarkMark_String benchmarks the String() method (includes validation).
func BenchmarkMark_String(b *testing.B) {
	marks := []Mark{0x00, 0x20, 0x40, 0x7E, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, mark := range marks {
			_ = mark.String()
		}
	}
}

// BenchmarkMark_IsPrintable benchmarks the IsPrintable() method.
func BenchmarkMark_IsPrintable(b *testing.B) {
	marks := []Mark{0x00, 0x20, 0x40, 0x7E, 0xFF}

	b.ResetTimer()

	for range b.N {
		for _, mark := range marks {
			_ = mark.IsPrintable()
		}
	}
}

// ----------------------------------------------------------------------------
//  Verification tests for correctness
// ----------------------------------------------------------------------------

// TestIsPrintableImplementations verifies all implementations produce same results.
func TestIsPrintableImplementations(t *testing.T) {
	t.Parallel()

	for i := range 256 {
		val := byte(i)
		expected := isPrintableCurrentImpl(val)

		if got := isPrintableSingleCompare(val); got != expected {
			t.Errorf("SingleCompare(0x%02X): got %v, want %v", val, got, expected)
		}

		if got := isPrintableBitMask(val); got != expected {
			t.Errorf("BitMask(0x%02X): got %v, want %v", val, got, expected)
		}

		if got := isPrintableLookupTable(val); got != expected {
			t.Errorf("LookupTable(0x%02X): got %v, want %v", val, got, expected)
		}
	}
}

func BenchmarkIsDate(b *testing.B) {
	tests := []string{
		// Valid dates
		"2024-01-01",
		"2024-12-31",
		"2024-02-29", // leap year
		"1234-56-78", // validation is off. valid XXXX-XX-XX format

		// Invalid format
		"20240101",
		"24-01-01",
		"2024/01/01",
		"2024-1-1",

		// Invalid dates
		"2024-02-30",
		"2024-13-01",
		"2024-00-10",
		"2024-04-31",

		// Non-date strings
		"date-is-10", // XXXX-XX-XX but not a date
		"hello-world",
		"1234567890",
		"",
	}

	b.Run("validation off", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			for _, input := range tests {
				_ = IsDate(input, false)
			}
		}
	})

	b.Run("validation on", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			for _, input := range tests {
				_ = IsDate(input, true)
			}
		}
	})
}

func BenchmarkIsPriorityLetter(b *testing.B) {
	fn1 := func(letter string) bool {
		if len(letter) != 1 {
			return false
		}

		ch := letter[0]
		return ch >= 'A' && ch <= 'Z'
	}

	letters := []string{"A", "M", "Z", "a", "1", "", "AA"}

	b.Run("Original", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			for _, letter := range letters {
				_ = IsPriorityLetter(letter)
			}
		}
	})

	b.Run("Fn1", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			for _, letter := range letters {
				_ = fn1(letter)
			}
		}
	})
}
