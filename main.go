package main

import (
	"sync"

	"meganex/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go nex.StartAuthenticationServer()
	go nex.StartSecureServer()

	wg.Wait()
}
