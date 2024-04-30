package rss

import (
	"os"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

type RSSFeedInfo struct {
	URL string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Feeds []RSSFeedInfo `yaml:"rss_feeds"`
}

type RSSFeed struct {
	FeedInfo RSSFeedInfo
	Feed *gofeed.Feed
}

func FetchAllFeeds() ([]*RSSFeed, error) {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	var allFeeds []*RSSFeed
	for _, rssFeed := range config.Feeds {
		feed, err := fetchFeed(rssFeed.URL)
		if err == nil {
			allFeeds = append(allFeeds, &RSSFeed{FeedInfo: rssFeed, Feed: feed})
		}
	}
	return allFeeds, nil
}

func fetchFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, nil
}