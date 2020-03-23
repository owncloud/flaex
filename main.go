package main

import (
	"log"
)

func main() {

	opts, err := ParseFile("/home/ilja/code/ocis/ocis-konnectd/pkg/flagset/flagset.go")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", opts)

}
