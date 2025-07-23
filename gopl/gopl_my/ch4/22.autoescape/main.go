// Autoescape demonstrates automatic HTML escaping in html/template.
// See page 117.
package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	const templ = `<p>A: {{.A}}</p><p>B: {{.B}}</p>`
	t := template.Must(template.New("escape").Parse(templ))
	var data struct {
		A string        // untrusted plain text
		B template.HTML // trusted HTML
	}
	data.A = "<b>Hello!</b>" // <p>A: &lt;b&gt;Hello!&lt;/b&gt; 污点变量，HTML实体转义
	data.B = "<b>Hello!</b>" // </p><p>B: <b>Hello!</b></p> 用户不可控变量，安全，不转义
	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
