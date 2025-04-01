package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sabrek15/gator/internal/database"
)


func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <feedname> <feedURL>", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)

	if err != nil {
		return fmt.Errorf("couldn't find the user: %w", err)
	}
	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: feedName,
		Url: feedURL,
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Println("Feed added successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetchfeeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't find the username: %w", err)
		}
		printFeeds(feed, user)
		fmt.Println()
		fmt.Println("=====================================")
	}
	return nil
}

func printFeed(feed database.Feed){
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func printFeeds(feed database.Feed, user database.User){
	fmt.Printf("* Name:          	%s\n", feed.Name)
	fmt.Printf("* URL:           	%s\n", feed.Url)
	fmt.Printf("* UserName:        	%s\n", user.Name)

}