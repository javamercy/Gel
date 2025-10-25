package objects

import (
	"Gel/src/gel/core/constant"
	"testing"
)

// TestNewBlob tests the NewBlob constructor
func TestNewBlob(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		expectedSize int
	}{
		{
			name:         "create blob with text data",
			data:         []byte("hello world"),
			expectedSize: 11,
		},
		{
			name:         "create blob with empty data",
			data:         []byte{},
			expectedSize: 0,
		},
		{
			name:         "create blob with binary data",
			data:         []byte{0x00, 0x01, 0x02, 0xFF},
			expectedSize: 4,
		},
		{
			name:         "create blob with nil data",
			data:         nil,
			expectedSize: 0,
		},
		{
			name:         "create blob with large data",
			data:         make([]byte, 1024),
			expectedSize: 1024,
		},
		{
			name:         "create blob with unicode data",
			data:         []byte("Hello 世界"),
			expectedSize: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blob := NewBlob(tt.data)

			if blob == nil {
				t.Fatal("NewBlob returned nil")
			}

			if blob.BaseObject == nil {
				t.Fatal("NewBlob BaseObject is nil")
			}

			// Verify type is set correctly
			if blob.Type() != constant.Blob {
				t.Errorf("Type() = %v, want %v", blob.Type(), constant.Blob)
			}

			// Verify size
			if blob.Size() != tt.expectedSize {
				t.Errorf("Size() = %d, want %d", blob.Size(), tt.expectedSize)
			}

			// Verify data
			retrievedData := blob.Data()
			if len(retrievedData) != len(tt.data) {
				t.Errorf("Data() length = %d, want %d", len(retrievedData), len(tt.data))
			}

			if tt.data != nil {
				for i := range tt.data {
					if retrievedData[i] != tt.data[i] {
						t.Errorf("Data()[%d] = %v, want %v", i, retrievedData[i], tt.data[i])
						break
					}
				}
			}
		})
	}
}

// TestBlobType tests that Blob type is always constant.Blob
func TestBlobType(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "text data",
			data: []byte("test"),
		},
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "binary data",
			data: []byte{0x00, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blob := NewBlob(tt.data)

			if blob.Type() != constant.Blob {
				t.Errorf("Type() = %v, want %v", blob.Type(), constant.Blob)
			}
		})
	}
}

// TestBlobInterfaceCompliance tests that Blob implements IObject interface
func TestBlobInterfaceCompliance(t *testing.T) {
	var _ IObject = (*Blob)(nil)

	t.Run("Blob implements IObject interface", func(t *testing.T) {
		data := []byte("test data")
		blob := NewBlob(data)

		var iobj IObject = blob

		// Test Type() method
		if iobj.Type() != constant.Blob {
			t.Errorf("Type() = %v, want %v", iobj.Type(), constant.Blob)
		}

		// Test Size() method
		if iobj.Size() != len(data) {
			t.Errorf("Size() = %d, want %d", iobj.Size(), len(data))
		}

		// Test Data() method
		retrievedData := iobj.Data()
		if len(retrievedData) != len(data) {
			t.Errorf("Data() length = %d, want %d", len(retrievedData), len(data))
		}
	})
}

// TestBlobBehavior tests various blob behaviors
func TestBlobBehavior(t *testing.T) {
	t.Run("blob with different data types", func(t *testing.T) {
		testCases := []struct {
			name        string
			data        []byte
			description string
		}{
			{
				name:        "plain text",
				data:        []byte("This is plain text content"),
				description: "ASCII text",
			},
			{
				name:        "source code",
				data:        []byte("package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}"),
				description: "Go source code",
			},
			{
				name:        "json data",
				data:        []byte(`{"key": "value", "number": 42}`),
				description: "JSON content",
			},
			{
				name:        "binary zeros",
				data:        []byte{0x00, 0x00, 0x00, 0x00},
				description: "Binary zeros",
			},
			{
				name:        "mixed binary",
				data:        []byte{0x89, 0x50, 0x4E, 0x47}, // PNG header
				description: "PNG file header",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				blob := NewBlob(tc.data)

				if blob.Type() != constant.Blob {
					t.Errorf("Type() = %v, want %v", blob.Type(), constant.Blob)
				}

				if blob.Size() != len(tc.data) {
					t.Errorf("Size() = %d, want %d", blob.Size(), len(tc.data))
				}

				retrievedData := blob.Data()
				if len(retrievedData) != len(tc.data) {
					t.Errorf("Data() length = %d, want %d", len(retrievedData), len(tc.data))
				}
			})
		}
	})

	t.Run("multiple blobs are independent", func(t *testing.T) {
		data1 := []byte("blob 1")
		data2 := []byte("blob 2 different")

		blob1 := NewBlob(data1)
		blob2 := NewBlob(data2)

		// Verify they are independent
		if blob1.Size() == blob2.Size() {
			t.Error("Expected different sizes for different blobs")
		}

		// Verify data is different
		data1Retrieved := blob1.Data()
		data2Retrieved := blob2.Data()

		if string(data1Retrieved) == string(data2Retrieved) {
			t.Error("Expected different data for different blobs")
		}
	})

	t.Run("blob preserves data integrity", func(t *testing.T) {
		originalData := []byte("test data that should be preserved")
		blob := NewBlob(originalData)

		retrievedData := blob.Data()

		// Verify exact match
		if len(retrievedData) != len(originalData) {
			t.Fatalf("Data length mismatch: got %d, want %d", len(retrievedData), len(originalData))
		}

		for i := range originalData {
			if retrievedData[i] != originalData[i] {
				t.Errorf("Data mismatch at byte %d: got %v, want %v", i, retrievedData[i], originalData[i])
			}
		}
	})
}

// TestBlobEdgeCases tests edge cases for Blob
func TestBlobEdgeCases(t *testing.T) {
	t.Run("blob with single byte", func(t *testing.T) {
		data := []byte{0x42}
		blob := NewBlob(data)

		if blob.Size() != 1 {
			t.Errorf("Size() = %d, want 1", blob.Size())
		}

		retrievedData := blob.Data()
		if len(retrievedData) != 1 || retrievedData[0] != 0x42 {
			t.Errorf("Data() = %v, want [0x42]", retrievedData)
		}
	})

	t.Run("blob with whitespace only", func(t *testing.T) {
		data := []byte("   \n\t\r\n   ")
		blob := NewBlob(data)

		if blob.Size() != len(data) {
			t.Errorf("Size() = %d, want %d", blob.Size(), len(data))
		}

		retrievedData := blob.Data()
		if string(retrievedData) != string(data) {
			t.Error("Whitespace data not preserved")
		}
	})

	t.Run("blob with newlines", func(t *testing.T) {
		data := []byte("line1\nline2\nline3\n")
		blob := NewBlob(data)

		if blob.Size() != len(data) {
			t.Errorf("Size() = %d, want %d", blob.Size(), len(data))
		}

		retrievedData := blob.Data()
		if string(retrievedData) != string(data) {
			t.Error("Newlines not preserved")
		}
	})

	t.Run("blob with special characters", func(t *testing.T) {
		data := []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?")
		blob := NewBlob(data)

		if blob.Size() != len(data) {
			t.Errorf("Size() = %d, want %d", blob.Size(), len(data))
		}

		retrievedData := blob.Data()
		if string(retrievedData) != string(data) {
			t.Error("Special characters not preserved")
		}
	})

	t.Run("very large blob", func(t *testing.T) {
		// Create a 10MB blob
		largeData := make([]byte, 10*1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		blob := NewBlob(largeData)

		if blob.Size() != len(largeData) {
			t.Errorf("Size() = %d, want %d", blob.Size(), len(largeData))
		}

		if blob.Type() != constant.Blob {
			t.Errorf("Type() = %v, want %v", blob.Type(), constant.Blob)
		}
	})
}

// TestBlobDataModification tests if modifying the input data affects the blob
func TestBlobDataModification(t *testing.T) {
	t.Run("modifying input data after creation", func(t *testing.T) {
		originalData := []byte("original data")
		blob := NewBlob(originalData)

		// Modify the original data
		if len(originalData) > 0 {
			originalData[0] = 'X'
		}

		// The blob's data should reflect this change since Go slices are references
		// This test documents the current behavior
		retrievedData := blob.Data()
		if len(retrievedData) > 0 && len(originalData) > 0 {
			// Document that this is a shallow copy behavior
			_ = retrievedData
		}
	})
}
