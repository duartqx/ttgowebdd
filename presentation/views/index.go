package views

import (
	"html/template"
	"io"
)

const (
	resultsPartial string = "task_results"
	rowPartial     string = "task_row"
	indexPartial   string = "index"
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
	return iv.tmplEngine.ExecuteTemplate(w, indexPartial, data)
}

func (iv IndexView) ExecuteRow(w io.Writer, data interface{}) error {
	return iv.tmplEngine.ExecuteTemplate(w, rowPartial, data)
}

func (iv IndexView) ExecuteResults(w io.Writer, data interface{}) error {
	return iv.tmplEngine.ExecuteTemplate(w, resultsPartial, data)
}
