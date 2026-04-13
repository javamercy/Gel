package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTreeHashHex = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
)

func TestNewCommitFromFields_Valid(t *testing.T) {
	author, err := NewIdentity("Author", "author@example.com", "1234567890", "+0000")
	require.NoError(t, err)
	committer, err := NewIdentity("Committer", "committer@example.com", "1234567890", "+0000")
	require.NoError(t, err)

	treeHash, err := NewHash(testTreeHashHex)
	require.NoError(t, err)

	fields := CommitFields{
		TreeHash:     treeHash,
		ParentHashes: []Hash{},
		Author:       author,
		Committer:    committer,
		Message:      "Initial commit",
	}
	commit, err := NewCommitFromFields(fields)
	require.NoError(t, err)
	assert.NotNil(t, commit)
	assert.NotEmpty(t, commit.Body())
}

func TestDeserializeCommit_Malformed(t *testing.T) {
	data := []byte("author Name <email> 123 +0000\ncommitter Name <email> 123 +0000\n\nMessage")
	_, err := DeserializeCommit(data)
	assert.ErrorIs(t, err, ErrInvalidCommitFormat)

	data = []byte("tree a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\ncommitter Name <email> 123 +0000\n\nMessage")
	_, err = DeserializeCommit(data)
	assert.ErrorIs(t, err, ErrInvalidCommitFormat)

	data = []byte("invalid a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\n\nMessage")
	_, err = DeserializeCommit(data)
	assert.ErrorIs(t, err, ErrInvalidCommitField)
}

func TestSerializeBody_Content(t *testing.T) {
	author, err := NewIdentity("Author", "author@example.com", "123", "+0000")
	require.NoError(t, err)
	committer, err := NewIdentity("Committer", "committer@example.com", "123", "+0000")
	require.NoError(t, err)

	treeHash, err := NewHash("0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	parentHash1, err := NewHash("1111111111111111111111111111111111111111111111111111111111111111")
	require.NoError(t, err)
	parentHash2, err := NewHash("2222222222222222222222222222222222222222222222222222222222222222")
	require.NoError(t, err)

	fields := CommitFields{
		TreeHash:     treeHash,
		ParentHashes: []Hash{parentHash1, parentHash2},
		Author:       author,
		Committer:    committer,
		Message:      "msg",
	}

	body := serializeBody(fields)
	s := string(body)

	assert.Contains(t, s, "tree 0000")
	assert.Contains(t, s, "parent 1111")
	assert.Contains(t, s, "parent 2222")
	assert.Contains(t, s, "author Author <author@example.com> 123 +0000")
	assert.Contains(t, s, "committer Committer <committer@example.com> 123 +0000")
	assert.Contains(t, s, "\n\nmsg")
}

func TestDeserializeCommit_Complex(t *testing.T) {
	body := "tree " + testTreeHashHex + "\n" +
		"parent b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\n" +
		"parent c1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\n" +
		"author A <a@e.c> 1 +0000\n" +
		"committer C <c@e.c> 2 +0000\n" +
		"\n" +
		"Multi-line\nMessage"

	commit, err := DeserializeCommit([]byte(body))
	require.NoError(t, err)

	assert.Equal(t, testTreeHashHex, commit.TreeHash.String())
	assert.Len(t, commit.ParentHashes, 2)
	assert.Equal(t, "b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2", commit.ParentHashes[0].String())
	assert.Equal(t, "A", commit.Author.Name)
	assert.Equal(t, "Multi-line\nMessage", commit.Message)
}

func TestNewCommitFromFields_Invalid(t *testing.T) {
	validID, err := NewIdentity("Author", "author@example.com", "1", "+0000")
	require.NoError(t, err)
	treeHash, err := NewHash(testTreeHashHex)
	require.NoError(t, err)

	tests := []struct {
		name   string
		fields CommitFields
	}{
		{
			name: "empty tree hash",
			fields: CommitFields{
				Author:    validID,
				Committer: validID,
				Message:   "m",
			},
		},
		{
			name: "invalid author",
			fields: CommitFields{
				TreeHash:     treeHash,
				Author:       Identity{},
				Committer:    validID,
				Message:      "m",
				ParentHashes: nil,
			},
		},
		{
			name: "empty parent hash",
			fields: CommitFields{
				TreeHash:     treeHash,
				ParentHashes: []Hash{{}},
				Author:       validID,
				Committer:    validID,
				Message:      "m",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := NewCommitFromFields(tt.fields)
				assert.Error(t, err)
			},
		)
	}
}

func TestDeserializeCommit_DuplicateFields(t *testing.T) {
	commitWithDuplicateTree := []byte(
		"tree " + testTreeHashHex + "\n" +
			"tree " + testTreeHashHex + "\n" +
			"author A <a@e.c> 1 +0000\n" +
			"committer C <c@e.c> 2 +0000\n\nmsg",
	)

	_, err := DeserializeCommit(commitWithDuplicateTree)
	assert.ErrorIs(t, err, ErrInvalidCommitFormat)
}

func TestCommit_BodyAndFieldsAreDefensiveCopies(t *testing.T) {
	author, err := NewIdentity("Author", "author@example.com", "1", "+0000")
	require.NoError(t, err)
	treeHash, err := NewHash(testTreeHashHex)
	require.NoError(t, err)
	parentHash, err := NewHash("b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	require.NoError(t, err)

	fields := CommitFields{
		TreeHash:     treeHash,
		ParentHashes: []Hash{parentHash},
		Author:       author,
		Committer:    author,
		Message:      "hello",
	}

	commit, err := NewCommitFromFields(fields)
	require.NoError(t, err)

	fields.ParentHashes[0] = Hash{}
	assert.Equal(t, parentHash, commit.ParentHashes[0])

	body := commit.Body()
	body[0] = 'x'
	assert.NotEqual(t, body[0], commit.Body()[0])
}
