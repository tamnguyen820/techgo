package rss

import (
	"os"
	"sync"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

type FetchService struct {
	configFilePath string
}

type RSSFeedInfo struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Feeds []RSSFeedInfo `yaml:"rss_feeds"`
}

type RSSFeed struct {
	FeedInfo RSSFeedInfo
	Feed     *gofeed.Feed
}

func NewFetchService(configFilePath string) *FetchService {
	if len(configFilePath) == 0 {
		configFilePath = "config.yml" // default config file path
	}
	return &FetchService{configFilePath: configFilePath}
}

func (service *FetchService) FetchAllFeeds() ([]*RSSFeed, error) {
	configFile, err := os.ReadFile(service.configFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	allFeeds, err := fetchAllFeedsParallel(config)
	return allFeeds, err
}

func fetchAllFeedsParallel(config Config) ([]*RSSFeed, error) {
	var wg sync.WaitGroup
	feedChan := make(chan *RSSFeed)

	// Spawn goroutines to fetch each feed concurrently
	for _, rssFeed := range config.Feeds {
		wg.Add(1)
		go func(rssFeed RSSFeedInfo) {
			defer wg.Done()
			feed, err := fetchFeed(rssFeed.URL)
			if err != nil {
				return
			}
			feedChan <- &RSSFeed{FeedInfo: rssFeed, Feed: feed}
		}(rssFeed)
	}

	// Aggregate fetched feeds
	var allFeeds []*RSSFeed
	go func() {
		for feed := range feedChan {
			allFeeds = append(allFeeds, feed)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(feedChan)

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
