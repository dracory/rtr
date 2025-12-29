package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/dracory/rtr"
)

//go:embed static
var staticFS embed.FS

func main() {
	router := rtr.NewRouter()

	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}

	router.AddRoute(rtr.GetStaticFS("/static/*", sub))

	router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
  <title>Static Embed Example</title>
</head>
<body>
  <h1>Embedded Static Files</h1>
  <p>These files are embedded into the binary using Go's embed.FS.</p>
  <ul>
    <li><a href="/static/style.css">/static/style.css</a></li>
    <li><a href="/static/script.js">/static/script.js</a></li>
    <li><a href="/static/data.json">/static/data.json</a></li>
  </ul>
</body>
</html>`)
	}))

	port := ":8080"
	fmt.Printf("Embedded static server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
