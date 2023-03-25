package main

import (
	_ "github.com/lib/pq"

	"github.com/rmiguelac/tasker/internal/api"
)

func main() {

	s := api.New(":8000")
	s.Run()
}
