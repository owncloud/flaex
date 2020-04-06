package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/owncloud/flaex/pkg/parsers"
)

// TemplateVariables holds all variables for the rendering
type TemplateVariables struct {
	Commands parsers.ParsedCommands
	Options  parsers.ParsedOptions
}

func main() {

	var commandPath = flag.String("command-path", "pkg/command/", "Path to command file or directory")
	var flagsetPath = flag.String("flagset-path", "pkg/flagset/", "Path to flagset file or directory")
	var templatePath = flag.String("template", "templates/CONFIGURATION.tmpl", "Path to go-template template file")
	flag.Parse()
	finfocp, err := os.Stat(*commandPath)
	if err != nil {
		log.Fatal(err)
	}

	finfofs, err := os.Stat(*flagsetPath)
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

	templateVariables := TemplateVariables{}

	if !finfocp.IsDir() {
		templateVariables.Commands, err = parsers.ParseCommandFile(*commandPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		templateVariables.Commands, err = parsers.ParseCommandDir(*commandPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !finfofs.IsDir() {
		templateVariables.Options, err = parsers.ParseFlagsetFile(*flagsetPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		templateVariables.Options, err = parsers.ParseFlagsetDir(*flagsetPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tpl.Execute(os.Stdout, templateVariables); err != nil {
		log.Fatal(err)
	}
}
