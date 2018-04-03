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

	namespace := "cibucks"
	setName := "userProfiles"
	key, err := as.NewKey(namespace, setName, "1") // user key can be of any supported type
	panicOnError(err)
	fmt.Println("ecco  " + key.String())

	// read it back!
	readPolicy := as.NewPolicy()
	rec, err := client.Get(readPolicy, key)
	panicOnError(err)
	fmt.Printf("eccoli userProfiles \n\n%#v\n", *rec)

	/*
				// define some bins
				bins := as.BinMap{
					"bin1": 42, // you can pass any supported type as bin value
					"bin2": "An elephant is a mouse with an operating system",
					"bin3": []interface{}{"Go", 17981},
				}

				// write the bins
				writePolicy := as.NewWritePolicy(0, 0)
				err = client.Put(writePolicy, key, bins)
				panicOnError(err)

				// read it back!
				readPolicy := as.NewPolicy()
				rec, err := client.Get(readPolicy, key)
				panicOnError(err)


		 			// Add to bin1
					err = client.Add(writePolicy, key, as.BinMap{"bin1": 1})
					panicOnError(err)

					rec2, err := client.Get(readPolicy, key)
					panicOnError(err)

					fmt.Println("value of %s: %v\n", "bin1", rec2.Bins["bin1"])

					// prepend and append to bin2
					err = client.Prepend(writePolicy, key, as.BinMap{"bin2": "Frankly:  "})
					panicOnError(err)
					err = client.Append(writePolicy, key, as.BinMap{"bin2": "."})
					panicOnError(err)

					rec3, err := client.Get(readPolicy, key)
					panicOnError(err)

					fmt.Println("value of %s: %v\n", "bin2", rec3.Bins["bin2"])

					// delete bin3
						err = client.Put(writePolicy, key, as.BinMap{"bin3": nil})
						rec4, err := client.Get(readPolicy, key)
						panicOnError(err)

						fmt.Println("bin3 does not exist anymore: %#v\n", *rec4)

						// check if key exists
						exists, err := client.Exists(readPolicy, key)
						panicOnError(err)
						fmt.Println("key exists in the database: %#v\n", exists)

						// delete the key, and check if key exists
						existed, err := client.Delete(writePolicy, key)
						panicOnError(err)
						fmt.Println("did key exist before delete: %#v\n", existed)

						exists, err = client.Exists(readPolicy, key)
						panicOnError(err)
						fmt.Println("key exists: %#v\n", exists)
	*/
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
