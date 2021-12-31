package main

import (
	"fmt"
	"log"

	"github.com/okamotoke/realaddress"
)

func main() {
	c, err := realaddress.GetRandomAddress()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.GetPostalCode())
	fmt.Println(c.GetPrefecture())
	fmt.Println(c.GetCity())
	fmt.Println(c.GetTown())
}
