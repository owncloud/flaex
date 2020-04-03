package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
)

func main() {

	var flagsetPath = flag.String("path", "pkg/flagset/", "Path to flagset file or directory")
	var templatePath = flag.String("template", "templates/CONFIGURATION.tmpl", "Path to go-template template file")
	flag.Parse()
	fi, err := os.Stat(*flagsetPath)
	if err != nil {
		log.Fatal(err)
	}

	tplContent, err := ioutil.ReadFile(*templatePath)
	if err != nil {
		log.Fatalf("unable to read template from %v: %v", *templatePath, err)
	}

	tpl := template.Must(
		template.New("").Funcs(sprig.GenericFuncMap()).Parse(string(tplContent)),
	)

	var opts ParsedOptions

	if !fi.IsDir() {
		opts, err = ParseFile(*flagsetPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		opts, err = ParseDir(*flagsetPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tpl.Execute(os.Stdout, opts); err != nil {
		log.Fatal(err)
	}
}
