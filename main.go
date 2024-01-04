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

func main() {
	viper.SetConfigName("minifluxauto")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
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
	for feed_id, expire_duration := range parsed_expires {
		expire_entries, err := client.FeedEntries(feed_id, &miniflux.Filter{Status: "unread", Before: time.Now().Add(-expire_duration).Unix()})
		if err != nil {
			fmt.Println("Error getting feed entries for feed_id", feed_id)
			continue
		}
		entries_actual := []int64{}
		for _, entry := range expire_entries.Entries {
			entries_actual = append(entries_actual, entry.ID)
		}
		fmt.Println("Marking", len(entries_actual), "entries for feed_id", feed_id)
		client.UpdateEntries(entries_actual, miniflux.EntryStatusRead)
	}
}
