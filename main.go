package main

import (
	_ "github.com/lib/pq"

	tserver "github.com/rmiguelac/tasker/internal/server/http"
)

func main() {
	tserver.HandleRequests()

}
