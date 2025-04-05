package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sabrek15/gator/internal/database"
)


func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <feedname> <feedURL>", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	
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

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("successfully created feed follow: ")
	PrintFeedFollow(follow)
	fmt.Println("======================================")

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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <feedURL>", cmd.Name)
	}
	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't find the feed: %w",err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("successfully created feed follow: ")
	PrintFeedFollow(follow)
	return nil
}

func hanlderFollowing(s *state, cmd command, user database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't fetch the following feed of login user: %w", err)
	}
	
	for _, follow := range following {
		PrintFeedFollowForUser(follow)
		fmt.Println()
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feedURL>", cmd.Name)
	}
	
	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	err = s.db.DeteleFeedFollow(context.Background(), database.DeteleFeedFollowParams{UserID: user.ID, FeedID: feed.ID})

	if err != nil {
		return 	fmt.Errorf("couldn't unfollow the feed: %w", err)
	}

	fmt.Println("Feed Unfollowed successfully")

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch the next feed: %w", err)
	}
	_ , err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now().UTC(),
		ID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't update the feed: %w", err)
	}

	feeds, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't find the feed: %w",err)
	}

	for _, feed := range feeds.Channel.Item {
		fmt.Printf("* Title: 	%s\n", feed.Title)	
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

func PrintFeedFollow(follow database.CreateFeedFollowRow){
	fmt.Printf("* ID:            %s\n", follow.ID)
	fmt.Printf("* Created:       %v\n", follow.CreatedAt)
	fmt.Printf("* Updated:       %v\n", follow.UpdatedAt)
	fmt.Printf("* FeedID:        %s\n", follow.FeedID)
	fmt.Printf("* UserID:        %s\n", follow.UserID)
}

func PrintFeedFollowForUser(follow database.GetFeedFollowsForUserRow){
	fmt.Printf("* FeedName:      %s\n", follow.FeedName)
	fmt.Printf("* UserName:		 %s\n", follow.UserName)
}