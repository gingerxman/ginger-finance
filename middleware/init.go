package middleware

import "github.com/gingerxman/eel"

func init() {
	eel.RegisterMiddleware(&eel.JWTMiddleware{})
	eel.RegisterMiddleware(&LoginMiddleware{})
}
