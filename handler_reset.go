package main

import (
	"context"
	"fmt"
)

func handlerRest(s *state, cmd command) error {
	if len(cmd.Args) >= 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
	
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Successfully reset the database.")
	return nil
}