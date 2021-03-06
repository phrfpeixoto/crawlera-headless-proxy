package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type data struct {
	Time, CA, Key string
}

var fileTemplate = `package main

// nolint: goimports, gochecknoglobals
// autogenerated by on {{ .Time }}

// DefaultCertCA is hardcoded TLS CA certificate
// nolint: gochecknoglobals
var DefaultCertCA = []byte(` + "`{{ .CA }}`" + `)

// DefaultPrivateKey is hardcoded TLS private key
// nolint: gochecknoglobals
var DefaultPrivateKey = []byte(` + "`{{ .Key }}`" + `)
`

func main() {
	caCertFile, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatal(err)
	}

	caKeyFile, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	caKey, err := ioutil.ReadFile(caKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := filepath.Abs(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	tpl := template.Must(template.New("tpl").Parse(fileTemplate))
	tpl.Execute(fp, data{ // nolint: errcheck
		Time: time.Now().Format(time.RFC1123Z),
		CA:   string(bytes.TrimSpace(caCert)),
		Key:  string(bytes.TrimSpace(caKey)),
	})
}
