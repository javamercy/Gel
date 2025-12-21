package domain

import (
	"Gel/core/constant"
	"bytes"
	"errors"
)

var (
	ErrInvalidCommitField = errors.New("invalid commit field")
)

const (
	CommitFieldTree      string = "tree"
	CommitFieldParent    string = "parent"
	CommitFieldAuthor    string = "author"
	CommitFieldCommitter string = "committer"
	CommitFieldMessage   string = "message"
)

type Identity struct {
	Name      string
	Email     string
	Timestamp string
	Timezone  string
}

func (identity *Identity) serialize() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(identity.Name)
	buffer.WriteString(constant.SpaceStr)
	buffer.WriteString(identity.Email)
	buffer.WriteString(constant.SpaceStr)
	buffer.WriteString(identity.Timestamp)
	buffer.WriteString(constant.SpaceStr)
	buffer.WriteString(identity.Timezone)

	return buffer.Bytes()
}

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
		fields: CommitFields{
			TreeHash:     "",
			ParentHashes: nil,
			Author:       Identity{},
			Committer:    Identity{},
			Message:      "",
		},
	}
}

func (commit *Commit) SerializeBody() []byte {
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
	buffer.WriteString(commit.fields.Message)
	return buffer.Bytes()
}

func DeserializeCommit(data []byte) (*Commit, error) {
	// TODO: implement deserialization logic
	return nil, nil
}

func isValidField(field string) error {
	switch field {
	case CommitFieldTree,
		CommitFieldParent,
		CommitFieldAuthor,
		CommitFieldCommitter,
		CommitFieldMessage:
		return nil
	}
	return ErrInvalidCommitField
}

func deserializeIdentity(data []byte, start int) (Identity, int, error) {
	i := start
	for i < len(data) && data[i] != constant.NewLineByte {
		i++
	}
	if i >= len(data) {
		return Identity{}, i, errors.New("invalid identity format")
	}

	fieldStr := string(data[start:i])
	if err := isValidField(fieldStr); err != nil {
		return Identity{}, i, err
	}
	i++

	deserializeIdentityField := func(delimiter byte) (string, error) {
		if i >= len(data) {
			return "", errors.New("invalid identity format")
		}
		partStart := i
		for i < len(data) && data[i] != delimiter {
			i++
		}
		if i >= len(data) {
			return "", errors.New("invalid identity format")
		}
		out := string(data[partStart:i])
		i++
		return out, nil
	}

	name, err := deserializeIdentityField(constant.SpaceByte)
	if err != nil {
		return Identity{}, i, err
	}

	email, err := deserializeIdentityField(constant.SpaceByte)
	if err != nil {
		return Identity{}, i, err
	}

	timestamp, err := deserializeIdentityField(constant.SpaceByte)
	if err != nil {
		return Identity{}, i, err
	}

	timezone, err := deserializeIdentityField(constant.NewLineByte)
	if err != nil {
		return Identity{}, i, err
	}

	return Identity{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
		Timezone:  timezone,
	}, i, nil
}

func deserializeMessage(data []byte, start int) string {
	message := string(data[start:])
	return message
}
