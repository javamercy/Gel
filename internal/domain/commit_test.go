package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommitFromFields_Valid(t *testing.T) {
	author := NewIdentity("Author", "author@example.com", "1234567890", "+0000")
	committer := NewIdentity("Committer", "committer@example.com", "1234567890", "+0000")
	fields := CommitFields{
		TreeHash:     "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
		ParentHashes: []string{},
		Author:       author,
		Committer:    committer,
		Message:      "Initial commit",
	}
	commit := NewCommitFromFields(fields)
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
	author := NewIdentity("Author", "author@example.com", "123", "+0000")
	committer := NewIdentity("Committer", "committer@example.com", "123", "+0000")
	fields := CommitFields{
		TreeHash:     "hash",
		ParentHashes: []string{"parent1", "parent2"},
		Author:       author,
		Committer:    committer,
		Message:      "msg",
	}

	body := SerializeBody(fields)
	s := string(body)

	assert.Contains(t, s, "tree hash")
	assert.Contains(t, s, "parent parent1")
	assert.Contains(t, s, "parent parent2")
	assert.Contains(t, s, "author Author <author@example.com> 123 +0000")
	assert.Contains(t, s, "committer Committer <committer@example.com> 123 +0000")
	assert.Contains(t, s, "\n\nmsg")
}

func TestDeserializeCommit_Complex(t *testing.T) {
	body := "tree a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\n" +
		"parent b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\n" +
		"parent c1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\n" +
		"author A <a@e.c> 1 +0000\n" +
		"committer C <c@e.c> 2 +0000\n" +
		"\n" +
		"Multi-line\nMessage"

	commit, err := DeserializeCommit([]byte(body))
	require.NoError(t, err)

	assert.Equal(t, "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2", commit.TreeHash)
	assert.Len(t, commit.ParentHashes, 2)
	assert.Equal(t, "b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2", commit.ParentHashes[0])
	assert.Equal(t, "A", commit.Author.Name)
	assert.Equal(t, "Multi-line\nMessage", commit.Message)
}
