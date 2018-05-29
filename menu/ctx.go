package menu

type Ctx struct {
	Meta map[string]interface{}
	Menu Menu
	err  error
}

func (ctx *Ctx) Error(err error) {
	ctx.err = err
}
