package main

import (
	"log"
)

func main() {

	opts, err := ParseFile("/home/ilja/code/ocis/ocis-konnectd/pkg/flagset/flagset.go")
	if err != nil {
		log.Fatal(err)
	}

	md := markdown(opts)

	log.Printf("%#v", md)

}

func markdown(opts []ParsedOption) map[string][]ParsedOption {
	byFunc := make(map[string][]ParsedOption)
	for k, o := range opts {
		byFunc[o.fnName] = append(byFunc[o.fnName], opts[k])
	}

	return byFunc

}
