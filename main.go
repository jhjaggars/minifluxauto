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

func expireEntries(entries miniflux.Entries, client *miniflux.Client) int {
	entries_actual := []int64{}
	for _, entry := range entries {
		entries_actual = append(entries_actual, entry.ID)
	}
	if len(entries_actual) > 0 {
		client.UpdateEntries(entries_actual, miniflux.EntryStatusRead)
	}
	return len(entries_actual)
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

	parsed_expires := make(map[int64]time.Duration)
	err = viper.UnmarshalKey("miniflux.feeds_expire", &parsed_expires)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	for feedId, expire_duration := range parsed_expires {
		expire_entries, err := client.FeedEntries(feedId, &miniflux.Filter{Status: "unread", Before: time.Now().Add(-expire_duration).Unix()})
		if err != nil {
			fmt.Printf("Error getting feed entries for feedId %d: %s", feedId, err)
			continue
		}
		fmt.Printf("Marking %d entries for feedId: %d\n", expireEntries(expire_entries.Entries, client), feedId)
	}

	parsed_expires = make(map[int64]time.Duration)
	err = viper.UnmarshalKey("miniflux.category_expire", &parsed_expires)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	for categoryId, expire_duration := range parsed_expires {
		expire_entries, err := client.CategoryEntries(categoryId, &miniflux.Filter{Status: "unread", Before: time.Now().Add(-expire_duration).Unix()})
		if err != nil {
			fmt.Printf("Error getting feed entries for categoryId %d: %s", categoryId, err)
			continue
		}
		fmt.Printf("Marking %d entries for categoryId: %d\n", expireEntries(expire_entries.Entries, client), categoryId)
	}
}
