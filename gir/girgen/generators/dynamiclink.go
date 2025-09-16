package generators

import (
	"github.com/go-gst/go-glib/gir/girgen/typesystem"
)

// WithDynamicLinking builds all generators needed for generated cgo dynamic linked code.
func WithDynamicLinking(repositories []*typesystem.Repository) (gens []*NamespaceGenerator) {
	for _, ns := range repositories {
		for _, ns := range ns.Namespaces {
			gens = append(gens, NewNamespaceGenerator(ns))
		}
	}

	return gens
}
