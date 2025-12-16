package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialize_Blob(t *testing.T) {
	data := []byte("hello world")
	blob := NewBlob(data)

	serialized := blob.Serialize()

	assert.Contains(t, string(serialized), "blob")
	assert.Contains(t, string(serialized), "11")
	assert.Contains(t, string(serialized), "hello world")
}

func TestSerialize_Tree(t *testing.T) {
	data := []byte("tree content")
	tree := NewTree(data)

	serialized := tree.Serialize()

	assert.Contains(t, string(serialized), "tree")
	assert.Contains(t, string(serialized), "12")
}

func TestBaseObject_Type(t *testing.T) {
	blob := NewBlob([]byte("test"))

	assert.Equal(t, ObjectTypeBlob, blob.Type())
	assert.True(t, blob.IsBlob())
	assert.False(t, blob.IsTree())
	assert.False(t, blob.IsCommit())
}

func TestBaseObject_Size(t *testing.T) {
	data := []byte("test data")
	blob := NewBlob(data)

	assert.Equal(t, len(data), blob.Size())
	assert.Equal(t, 9, blob.Size())
}

func TestBaseObject_Data(t *testing.T) {
	data := []byte("test data")
	blob := NewBlob(data)

	assert.Equal(t, data, blob.Data())
}

func TestDeserializeObject_ValidBlob(t *testing.T) {
	original := NewBlob([]byte("hello"))
	serialized := original.Serialize()

	deserialized, err := DeserializeObject(serialized)

	require.NoError(t, err)
	assert.NotNil(t, deserialized)
	assert.Equal(t, ObjectTypeBlob, deserialized.Type())
	assert.Equal(t, 5, deserialized.Size())
	assert.Equal(t, []byte("hello"), deserialized.Data())
	assert.True(t, deserialized.IsBlob())
}

func TestDeserializeObject_ValidTree(t *testing.T) {
	original := NewTree([]byte("tree data"))
	serialized := original.Serialize()

	deserialized, err := DeserializeObject(serialized)

	require.NoError(t, err)
	assert.NotNil(t, deserialized)
	assert.Equal(t, ObjectTypeTree, deserialized.Type())
	assert.Equal(t, 9, deserialized.Size())
	assert.True(t, deserialized.IsTree())
}

func TestDeserializeObject_NoNullByte(t *testing.T) {
	invalidData := []byte("blob 5 hello")

	_, err := DeserializeObject(invalidData)

	assert.ErrorIs(t, err, ErrNoNullByteFound)
}

func TestDeserializeObject_SizeMismatch(t *testing.T) {
	invalidData := []byte("blob 10\x00hello")

	_, err := DeserializeObject(invalidData)

	assert.ErrorIs(t, err, ErrObjectSizeMismatch)
}

func TestDeserializeObject_NoSpaceInHeader(t *testing.T) {
	invalidData := []byte("blob5\x00hello")

	_, err := DeserializeObject(invalidData)

	assert.ErrorIs(t, err, ErrNoSpaceInHeader)
}

func TestDeserializeObject_UnknownType(t *testing.T) {
	invalidData := []byte("unknown 5\x00hello")

	_, err := DeserializeObject(invalidData)

	assert.ErrorIs(t, err, ErrUnknownObjectType)
}

func TestDeserializeObject_InvalidSizeFormat(t *testing.T) {
	invalidData := []byte("blob abc\x00hello")

	_, err := DeserializeObject(invalidData)

	assert.ErrorIs(t, err, ErrInvalidSizeFormat)
}

func TestSerializeDeserializeObject_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		typ  ObjectType
	}{
		{"empty blob", []byte{}, ObjectTypeBlob},
		{"small blob", []byte("test"), ObjectTypeBlob},
		{"large blob", []byte("this is a much larger content for testing"), ObjectTypeBlob},
		{"blob with newlines", []byte("line1\nline2\nline3"), ObjectTypeBlob},
		{"tree", []byte("tree content"), ObjectTypeTree},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var original IObject
			if tt.typ == ObjectTypeBlob {
				original = NewBlob(tt.data)
			} else {
				original = NewTree(tt.data)
			}

			serialized := original.Serialize()
			deserialized, err := DeserializeObject(serialized)

			require.NoError(t, err)
			assert.Equal(t, original.Type(), deserialized.Type())
			assert.Equal(t, original.Size(), deserialized.Size())
			assert.Equal(t, original.Data(), deserialized.Data())
		})
	}
}

func TestDeserializeObject_EmptyData(t *testing.T) {
	var emptyData []byte

	_, err := DeserializeObject(emptyData)

	assert.ErrorIs(t, err, ErrNoNullByteFound)
}

func TestSerialize_Format(t *testing.T) {
	blob := NewBlob([]byte("hello"))
	serialized := blob.Serialize()

	expected := "blob 5\x00hello"
	assert.Equal(t, expected, string(serialized))
}
