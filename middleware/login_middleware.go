package middleware

import (
	"github.com/gingerxman/eel"
)

type LoginMiddleware struct {
	eel.Middleware
}

func (this *LoginMiddleware) ProcessRequest(ctx *eel.Context) {
	//log.Logger.Info("i am in login middleware process request")
}

