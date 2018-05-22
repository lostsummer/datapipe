package handlers

import (
	"github.com/devfeel/dotweb"
)

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx.WriteString("welcome to datapipe | version=1")
	return nil
}
