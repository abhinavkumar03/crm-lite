package database

import (
	"fmt"
	"net/url"
)

// MigrationURL builds a golang-migrate connection URL using the pgx v5 driver
// scheme ("pgx5"). Credentials are URL-encoded so passwords with special
// characters are handled correctly.
func MigrationURL(host, port, user, password, dbName, sslMode string) string {
	u := &url.URL{
		Scheme: "pgx5",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/" + dbName,
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}
