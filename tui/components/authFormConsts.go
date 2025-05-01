package components

import "fmt"


type AuthMode int

// it is very important to keep the order of (SignUp = 0) because the authForm depends on it
const (
	SignUp AuthMode = iota
	SignIn

	MaxMode
)

func (m AuthMode) String() string {
	switch m {
	case SignIn:
		return "Sign-In"
	case SignUp:
		return "Sign-Up"
	default:
		return fmt.Sprintf("AuthMode(%d)", m)
	}
}
