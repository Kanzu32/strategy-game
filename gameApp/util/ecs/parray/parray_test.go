package parray

import (
	"strategy-game/util/ecs/psize"
	"testing"
)

func TestCreatePageArray(t *testing.T) {
	tests := []struct {
		name     string
		pageSize psize.PageSize
		wantSize int
	}{
		{"Small page size", psize.Page16, 65536},
		{"Medium page size", psize.Page256, 65536},
		{"Large page size", psize.Page1024, 65536},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pa := CreatePageArray(tt.pageSize)

			if pa.pageSize != uint16(tt.pageSize) {
				t.Errorf("CreatePageArray() pageSize = %v, want %v", pa.pageSize, tt.pageSize)
			}
			if pa.arraySize != tt.wantSize {
				t.Errorf("CreatePageArray() arraySize = %v, want %v", pa.arraySize, tt.wantSize)
			}
			if len(pa.data) != tt.wantSize/int(tt.pageSize) {
				t.Errorf("CreatePageArray() data length = %v, want %v", len(pa.data), tt.wantSize/int(tt.pageSize))
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		pageSize psize.PageSize
		want     int
	}{
		{"Size 32", psize.Page32, 65536},
		{"Size 128", psize.Page128, 65536},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pa := CreatePageArray(tt.pageSize)
			if got := pa.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAndGet(t *testing.T) {
	tests := []struct {
		name     string
		pageSize psize.PageSize
		index    uint16
		value    int
	}{
		{"First element", psize.Page32, 0, 42},
		{"Middle element", psize.Page32, 20, 123},
		{"Last element", psize.Page32, 31, -1},
		{"Cross-page boundary", psize.Page32, 200, 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pa := CreatePageArray(tt.pageSize)

			// Test initial get (should return -1 for unset values)
			if got := pa.Get(tt.index); got != -1 {
				t.Errorf("Initial Get() = %v, want -1", got)
			}

			// Test Set followed by Get
			pa.Set(tt.index, tt.value)
			if got := pa.Get(tt.index); got != tt.value {
				t.Errorf("After Set(), Get() = %v, want %v", got, tt.value)
			}

			// Verify page creation
			pageNumber := tt.index / uint16(tt.pageSize)
			if pa.data[pageNumber] == nil {
				t.Error("Set() should create the page")
			}
		})
	}
}

func TestPageInitialization(t *testing.T) {
	pa := CreatePageArray(psize.Page64)
	testIndex := uint16(42) // Arbitrary index

	pa.Set(testIndex, 100)

	pageNumber := testIndex / uint16(psize.Page64)
	page := pa.data[pageNumber]

	// Verify page initialization pattern
	for i := 0; i < len(page); i++ {
		if i == int(testIndex%uint16(psize.Page64)) {
			if page[i] != 100 {
				t.Errorf("Set value not preserved, got %d want 100", page[i])
			}
		} else if page[i] != -1 {
			t.Errorf("Uninitialized page element not -1, got %d at index %d", page[i], i)
		}
	}
}

func TestString(t *testing.T) {
	pa := CreatePageArray(psize.Page1024)
	pa.Set(0, 1)
	pa.Set(1023, 2)

	str := pa.String()
	if str == "" {
		t.Error("String() returned empty string")
	}

	// Basic sanity checks of string output
	expectedParts := []string{
		"Page size: 1024",
		"Array size: 65536",
		"[[1 -1 -1", // Beginning of first page
		"2]",        // End of first page
	}

	for _, part := range expectedParts {
		if !contains(str, part) {
			t.Errorf("String() output missing expected part: %q", part)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
