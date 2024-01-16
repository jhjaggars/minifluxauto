package main

import (
	"fmt"
	"time"

	viper "github.com/spf13/viper"
	miniflux "miniflux.app/client"
)

type ExpireElement struct {
	expire_duration time.Duration `mapstructure:"expire_duration"`
}

func expire(client *miniflux.Client, feed_id int64, expire_duration time.Duration) {
	expire_entries, err := client.FeedEntries(feed_id, &miniflux.Filter{Status: "unread", Before: time.Now().Add(-expire_duration).Unix()})
	if err != nil {
		fmt.Printf("Error getting feed entries for feed_id: %d: %v\n", feed_id, err)
		return
	}
	entries_actual := []int64{}
	for _, entry := range expire_entries.Entries {
		entries_actual = append(entries_actual, entry.ID)
	}
	if len(entries_actual) > 0 {
		fmt.Println("Marking", len(entries_actual), "entries for feed_id", feed_id)
		client.UpdateEntries(entries_actual, miniflux.EntryStatusRead)
	} else {
		fmt.Println("No entries to mark for feed_id", feed_id)
	}
}

func main() {
	viper.SetConfigName("minifluxauto")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	client := miniflux.New(viper.GetString("miniflux.url"), viper.GetString("miniflux.token"))

	// categories
	category_expires := make(map[int64]time.Duration)
	err = viper.UnmarshalKey("miniflux.categories_expire", &category_expires)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	for category_id, expire_duration := range category_expires {
		feeds, err := client.CategoryFeeds(category_id)
		if err != nil {
			fmt.Printf("Error getting feeds for category_id: %d: %v\n", category_id, err)
			continue
		}
		fmt.Println("Expiring Category", category_id, "with", len(feeds), "feeds")
		for _, feed := range feeds {
			expire(client, feed.ID, expire_duration)
		}
	}

	// individual feeds
	parsed_expires := make(map[int64]time.Duration)
	err = viper.UnmarshalKey("miniflux.feeds_expire", &parsed_expires)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	for feed_id, expire_duration := range parsed_expires {
		expire(client, feed_id, expire_duration)
	}
}
