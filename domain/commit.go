package domain

import (
	"Gel/core/constant"
	"bytes"
	"errors"
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
	*BaseObject
	fields CommitFields
}

func NewCommit(data []byte) *Commit {
	return &Commit{
		BaseObject: &BaseObject{
			objectType: ObjectTypeCommit,
			data:       data,
		},
	}
}
func NewCommitFromFields(data []byte, fields CommitFields) *Commit {
	return &Commit{
		BaseObject: &BaseObject{
			objectType: ObjectTypeCommit,
			data:       data,
		},
		fields: fields,
	}
}

func (commit *Commit) SerializeBody() []byte {
	// SerializeBody assumes the commit fields have been validated by the caller.

	var buffer bytes.Buffer
	buffer.WriteString(CommitFieldTree)
	buffer.WriteByte(constant.SpaceByte)
	buffer.WriteString(commit.fields.TreeHash)
	buffer.WriteByte(constant.NewLineByte)

	for _, parentHash := range commit.fields.ParentHashes {
		buffer.WriteString(CommitFieldParent)
		buffer.WriteByte(constant.SpaceByte)
		buffer.WriteString(parentHash)
		buffer.WriteByte(constant.NewLineByte)
	}
	buffer.WriteString(CommitFieldAuthor)
	buffer.WriteByte(constant.SpaceByte)
	buffer.Write(commit.fields.Author.serialize())
	buffer.WriteByte(constant.NewLineByte)
	buffer.WriteString(CommitFieldCommitter)
	buffer.WriteByte(constant.SpaceByte)
	buffer.Write(commit.fields.Committer.serialize())
	buffer.WriteByte(constant.NewLineByte)
	buffer.WriteByte(constant.NewLineByte)
	buffer.WriteString(commit.fields.Message)
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
		if data[i] == constant.NewLineByte {
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
	return NewCommitFromFields(data, fields), nil
}

func deserializeFieldStr(data []byte, start int) (string, int, error) {
	i := start
	for i < len(data) && data[i] != constant.SpaceByte {
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
	for i < len(data) && data[i] != constant.NewLineByte {
		i++
	}
	if i >= len(data) {
		return "", i, ErrInvalidCommitFormat
	}
	hexHash := string(data[start:i])
	// TODO: validate hash
	return hexHash, i + 1, nil
}

func deserializeIdentity(data []byte, start int) (Identity, int, error) {
	i := start
	lineEnd := i
	for lineEnd < len(data) && data[lineEnd] != constant.NewLineByte {
		lineEnd++
	}
	if lineEnd >= len(data) {
		return Identity{}, i, ErrInvalidCommitFormat
	}

	name, i, err := parseDelimitedField(data, i, constant.SpaceByte)
	if err != nil {
		return Identity{}, i, err
	}

	email, i, err := parseDelimitedField(data, i, constant.SpaceByte)
	if err != nil {
		return Identity{}, i, err
	}

	timestamp, i, err := parseDelimitedField(data, i, constant.SpaceByte)
	if err != nil {
		return Identity{}, i, err
	}

	timezone, i, err := parseDelimitedField(data, i, constant.NewLineByte)
	if err != nil {
		return Identity{}, i, err
	}

	return Identity{name, email, timestamp, timezone}, i, nil
}

func parseDelimitedField(data []byte, start int, delimiter byte) (value string, nextPosition int, err error) {
	i := start
	if i >= len(data) {
		return "", i, ErrInvalidCommitFormat
	}

	fieldStart := i
	for i < len(data) && data[i] != delimiter {
		i++
	}
	if i >= len(data) {
		return "", i, ErrInvalidCommitFormat
	}

	return string(data[fieldStart:i]), i + 1, nil
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
