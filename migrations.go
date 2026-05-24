package test_em

import "embed"

//go:embed migrations/*.sql
var MigrationFS embed.FS
