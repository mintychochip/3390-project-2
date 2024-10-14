package main

import (
	"api-3390/config"
	"api-3390/database"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Connection(cfg)

	fmt.Println(cfg)
}
func getConfig() (*config.Config, error) {
	if len(os.Args) > 1 {
		cfg, err := config.Load(os.Args[1])
		return cfg, err
	}
	cfg, err := config.LoadFromEnv()
	return cfg, err

}
