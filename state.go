package cantainer

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
)

type state struct {
	db     *sql.DB
	socket Socket
}

func NewState(socket Socket, dir string, seeds []string) (*state, error) {

	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	options := []app.Option{
		app.WithAddress(socket.ExtendedAddress()),
		app.WithLogFunc(func(l client.LogLevel, format string, a ...interface{}) {
			slog.Debug(fmt.Sprintf("%s: %s\n", l.String(), format), a...)
		}),
		app.WithNetworkLatency(50 * time.Millisecond),
	}

	var remoteSeeds []string
	for _, seed := range seeds {
		if seed != socket.ExtendedAddress() {
			remoteSeeds = append(remoteSeeds, seed)
		}
	}

	if len(remoteSeeds) != 0 {
		options = append(options, app.WithCluster(remoteSeeds))
	}

	app, err := app.New(dir, options...)
	if err != nil {
		return nil, err
	}

	if err := app.Ready(context.Background()); err != nil {
		return nil, err
	}
	slog.Info("database connection is ready")

	db, err := app.Open(context.Background(), "state")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS nodes (address TEXT, port BIGINT, UNIQUE(address, port));")
	if err != nil {
		return nil, err
	}

	return &state{
		db:     db,
		socket: socket,
	}, nil
}

func (s *state) RegisterNode() error {
	_, err := s.db.Exec("INSERT INTO nodes VALUES (?, ?);", s.socket.Address, s.socket.Port)

	return err
}
