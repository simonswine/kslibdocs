package docgen

import (
	"strings"

	"github.com/gobuffalo/plush"

	rice "github.com/GeertJohan/go.rice"
)

//go:generate rice embed-go

type argumentType struct {
	propertyName string
	typeName     string
	typeDef      string
}

func (at *argumentType) ToDoc() (string, error) {
	box, err := rice.FindBox("templates")
	if err != nil {
		return "", err
	}

	tmpl, err := box.String("argument_type.md")
	if err != nil {
		return "", err
	}

	parts := strings.Split(at.typeDef, ".")
	typeURL := strings.Join(parts, "/")

	ctx := plush.NewContext()
	ctx.Set("propertyName", at.propertyName)
	ctx.Set("typeName", at.typeName)
	ctx.Set("typeDef", at.typeDef)
	ctx.Set("typeURL", typeURL)

	return plush.Render(tmpl, ctx)
}
