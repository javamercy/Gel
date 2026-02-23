package domain

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrInvalidCommitFormat = errors.New("invalid commit format")
	ErrInvalidCommitField  = errors.New("invalid commit field")
)

const (
	CommitFieldTree      string = "tree"
	CommitFieldParent    string = "parent"
	CommitFieldAuthor    string = "author"
	CommitFieldCommitter string = "committer"
)

type CommitFields struct {
	TreeHash     string
	ParentHashes []string
	Author       Identity
	Committer    Identity
	Message      string
}

type Commit struct {
	body []byte
	CommitFields
}

func (commit *Commit) Body() []byte {
	return commit.body
}

func (commit *Commit) Type() ObjectType {
	return ObjectTypeCommit
}

func (commit *Commit) Size() int {
	return len(commit.body)
}

func (commit *Commit) Serialize() []byte {
	return SerializeObject(ObjectTypeCommit, commit.body)
}

func NewCommit(body []byte) (*Commit, error) {
	fields, err := DeserializeCommit(body)
	if err != nil {
		return nil, err
	}
	return fields, nil
}
func NewCommitFromFields(commitFields CommitFields) *Commit {
	return &Commit{
		body:         SerializeBody(commitFields),
		CommitFields: commitFields,
	}
}

func SerializeBody(fields CommitFields) []byte {
	// SerializeBody assumes the commit fields have been validated by the caller.

	var buffer bytes.Buffer
	buffer.WriteString(CommitFieldTree)
	buffer.WriteString(" ")
	buffer.WriteString(fields.TreeHash)
	buffer.WriteString("\n")

	for _, parentHash := range fields.ParentHashes {
		buffer.WriteString(CommitFieldParent)
		buffer.WriteString(" ")
		buffer.WriteString(parentHash)
		buffer.WriteString("\n")
	}
	buffer.WriteString(CommitFieldAuthor)
	buffer.WriteString(" ")
	buffer.Write(fields.Author.serialize())
	buffer.WriteString("\n")
	buffer.WriteString(CommitFieldCommitter)
	buffer.WriteString(" ")
	buffer.Write(fields.Committer.serialize())
	buffer.WriteString("\n")
	buffer.WriteString("\n")
	buffer.WriteString(fields.Message)

	return buffer.Bytes()
}

func DeserializeCommit(data []byte) (*Commit, error) {
	i := 0
	var fields CommitFields
	hasTree := false
	hasAuthor := false
	hasCommitter := false
	hasMessage := false
	for i < len(data) {
		if data[i] == '\n' {
			i++
			fields.Message = string(data[i:])
			hasMessage = true
			break
		}
		fieldStr, nextI, err := deserializeFieldStr(data, i)
		if err != nil {
			return nil, err
		}
		i = nextI
		switch fieldStr {
		case CommitFieldTree:
			hexHash, nextI, err := deserializeTreeOrParent(data, i)
			if err != nil {
				return nil, err
			}
			fields.TreeHash = hexHash
			hasTree = true
			i = nextI
		case CommitFieldParent:
			hexHash, nextI, err := deserializeTreeOrParent(data, i)
			if err != nil {
				return nil, err
			}
			fields.ParentHashes = append(fields.ParentHashes, hexHash)
			i = nextI
		case CommitFieldAuthor:
			author, nextI, err := deserializeIdentity(data, i)
			if err != nil {
				return nil, err
			}
			fields.Author = author
			hasAuthor = true
			i = nextI
		case CommitFieldCommitter:
			committer, nextI, err := deserializeIdentity(data, i)
			if err != nil {
				return nil, err
			}
			fields.Committer = committer
			hasCommitter = true
			i = nextI
		}
	}
	if !hasTree || !hasAuthor || !hasCommitter || !hasMessage {
		return nil, ErrInvalidCommitFormat
	}
	return NewCommitFromFields(fields), nil
}

func deserializeFieldStr(data []byte, start int) (string, int, error) {
	i := start
	for i < len(data) && data[i] != ' ' {
		i++
	}

	if i >= len(data) {
		return "", i, ErrInvalidCommitFormat
	}

	fieldStr := string(data[start:i])
	if ok := isValidCommitField(fieldStr); !ok {
		return "", i, ErrInvalidCommitField
	}
	return fieldStr, i + 1, nil
}

func deserializeTreeOrParent(data []byte, start int) (string, int, error) {
	i := start
	for i < len(data) && data[i] != '\n' {
		i++
	}
	if i >= len(data) {
		return "", i, ErrInvalidCommitFormat
	}
	hexHash := string(data[start:i])

	// TODO: refactor here to use a common validator
	if len(hexHash) != SHA256HexLength {
		return "", i, fmt.Errorf("invalid hash length: got %d, expected %d", len(hexHash), SHA256HexLength)
	}
	for _, c := range hexHash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return "", i, fmt.Errorf("invalid hash character: '%c'", c)
		}
	}

	return hexHash, i + 1, nil
}

func deserializeIdentity(data []byte, start int) (Identity, int, error) {
	i := start
	lineEnd := i
	for lineEnd < len(data) && data[lineEnd] != '\n' {
		lineEnd++
	}
	if lineEnd >= len(data) {
		return Identity{}, i, ErrInvalidCommitFormat
	}

	emailStart := i
	for emailStart < lineEnd && data[emailStart] != '<' {
		emailStart++
	}
	if emailStart >= lineEnd {
		return Identity{}, i, ErrInvalidCommitFormat
	}

	nameBytes := bytes.TrimSpace(data[i:emailStart])
	name := string(nameBytes)

	emailEnd := emailStart + 1
	for emailEnd < lineEnd && data[emailEnd] != '>' {
		emailEnd++
	}
	if emailEnd >= lineEnd {
		return Identity{}, i, ErrInvalidCommitFormat
	}

	email := string(data[emailStart+1 : emailEnd])

	i = emailEnd + 1
	for i < lineEnd && data[i] == ' ' {
		i++
	}

	timestampStart := i
	for i < lineEnd && data[i] != ' ' {
		i++
	}
	if i >= lineEnd {
		return Identity{}, i, ErrInvalidCommitFormat
	}
	timestamp := string(data[timestampStart:i])

	i++
	if i >= lineEnd {
		return Identity{}, i, ErrInvalidCommitFormat
	}
	timezone := string(data[i:lineEnd])

	return Identity{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
		Timezone:  timezone,
	}, lineEnd + 1, nil
}

func isValidCommitField(field string) bool {
	switch field {
	case CommitFieldTree,
		CommitFieldParent,
		CommitFieldAuthor,
		CommitFieldCommitter:
		return true
	}
	return false
}
