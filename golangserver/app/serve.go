package main

import (
	"fmt"

	as "github.com/aerospike/aerospike-client-go"
)

func main() {
	// define a client to connect to
	client, err := as.NewClient("aerospike", 3000)
	panicOnError(err)
	fmt.Println("connesso")

	userID := 1
	// prendo la chiave
	namespace := "cibucks"
	setName := "userProfiles"
	key, err := as.NewKey(namespace, setName, as.IntegerValue(userID)) // user key can be of any supported type
	panicOnError(err)

	// read it!
	readPolicy := as.NewPolicy()
	rec, err := client.Get(readPolicy, key)
	panicOnError(err)
	fmt.Printf("here is userProfiles \n\n%#v\n", *rec)
}

func panicOnError(err error) {
	if err != nil {
		println(err)
		panic(err)
	}
}
