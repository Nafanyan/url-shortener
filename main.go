package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	// Теперь можно использовать cfg
	fmt.Println(cfg)
}
