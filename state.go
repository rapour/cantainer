package cantainer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS containers (address TEXT UNIQUE, node TEXT);")
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

func (s *state) RegisterContainer(address *netip.Addr) error {

	_, err := s.db.Exec("INSERT INTO containers VALUES (?, ?);", address.String(), s.socket.ExtendedAddress())

	return err
}

func (s *state) UnregisterContainer(address *netip.Addr) error {

	_, err := s.db.Exec("DELETE FROM containers WHERE address=?", address.String())

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

// TODO: Get container addresses for a specific subnet using SQL
func (s *state) Containers() ([]netip.Addr, error) {

	rows, err := s.db.Query("SELECT address from containers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []netip.Addr{}
	for rows.Next() {
		var addrStr string
		if err := rows.Scan(&addrStr); err != nil {
			return nil, err
		}

		addr, err := netip.ParseAddr(addrStr)
		if err != nil {
			return nil, err
		}

		results = append(results, addr)
	}

	return results, rows.Err()
}
