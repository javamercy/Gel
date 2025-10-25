package objects

import (
	"Gel/src/gel/core/constant"
	"testing"
)

// TestBaseObjectType tests the Type() method of BaseObject
func TestBaseObjectType(t *testing.T) {
	tests := []struct {
		name         string
		objectType   constant.ObjectType
		data         []byte
	}{
		{
			name:       "blob type",
			objectType: constant.Blob,
			data:       []byte("test data"),
		},
		{
			name:       "tree type",
			objectType: constant.Tree,
			data:       []byte("tree data"),
		},
		{
			name:       "commit type",
			objectType: constant.Commit,
			data:       []byte("commit data"),
		},
		{
			name:       "empty data with blob type",
			objectType: constant.Blob,
			data:       []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &BaseObject{
				objectType: tt.objectType,
				data:       tt.data,
			}

			result := obj.Type()

			if result != tt.objectType {
				t.Errorf("Type() = %v, want %v", result, tt.objectType)
			}
		})
	}
}

// TestBaseObjectSize tests the Size() method of BaseObject
func TestBaseObjectSize(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		expectedSize int
	}{
		{
			name:         "empty data",
			data:         []byte{},
			expectedSize: 0,
		},
		{
			name:         "small text data",
			data:         []byte("hello"),
			expectedSize: 5,
		},
		{
			name:         "larger text data",
			data:         []byte("hello world this is a test"),
			expectedSize: 26,
		},
		{
			name:         "binary data",
			data:         []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expectedSize: 5,
		},
		{
			name:         "unicode data",
			data:         []byte("Hello 世界"),
			expectedSize: 12, // UTF-8 encoding
		},
		{
			name:         "nil data",
			data:         nil,
			expectedSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &BaseObject{
				objectType: constant.Blob,
				data:       tt.data,
			}

			result := obj.Size()

			if result != tt.expectedSize {
				t.Errorf("Size() = %d, want %d", result, tt.expectedSize)
			}
		})
	}
}

// TestBaseObjectData tests the Data() method of BaseObject
func TestBaseObjectData(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "text data",
			data: []byte("hello world"),
		},
		{
			name: "binary data",
			data: []byte{0x00, 0x01, 0x02, 0xFF},
		},
		{
			name: "large data",
			data: make([]byte, 1024),
		},
		{
			name: "nil data",
			data: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &BaseObject{
				objectType: constant.Blob,
				data:       tt.data,
			}

			result := obj.Data()

			// Check length
			if len(result) != len(tt.data) {
				t.Errorf("Data() length = %d, want %d", len(result), len(tt.data))
			}

			// Check content
			if tt.data != nil && result != nil {
				for i := range tt.data {
					if result[i] != tt.data[i] {
						t.Errorf("Data()[%d] = %v, want %v", i, result[i], tt.data[i])
						break
					}
				}
			}
		})
	}
}

// TestBaseObjectInterfaceCompliance tests that BaseObject implements IObject interface
func TestBaseObjectInterfaceCompliance(t *testing.T) {
	var _ IObject = (*BaseObject)(nil)
	
	t.Run("BaseObject implements IObject", func(t *testing.T) {
		obj := &BaseObject{
			objectType: constant.Blob,
			data:       []byte("test"),
		}

		// Test all interface methods
		var iobj IObject = obj

		// Type() method
		objType := iobj.Type()
		if objType != constant.Blob {
			t.Errorf("Type() = %v, want %v", objType, constant.Blob)
		}

		// Size() method
		size := iobj.Size()
		if size != 4 {
			t.Errorf("Size() = %d, want 4", size)
		}

		// Data() method
		data := iobj.Data()
		expected := []byte("test")
		if len(data) != len(expected) {
			t.Errorf("Data() length = %d, want %d", len(data), len(expected))
		}
	})
}

// TestBaseObjectImmutability tests that returned data is the actual slice
func TestBaseObjectDataReference(t *testing.T) {
	t.Run("Data() returns reference to internal data", func(t *testing.T) {
		originalData := []byte("original")
		obj := &BaseObject{
			objectType: constant.Blob,
			data:       originalData,
		}

		returnedData := obj.Data()

		// Verify it's the same underlying array (shallow copy check)
		if len(returnedData) != len(originalData) {
			t.Errorf("Data length mismatch")
		}

		for i := range originalData {
			if returnedData[i] != originalData[i] {
				t.Errorf("Data content mismatch at index %d", i)
			}
		}
	})
}

// TestBaseObjectEdgeCases tests edge cases
func TestBaseObjectEdgeCases(t *testing.T) {
	t.Run("zero value BaseObject", func(t *testing.T) {
		var obj BaseObject

		// Should have empty/zero values
		if obj.Type() != "" {
			t.Errorf("Zero value Type() = %v, want empty", obj.Type())
		}

		if obj.Size() != 0 {
			t.Errorf("Zero value Size() = %d, want 0", obj.Size())
		}

		if obj.Data() != nil {
			t.Errorf("Zero value Data() = %v, want nil", obj.Data())
		}
	})

	t.Run("BaseObject with very large data", func(t *testing.T) {
		// Create a large byte slice (1MB)
		largeData := make([]byte, 1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		obj := &BaseObject{
			objectType: constant.Blob,
			data:       largeData,
		}

		if obj.Size() != 1024*1024 {
			t.Errorf("Size() = %d, want %d", obj.Size(), 1024*1024)
		}

		retrievedData := obj.Data()
		if len(retrievedData) != len(largeData) {
			t.Errorf("Data() length = %d, want %d", len(retrievedData), len(largeData))
		}
	})
}
