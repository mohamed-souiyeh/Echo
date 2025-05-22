package messages

import db "echo/db/sqlc_generated"

type SignUpAttemptMsg struct {
	Username  string
	Passwords []string
}

type SignInAttemptMsg struct {
	Username string
	Password string
}

type AuthSuccessMsg struct {
	User db.User
}

type AuthFailedMsg struct {
	Reason      string
	DebugReason string
}

type AccessChatMsg struct {
	User db.User
}

type LogoutMsg struct {
	Username string
}
