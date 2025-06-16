package authsql

import (
	_ "embed"
)

//go:embed update_auth_code_for_sending_validation_email.sql
var UpdateAuthCodeQuery string
