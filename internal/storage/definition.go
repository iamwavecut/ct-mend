// Package storage Adapters, entities and definition of storage
package storage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/iamwavecut/ct-mend/internal/config"
	"github.com/iamwavecut/ct-mend/tools"
)

const (
	TypeSQLite Type = iota + 1
	TypeMongoDB
)

type (
	Type int

	Adapter interface {
		Init(ctx context.Context, connAddr string) error

		SelectClients() ([]*Client, error)
		GetClient(id int) (*Client, error)
		UpsertClient(client *Client) (*Client, error)
		DeleteClient(id int) error

		SelectProjects() ([]*Project, error)
		GetProject(id int) (*Project, error)
		SelectProjectsOfClient(id int) ([]*Project, error)
		UpsertProject(project *Project) (*Project, error)
		DeleteProject(id int) error
	}

	ErrNotFound struct {
		cause error
	}
	ErrNilEntity struct {
		cause error
	}
)

func (e ErrNotFound) Error() string {
	msg := "not found"
	if e.cause != nil {
		msg += ": " + e.cause.Error()
	}
	return msg
}

func (e ErrNilEntity) Error() string {
	msg := "nil entity passed"
	if e.cause != nil {
		msg += ": " + e.cause.Error()
	}
	return msg
}

func New(ctx context.Context, cfg *config.Storage) (Adapter, error) {
	if cfg == nil {
		return nil, errors.New("no storage config provided")
	}
	sType, err := typeFromString(cfg.Type)
	if !tools.Try(err) {
		return nil, err
	}
	var adapter Adapter
	switch {
	case sType == TypeSQLite:
		adapter = &SQLite{}
	case sType == TypeMongoDB:
		adapter = &MongoDB{}
	default:
		return nil, errors.New("unknown storage type " + cfg.Type)
	}
	err = adapter.Init(ctx, cfg.Addr)
	if !tools.Try(err) {
		return nil, err
	}
	return adapter, nil
}

func typeFromString(toFind string) (Type, error) {
	if foundType, ok := map[string]Type{
		"sqlite":  TypeSQLite,
		"mongodb": TypeMongoDB,
	}[toFind]; ok {
		return foundType, nil
	}
	return 0, errors.New("unrecognized type")
}

func passNotFound(err error) error {
	switch err {
	case sql.ErrNoRows:
		return ErrNotFound{err}
	case mongo.ErrNoDocuments:
		return ErrNotFound{err}
	}
	return err
}
