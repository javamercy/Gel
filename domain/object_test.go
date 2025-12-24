package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialize_Blob(t *testing.T) {
	body := []byte("hello world")
	blob := NewBlob(body)

	serializedBlob := blob.Serialize()

	assert.Contains(t, string(serializedBlob), "blob")
	assert.Contains(t, string(serializedBlob), "11")
	assert.Contains(t, string(serializedBlob), "hello world")
}

func TestSerialize_Tree(t *testing.T) {
	body := []byte("tree content")
	tree := NewTree(body)

	serializedTree := tree.Serialize()

	assert.Contains(t, string(serializedTree), "tree")
	assert.Contains(t, string(serializedTree), "12")
}

func TestBaseObject_Type(t *testing.T) {
	blob := NewBlob([]byte("test"))

	assert.Equal(t, ObjectTypeBlob, blob.Type())
}

func TestBaseObject_Size(t *testing.T) {
	body := []byte("test data")
	blob := NewBlob(body)

	assert.Equal(t, len(body), blob.Size())
	assert.Equal(t, 9, blob.Size())
}

func TestBaseObject_Data(t *testing.T) {
	body := []byte("test data")
	blob := NewBlob(body)

	assert.Equal(t, body, blob.Body())
}

func TestDeserializeObject_ValidBlob(t *testing.T) {
	body := []byte("hello")
	blob := NewBlob(body)
	serializedBlob := blob.Serialize()

	object, err := DeserializeObject(serializedBlob)
	require.NoError(t, err)
	assert.NotNil(t, object)

	deserializedBlob, ok := object.(*Blob)
	require.True(t, ok)
	assert.Equal(t, ObjectTypeBlob, deserializedBlob.Type())
	assert.Equal(t, len(body), deserializedBlob.Size())
	assert.Equal(t, body, deserializedBlob.Body())
}

func TestDeserializeObject_ValidTree(t *testing.T) {
	body := []byte("tree data")
	tree := NewTree(body)
	serializedTree := tree.Serialize()

	object, err := DeserializeObject(serializedTree)
	require.NoError(t, err)
	assert.NotNil(t, object)

	deserializedTree, ok := object.(*Tree)
	require.True(t, ok)
	assert.Equal(t, len(body), deserializedTree.Size())
}

func TestDeserializeObject_NoNullByte(t *testing.T) {
	invalidBody := []byte("blob 5 hello")

	_, err := DeserializeObject(invalidBody)

	assert.ErrorIs(t, err, ErrNoNullByteFound)
}

func TestDeserializeObject_SizeMismatch(t *testing.T) {
	invalidBody := []byte("blob 10\x00hello")

	_, err := DeserializeObject(invalidBody)

	assert.ErrorIs(t, err, ErrObjectSizeMismatch)
}

func TestDeserializeObject_NoSpaceInHeader(t *testing.T) {
	invalidBody := []byte("blob5\x00hello")

	_, err := DeserializeObject(invalidBody)

	assert.ErrorIs(t, err, ErrNoSpaceInHeader)
}

func TestDeserializeObject_UnknownType(t *testing.T) {
	invalidBody := []byte("unknown 5\x00hello")

	_, err := DeserializeObject(invalidBody)

	assert.ErrorIs(t, err, ErrUnknownObjectType)
}

func TestDeserializeObject_InvalidSizeFormat(t *testing.T) {
	invalidBody := []byte("blob abc\x00hello")

	_, err := DeserializeObject(invalidBody)

	assert.ErrorIs(t, err, ErrInvalidSizeFormat)
}

func TestSerializeDeserializeObject_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		body  []byte
		type_ ObjectType
	}{
		{"empty blob", []byte{}, ObjectTypeBlob},
		{"small blob", []byte("test"), ObjectTypeBlob},
		{"large blob", []byte("this is a much larger content for testing"), ObjectTypeBlob},
		{"blob with newlines", []byte("line1\nline2\nline3"), ObjectTypeBlob},
		{"tree", []byte("tree content"), ObjectTypeTree},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var object IObject
			if tt.type_ == ObjectTypeBlob {
				object = NewBlob(tt.body)
			} else {
				object = NewTree(tt.body)
			}

			serializedObject := object.Serialize()
			deserializedObject, err := DeserializeObject(serializedObject)

			require.NoError(t, err)
			assert.Equal(t, object.Type(), deserializedObject.Type())
			assert.Equal(t, object.Size(), deserializedObject.Size())

			if tt.type_ == ObjectTypeBlob {
				blob := object.(*Blob)
				deserializedBlob := deserializedObject.(*Blob)
				assert.Equal(t, blob.Body(), deserializedBlob.Body())
			} else {
				tree := object.(*Tree)
				deserializedTree := deserializedObject.(*Tree)
				assert.Equal(t, tree.Body(), deserializedTree.Body())
			}
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
	serializedBlob := blob.Serialize()

	expectedData := "blob 5\x00hello"
	assert.Equal(t, expectedData, string(serializedBlob))
}
