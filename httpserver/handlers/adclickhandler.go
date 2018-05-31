package handlers

import "github.com/devfeel/dotweb"

func ADClick(ctx dotweb.Context) error {
	return adbase(ctx, "ADClick")
}
