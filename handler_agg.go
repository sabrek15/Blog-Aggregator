package main

import (
	"fmt"
	"time"
)


func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected one command argument")
	}

	timebetweenreqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse time duration: %w", err)
	}

	fmt.Printf("collecting feeds every: %s", timebetweenreqs)
	ticker := time.NewTicker(timebetweenreqs)
	for ; ; <- ticker.C {
		scrapeFeeds(s)
	}

}