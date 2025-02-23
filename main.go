package main

import (
	"fmt"
	"log"

	config "github.com/KrishKoria/GatorConfig"
)
func main()  {
    cfg, err := config.Read()
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    err = cfg.SetUser("Krish")
    if err != nil {
        log.Fatalf("Error setting user: %v", err)
    }

    updatedCfg, err := config.Read()
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    fmt.Printf("Updated Config: %+v\n", updatedCfg)
}