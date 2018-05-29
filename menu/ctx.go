package menu

import "strings"

type Ctx struct {
	Menu  Menu
	Query string
	err   error
	stop  bool
}

func (ctx *Ctx) Error(err error) {
	ctx.err = err
}

func (ctx *Ctx) Stop() {
	ctx.stop = true
}

func (ctx Ctx) WithQuery(query string) Ctx {
	ctx.Query = query
	return ctx
}

func (ctx Ctx) Args() (head string, tail []string) {
	var tokens = strings.Fields(ctx.Query)
	if len(tokens) == 0 {
		return "", []string{}
	}
	return tokens[0], tokens[1:]
}
