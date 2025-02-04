package stack

import "github.com/buildtool/scaffold/pkg/templating"

type Stack interface {
	Scaffold(dir string, data templating.TemplateData) error
	Name() string
}

var Stacks = map[string]Stack{
	"none":  &None{},
	"go":    &Go{},
	"scala": &Scala{},
}
