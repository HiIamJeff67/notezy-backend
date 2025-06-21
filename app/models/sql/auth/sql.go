package authsql

import (
	_ "embed"
)

//go:embed update_auth_code_by_email.sql
var UpdateAuthCodeQuery string

//go:embed reset_email.sql
var ResetEmailQuery string

//go:embed validate_email_by_auth_code.sql
var ValidateEmailQuery string

//go:embed delete_me.sql
var DeleteMeQuery string
