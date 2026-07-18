// Package migrations embeds the SQL migration files so the migrate binary is
// fully self-contained: it needs no migration files on disk at runtime, which
// makes it safe to run as a Kubernetes init container or a Render pre-deploy
// command.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
