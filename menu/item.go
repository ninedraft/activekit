package menu

type MenuItem struct {
	Label  string
	Action func(Ctx)
}

func (item *MenuItem) String() string {
	return item.Label
}
