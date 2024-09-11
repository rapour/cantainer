package cantainer

import (
	"context"
	"database/sql"
	"os"

	"github.com/canonical/go-dqlite/app"
)

type state struct {
	db *sql.DB
}

func NewState(address string, dir string, seeds []string) (*state, error) {

	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	options := []app.Option{
		app.WithAddress(address),
		app.WithCluster(seeds),
	}

	app, err := app.New(dir, options...)
	if err != nil {
		return nil, err
	}

	if err := app.Ready(context.Background()); err != nil {
		return nil, err
	}

	db, err := app.Open(context.Background(), "state")
	if err != nil {
		return nil, err
	}

	return &state{
		db: db,
	}, nil
}
