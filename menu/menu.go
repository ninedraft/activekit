package menu

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"strings"
)

type Menu struct {
	Title        string
	Promt        string
	History      []string
	Items        MenuItems
	Meta         map[string]interface{}
	QueryHandler func(*Ctx)
	once         sync.Once
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
		if menu.Meta == nil {
			menu.Meta = map[string]interface{}{}
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
		if menu.History != nil {
			menu.History = append(menu.History, input)
		}
		var itemIndex int
		if ind, ok := optionSet[input]; ok {
			itemIndex = ind
		}
		if _, err = fmt.Sscan(input, &itemIndex); err == nil {
			itemIndex-- // decrement is very important, do not change!
		}
		if itemIndex >= 0 && itemIndex < len(menu.Items) {
			item, ctx := menu.RunItem(itemIndex)
			if ctx.err != nil || ctx.stop {
				return item, ctx.err
			}
			continue
		}
		if menu.QueryHandler == nil {
			fmt.Printf("Option %q not found\n", input)
		} else {
			var ctx = menu.Context().WithQuery(input)
			menu.QueryHandler(&ctx)
			if ctx.err != nil || ctx.stop {
				return nil, ctx.err
			}
		}
	}
}

func (menu *Menu) RunItem(i int) (*MenuItem, Ctx) {
	var item = menu.Items[i]
	var ctx = menu.Context()
	if item.Action != nil {
		item.Action(&ctx)
	}
	return item, ctx
}
