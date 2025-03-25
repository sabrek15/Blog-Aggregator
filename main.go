package main

import (
	"fmt"
	"github.com/sabrek15/Blog-Aggregator/internal/config"
)


func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		return
	}
	
	err = cfg.SetUser("sabrek15")
	if err != nil {
		fmt.Println("error updating config: ", err)
		return
	}

	updatedCfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		return
	}
	
	// fmt.Print(updatedCfg.DBURL)

	fmt.Printf("updated config: %+v\n", updatedCfg)
}