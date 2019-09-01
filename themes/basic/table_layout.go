package basic

import (
	"github.com/jakemakesstuff/gxui"
	"github.com/jakemakesstuff/gxui/mixins"
)

func CreateTableLayout(theme *Theme) gxui.TableLayout {
	l := &mixins.TableLayout{}
	l.Init(l, theme)
	return l
}
