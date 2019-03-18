/*
Vanityd serves Go packages via vanity URLs. The default settings defined in
the "config" varaible may be overwitten by an /etc/vanityd.conf file in
the following format:

	# Lines starting with a pound sign are comments.
	vanityUrl=example.com
	repoUrl=https://github.com/username
	port=1234

See "go held importpath" for more information on custom import paths.
*/
package main // import "go.linskey.org/vanityd"

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const configFile = "/etc/vanityd.conf"

var config = map[string]string{
	"vanityUrl": "example.com",
	"repoUrl":   "https://github.com/username",
	"protocol":  "git",
	"port":      "8000",
}

const tpl = `
<html>
<head>
<meta name="go-import" content="{{.vanityUrl}}/{{.path}} {{.protocol}} {{.repoUrl}}/{{.path}}">
</head>
</html>
`

var html = template.Must(template.New("output").Parse(tpl))

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	addr := "localhost:" + config["port"]
	log.Fatal(http.ListenAndServe(addr, nil))
}

func initConfig() error {
	file, err := os.Open(configFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return fmt.Errorf(
				"%s: Invalid syntax on line %d: %q\n", configFile, i, line)
		}
		config[parts[0]] = parts[1]
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	isGoGet := r.URL.Query().Get("go-get") == "1"
	if !isGoGet {
		repo := config["repoUrl"] + "/" + path
		http.Redirect(w, r, repo, 302)
		return
	}

	config["path"] = path
	if err := html.Execute(w, config); err != nil {
		log.Fatal(err)
	}
}
