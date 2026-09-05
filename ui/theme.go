package ui

import (
	"gioui.org/widget/material"
	"github.com/chapar-rest/chapar/ui/chapartheme"
)

var Theme *chapartheme.Theme

func SetupTheme() {
	Theme = chapartheme.New(material.NewTheme())
}
