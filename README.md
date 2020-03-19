# go-call-em-all
Go library for interacting with call-em-all's API. This implements the minimum needed to create contacts & groups.

You'll have to talk to them about getting access for your company.

## Use
```go
package main

import (
	"log"

	callemall "github.com/MikeAlbertFleetSolutions/go-call-em-all"
)

func main() {
	// connection to call-em-all
	callemallClient, err := callemall.NewClient("appkey", "secretkey", "authtoken", "https://rest.call-em-all.com")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	l, err := callemallClient.CreateList("test list")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	_, err = callemallClient.CreateContact("test", "user", "555-555-5555", []callemall.List{l})
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
```