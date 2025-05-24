package main

import (
	"errors"
	"log"

	"github.com/itchyny/gojq"
)

func main() {
	query, err := gojq.Parse(".foo | ..")
	if err != nil {
		log.Fatalln(err)
	}
	input := map[string]any{"foo": []any{1, 2, 3}}
	iter := query.Run(input) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			hErr := &gojq.HaltError{}
			if errors.As(err, &hErr) {
				break
			}
			log.Fatalln(hErr)
		}
		log.Printf("%#v\n", v)
	}
}
