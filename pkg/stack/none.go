package stack

import "github.com/buildtool/scaffold/pkg/templating"

type None struct{}

func (n *None) Scaffold(dir string, data templating.TemplateData) error {
	return nil
}

func (n *None) Name() string {
	return "none"
}

var _ Stack = &None{}
