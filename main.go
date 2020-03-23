package main

import (
	"flag"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {

	var flagsetPath = flag.String("path", "", "Path to flagset file or directory")
	var templatePath = flag.String("template", "flaex.tpl", "Path to go-template template file")
	flag.Parse()
	fi, err := os.Stat(*flagsetPath)
	if err != nil {
		log.Fatal(err)
	}

	if !fi.IsDir() {
		tplContent, err := ioutil.ReadFile(*templatePath)
		if err != nil {
			log.Fatalf("unable to read template from %v: %v", *templatePath, err)
		}

		tpl, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(string(tplContent))
		if err != nil {
			log.Fatal(err)
		}

		opts, err := ParseFile(*flagsetPath)
		if err != nil {
			log.Fatal(err)
		}

		if err := tpl.Execute(os.Stdout, opts); err != nil {
			log.Fatal(err)
		}
	}
}
