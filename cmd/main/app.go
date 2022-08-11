package main

import (
	"fmt"

	"github.com/VrMolodyakov/vote-service/internal/config"
)

func main() {
	fmt.Println("start")
	cfg := config.GetConfig()
	fmt.Println(cfg)
}
