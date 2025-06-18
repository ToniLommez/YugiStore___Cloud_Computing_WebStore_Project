// Sample run-helloworld is a minimal Cloud Run service.
// https://yugistore-650237781960.us-central1.run.app

package main

import (
	"yugistore/server"
)

func main() {
	server.StartServer()
}
