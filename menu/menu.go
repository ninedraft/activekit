package menu

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"strings"
)

type Menu struct {
	Title               string
	Promt               string
	History             []string
	Items               MenuItems
	Meta                map[string]interface{}
	CustomOptionHandler func(string) error
	QueryHandler        func(Ctx)
	once                sync.Once
}

func (menu *Menu) init() {
	menu.once.Do(func() {
		if menu.History == nil {
			menu.History = make([]string, 0, 16)
		}
		if menu.Title == "" {
			menu.Title = "What's next?"
		}
		if menu.Promt == "" {
			menu.Promt = "Choose wisely: "
		}
		menu.Items = menu.Items.NotNil()
	})
}

func (menu *Menu) scanLine() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		if scanner.Text() == "" {
			continue
		}
		return scanner.Text(), nil
	}
	return "", nil
}

func (menu *Menu) Context() Ctx {
	return Ctx{
		Menu: *menu,
	}
}

func (menu *Menu) Run() (*MenuItem, error) {
	menu.init()
	optionSet := map[string]int{}
	for i, item := range menu.Items {
		optionSet[item.Label] = i
	}
	for {
		fmt.Printf("%s\n", menu.Title)
		for i, item := range menu.Items {
			fmt.Printf("%d) %s\n", i+1, item.String())
		}
		fmt.Printf("%s", menu.Promt)
		input, err := menu.scanLine()
		if err != nil {
			return nil, err
		}
		input = strings.TrimSpace(input)
		if ind, ok := optionSet[input]; ok {
			item := menu.Items[ind]
			if item.Action == nil {
				return item, nil
			}
			var ctx = menu.Context()
			item.Action(ctx)
			return item, ctx.err
		}
		ind := 0
		if _, err = fmt.Sscan(input, &ind); err == nil && (ind > 0 && ind <= len(menu.Items)) {
			item := menu.Items[ind-1] // -1 is very important, do not change!
			if item.Action == nil {
				return item, nil
			}
			var ctx = menu.Context()
			item.Action(ctx)
			return item, ctx.err
		}
		if menu.CustomOptionHandler == nil {
			fmt.Printf("Option %q not found\n", input)
		} else {
			return nil, menu.CustomOptionHandler(input)
		}
	}
}
