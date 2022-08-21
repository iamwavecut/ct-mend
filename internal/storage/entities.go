package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Client struct {
		_ primitive.ObjectID `bson:"_id"`

		ID       *int           `json:"id,omitempty" bson:"id" db:"id"`
		Name     string         `json:"name" bson:"name" db:"name"`
		Settings ClientSettings `json:"settings" bson:"settings"`
	}

	ClientSettings struct {
		CodeScanInterval time.Duration `json:"code_scan_interval" bson:"code_scan_interval" db:"code_scan_interval"`
	}

	Project struct {
		_ primitive.ObjectID `bson:"_id"`

		ID       *int   `json:"id,omitempty" bson:"id" db:"id"`
		ClientID *int   `json:"client_id,omitempty" bson:"client_id" db:"client_id"`
		Name     string `json:"name" bson:"name" db:"name"`
	}
)
