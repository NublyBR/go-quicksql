package quicksql

import (
	"testing"
)

func TestQuote(t *testing.T) {
	// Test cases for nil data
	t.Run("NilData", func(t *testing.T) {
		want := "NULL"
		got := Quote(nil)
		if got != want {
			t.Errorf("Quote(nil) = %s; want %s", got, want)
		}
	})

	// Test cases for slice of bytes
	t.Run("BytesSlice", func(t *testing.T) {
		data := []byte{65, 66, 67}
		want := "0x414243"
		got := Quote(data)
		if got != want {
			t.Errorf("Quote(%v) = %s; want %s", data, got, want)
		}
	})

	// Test cases for slice of runes
	t.Run("RunesSlice", func(t *testing.T) {
		data := []rune{'A', 'B', 'C'}
		want := "'ABC'"
		got := Quote(data)
		if got != want {
			t.Errorf("Quote(%v) = %s; want %s", data, got, want)
		}
	})

	// Test cases for string
	t.Run("String", func(t *testing.T) {
		data := "Hello, World!"
		want := "'Hello, World!'"
		got := Quote(data)
		if got != want {
			t.Errorf("Quote(%v) = %s; want %s", data, got, want)
		}
	})

	// Test cases for unsigned integers
	t.Run("UnsignedIntegers", func(t *testing.T) {
		data := uint32(42)
		want := "42"
		got := Quote(data)
		if got != want {
			t.Errorf("Quote(%v) = %s; want %s", data, got, want)
		}
	})

	// Test cases for signed integers
	t.Run("SignedIntegers", func(t *testing.T) {
		data := int64(-123)
		want := "-123"
		got := Quote(data)
		if got != want {
			t.Errorf("Quote(%v) = %s; want %s", data, got, want)
		}
	})

	// Test cases for floating point numbers
	t.Run("FloatingPoint", func(t *testing.T) {
		data := 3.14
		want := "3.140000"
		got := Quote(data)
		if got != want {
			t.Errorf("Quote(%v) = %s; want %s", data, got, want)
		}
	})
}
