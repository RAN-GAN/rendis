package main

import (
	"fmt"
	"log"

	rendis "github.com/RAN-GAN/client/go"
)

func main() {
	client, err := rendis.New("ws://localhost:8080", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	err = client.Set("name", "Pradeep")
	if err != nil {
		log.Fatal(err)
	}

	name, err := client.Get("name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name:", name)

	ok, err := client.Exists("name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists:", ok)

	ttl, err := client.TTL("name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TTL:", ttl)

	ok, err = client.Expire("name", 60)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Expire:", ok)

	deleted, err := client.Del("name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted:", deleted)
}
