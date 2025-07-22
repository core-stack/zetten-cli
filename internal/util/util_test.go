package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInSlice(t *testing.T) {
	t.Run("find existing element", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		val, found := FindInSlice(slice, func(x int) bool { return x == 3 })
		assert.True(t, found)
		assert.Equal(t, 3, val)
	})

	t.Run("find non-existing element", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		val, found := FindInSlice(slice, func(s string) bool { return s == "d" })
		assert.False(t, found)
		assert.Equal(t, "", val)
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []float64{}
		val, found := FindInSlice(slice, func(f float64) bool { return f > 0 })
		assert.False(t, found)
		assert.Equal(t, 0.0, val)
	})

	t.Run("find struct element", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		people := []Person{
			{"Alice", 30},
			{"Bob", 25},
		}
		val, found := FindInSlice(people, func(p Person) bool { return p.Name == "Bob" })
		assert.True(t, found)
		assert.Equal(t, "Bob", val.Name)
		assert.Equal(t, 25, val.Age)
	})
}

func TestFindInMap(t *testing.T) {
	t.Run("find existing value", func(t *testing.T) {
		m := map[string]int{
			"one": 1,
			"two": 2,
		}
		val, found := FindInMap(m, func(k string, v int) bool { return k == "two" })
		assert.True(t, found)
		assert.Equal(t, 2, val)
	})

	t.Run("find by value", func(t *testing.T) {
		m := map[int]string{
			1: "one",
			2: "two",
		}
		val, found := FindInMap(m, func(k int, v string) bool { return v == "one" })
		assert.True(t, found)
		assert.Equal(t, "one", val)
	})

	t.Run("non-existing key", func(t *testing.T) {
		m := map[string]bool{
			"true":  true,
			"false": false,
		}
		val, found := FindInMap(m, func(k string, v bool) bool { return k == "maybe" })
		assert.False(t, found)
		assert.False(t, val)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[float64]string{}
		val, found := FindInMap(m, func(k float64, v string) bool { return true })
		assert.False(t, found)
		assert.Equal(t, "", val)
	})
}

func TestMergeMap(t *testing.T) {
	t.Run("merge two maps", func(t *testing.T) {
		m1 := map[string]int{"a": 1, "b": 2}
		m2 := map[string]int{"b": 3, "c": 4}
		merged := MergeMap(m1, m2)
		assert.Equal(t, 3, len(merged))
		assert.Equal(t, 1, merged["a"])
		assert.Equal(t, 2, merged["b"]) // m1 takes precedence
		assert.Equal(t, 4, merged["c"])
	})

	t.Run("merge three maps", func(t *testing.T) {
		m1 := map[int]string{1: "one"}
		m2 := map[int]string{2: "two"}
		m3 := map[int]string{3: "three"}
		merged := MergeMap(m1, m2, m3)
		assert.Equal(t, 3, len(merged))
	})

	t.Run("merge empty maps", func(t *testing.T) {
		merged := MergeMap(map[string]bool{}, map[string]bool{})
		assert.Empty(t, merged)
	})

	t.Run("merge with empty map", func(t *testing.T) {
		m := map[string]float64{"pi": 3.14}
		merged := MergeMap(m, map[string]float64{})
		assert.Equal(t, 1, len(merged))
		assert.Equal(t, 3.14, merged["pi"])
	})
}

func TestOr(t *testing.T) {
	t.Run("first non-empty", func(t *testing.T) {
		result := Or("", "second", "third")
		assert.Equal(t, "second", result)
	})

	t.Run("all non-empty", func(t *testing.T) {
		result := Or("first", "second")
		assert.Equal(t, "first", result)
	})

	t.Run("with whitespace", func(t *testing.T) {
		result := Or("  ", "\t", "valid")
		assert.Equal(t, "valid", result)
	})

	t.Run("all empty", func(t *testing.T) {
		result := Or("", "", "")
		assert.Equal(t, "", result)
	})

	t.Run("single empty", func(t *testing.T) {
		result := Or("")
		assert.Equal(t, "", result)
	})

	t.Run("single non-empty", func(t *testing.T) {
		result := Or("only")
		assert.Equal(t, "only", result)
	})
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://example.com/path?query=param", true},
		{"valid with port", "http://localhost:8080", true},
		{"missing scheme", "example.com", false},
		{"missing host", "http://", false},
		{"invalid scheme", "ftp://example.com", true},
		{"invalid format", "not a url", false},
		{"empty string", "", false},
		{"with basic auth", "http://user:pass@example.com", true},
		{"with fragment", "https://example.com/#section", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidURL(tt.url))
		})
	}
}
