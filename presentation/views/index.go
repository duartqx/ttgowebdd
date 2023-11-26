package views

import (
	"html/template"
	"io"
)

const (
	rowPartial string = "task_row"
	index      string = "index"
)

type IndexView struct {
	tmplEngine *template.Template
}

func NewIndexView(tmplEngine *template.Template) *IndexView {
	return &IndexView{
		tmplEngine: tmplEngine,
	}
}

func (iv IndexView) Execute(w io.Writer, data interface{}) error {
	return iv.tmplEngine.ExecuteTemplate(w, index, data)
}

func (iv IndexView) ExecuteRow(w io.Writer, data interface{}) error {
	return iv.tmplEngine.ExecuteTemplate(w, rowPartial, data)
}
