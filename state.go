package cantainer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
)

type state struct {
	db     *sql.DB
	app    *app.App
	socket Socket
}

func NewState(socket Socket, dir string, seeds []string) (*state, error) {

	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	options := []app.Option{
		app.WithAddress(socket.ExtendedAddress()),
		app.WithLogFunc(func(l client.LogLevel, format string, a ...interface{}) {
			slog.Info(fmt.Sprintf("%s: %s\n", l.String(), format), a...)
		}),
		app.WithNetworkLatency(10 * time.Millisecond),
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
		app:    app,
		socket: socket,
	}, nil
}

func (s *state) RegisterNode() error {
	_, err := s.db.Exec("INSERT INTO nodes VALUES (?, ?);", s.socket.Address, s.socket.Port)

	return err
}

func (s *state) Shutdown(ctx context.Context) error {

	slog.Info("running state shutdown")
	err := s.app.Handover(ctx)

	closeErr := s.app.Close()
	if closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	return err
}

func (s *state) UnregisterNode() error {
	_, err := s.db.Exec("DELETE FROM nodes WHERE address=?", s.socket.Address)

	return err
}

func (s *state) PeerAddresses() (map[string]bool, error) {

	rows, err := s.db.Query("SELECT address from nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]bool)
	for rows.Next() {
		var addr string
		if err := rows.Scan(&addr); err != nil {
			return nil, err
		}
		results[addr] = true
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// remove node's own ip address from the peers list
	delete(results, s.socket.Address)

	return results, nil
}
