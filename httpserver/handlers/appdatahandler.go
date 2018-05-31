package handlers

import "github.com/devfeel/dotweb"

func AppData(ctx dotweb.Context) error {
	return userlogbase(ctx, "AppData")
}
