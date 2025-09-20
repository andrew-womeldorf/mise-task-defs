package users

import _ "embed"

//go:embed schema.sql
var Migrations string
