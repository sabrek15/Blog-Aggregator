package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch the next feed: %w", err)
	}
	_ , err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now().UTC(),
		ID: feedToFetch.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't update the feed: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("couldn't find the feed: %w",err)
	}

	for _, feed := range rssFeed.Channel.Item {
		parsedPubTime, err := time.Parse(time.RFC1123Z, feed.PubDate)
		if err != nil {
			fmt.Printf("couldn't parse publish time of post")
			continue
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Title: feed.Title,
			Url: feed.Link,
			Description: feed.Description,
			PublishedAt: parsedPubTime,
			FeedID: feedToFetch.ID,
		})

		if err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
				continue
			}
			fmt.Printf("an error while adding post to Database: %s", feed.Title)
			continue
		}
	}

	return nil
}

func hanlderBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args)==1 {
		i, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("couldn't convert string to integer: %w", err)
		}
		limit = i
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments: expected 1")
	}

	posts, err := s.db.GetUserPosts(context.Background(), database.GetUserPostsParams{
		ID: user.ID,
		Limit: int32(limit),
	})

	if err !=  nil {
		return fmt.Errorf("couldn't get user posts: %w", err)
	}
	
	fmt.Printf("Posts of user:\t%v\n", user.Name)
	fmt.Println("--------------------------------")
	for _, post := range posts {
		fmt.Printf("Title:\t%v\n", post.Title)
		fmt.Printf("Desc:\t%v\n", post.Description)
		fmt.Printf("Feed:\t%v\n", post.FeedName)
		fmt.Printf("----------------------------\n");
		
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