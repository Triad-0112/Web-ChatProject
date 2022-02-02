package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"
)

//template struct
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//Writing response with file and stream it on HTTP
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Go web server is do concurrent, regardless go routines are calling ServeHTTP
	t.once.Do(func() {
		// Joining folder with slash filepath for ex in my case templates/filename
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	//Execute is writing w which was the response in byte.
	t.templ.Execute(w, data)
}
