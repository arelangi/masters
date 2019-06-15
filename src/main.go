package main

import (
	"fmt"

	"gopkg.in/jdkato/prose.v2"
)

func main() {
	doc, _ := prose.NewDocument("John gave a speech in Hyderabad for TDP")
	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
		// Lebron James PERSON
		// Los Angeles GPE
	}
}
