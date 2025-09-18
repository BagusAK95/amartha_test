package html

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"
)

type HtmlTemplate struct {
	template *template.Template
}

func NewTemplate() (*HtmlTemplate, error) {
	tmpl, err := template.New("baseTemplate").Funcs(template.FuncMap{
		"FormatNumber":   FormatNumber,
		"FormatCurrency": FormatCurrency,
		"FormatDate":     FormatDate,
	}).ParseGlob("./templates/**/*.html")
	if err != nil {
		return nil, err
	}

	return &HtmlTemplate{
		template: tmpl,
	}, nil
}

func (tmpl *HtmlTemplate) Execute(wr io.Writer, file string, data any) (err error) {
	err = tmpl.template.ExecuteTemplate(wr, file, data)
	if err != nil {
		return
	}

	return
}

func FormatNumber(amount float64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	rounded := fmt.Sprintf("%.0f", amount)

	var formatted []string
	for i := len(rounded); i > 0; i -= 3 {
		if i-3 > 0 {
			formatted = append([]string{rounded[i-3 : i]}, formatted...)
		} else {
			formatted = append([]string{rounded[:i]}, formatted...)
		}
	}

	return sign + strings.Join(formatted, ".")
}

func FormatCurrency(amount float64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	rounded := fmt.Sprintf("%.0f", amount)

	var formatted []string
	for i := len(rounded); i > 0; i -= 3 {
		if i-3 > 0 {
			formatted = append([]string{rounded[i-3 : i]}, formatted...)
		} else {
			formatted = append([]string{rounded[:i]}, formatted...)
		}
	}

	return sign + "Rp" + strings.Join(formatted, ".")
}

func FormatDate(date time.Time) string {
	return date.Format("02 January 2006")
}
