package main

import (
	"fmt"
	"github.com/slavik22/urlShortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
