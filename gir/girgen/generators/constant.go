package generators

import (
	"fmt"

	"github.com/go-gst/go-glib/gir/girgen/file"
	"github.com/go-gst/go-glib/gir/girgen/typesystem"
)

type ConstantGenerator struct {
	Doc SubGenerator

	*typesystem.Constant
}

func (g *ConstantGenerator) Generate(w *file.Package) {
	g.Doc.Generate(w.Go())

	fmt.Fprintf(w.Go(), "const %s = %s\n", g.GoIndentifier(), g.GoValue)
}

func NewConstantGenerator(constant *typesystem.Constant) *ConstantGenerator {
	return &ConstantGenerator{
		Doc:      NewGoDocGenerator(constant),
		Constant: constant,
	}
}
