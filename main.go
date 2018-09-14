package main

import (
	"gopkg.in/ini.v1"
	"github.com/valyala/fasthttp"
	"github.com/onebone/onessentials-go"
	"log"
	"net/url"
	"fmt"
	"path"
)

func main() {
	onessentials.CopyFile(path.Join("resources", "redir.ini"), "redir.ini")

	cfg, err := ini.Load("redir.ini")
	if err != nil {
		log.Panic(err)
	}

	for _, k := range cfg.Section("").Keys() {
		if _, err := url.ParseRequestURI(k.Value()); err != nil {
			log.Panic(fmt.Sprintf("Invalid URI given: %s=%s", k.String(), k.Value()))
		}
	}

	type Config struct {
		Port	int `json:"port"`
		Notify	bool `json:"notify-bot"`
	}
	var config Config
	err = onessentials.InitConfig(&config)
	if err != nil {
		log.Panic(err)
	}

	fasthttp.ListenAndServe(fmt.Sprintf(":%d", config.Port), func(ctx *fasthttp.RequestCtx) {
		log.Printf("Got access to: %s", string(ctx.Request.URI().Path()))

		key, err := cfg.Section("").GetKey(string(ctx.Request.URI().Path()))
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Response.BodyWriter().Write([]byte("There is nothing to show here."))

			return
		}

		ctx.Redirect(string(key.Value()), fasthttp.StatusTemporaryRedirect)
	})
}