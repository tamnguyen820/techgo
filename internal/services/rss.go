package services

import (
	"os"
	"sync"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

type RSSService struct {
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

var (
	rssService     *RSSService
	rssServiceOnce sync.Once
)

// Singleton
func NewRSSService(configFilePath string) *RSSService {
	rssServiceOnce.Do(func() {
		if len(configFilePath) == 0 {
			configFilePath = "config.yml" // default config file path
		}
		rssService = &RSSService{configFilePath}
	})
	return rssService
}

func (service *RSSService) FetchAllFeeds() ([]*RSSFeed, error) {
	configFile, err := os.ReadFile(service.configFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	allFeeds, err := fetchAllFeedsInParallel(config)
	if err != nil {
		return nil, err
	}

	return allFeeds, nil
}

func fetchAllFeedsInParallel(config Config) ([]*RSSFeed, error) {
	numFeeds := len(config.Feeds)
	var wg sync.WaitGroup
	wg.Add(numFeeds)
	allFeeds := make([]*RSSFeed, numFeeds)
	errChan := make(chan error, numFeeds)

	// Spawn goroutines to fetch each feed concurrently
	for i, rssFeed := range config.Feeds {
		go func(index int, rssFeed RSSFeedInfo) {
			defer wg.Done()
			feed, err := fetchFeed(rssFeed.URL, gofeed.NewParser())
			if err != nil {
				errChan <- err
				return
			}
			allFeeds[index] = &RSSFeed{FeedInfo: rssFeed, Feed: feed}
		}(i, rssFeed)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return allFeeds, nil
}

type FeedParser interface {
	ParseURL(url string) (*gofeed.Feed, error)
}

func fetchFeed(url string, parser FeedParser) (*gofeed.Feed, error) {
	feed, err := parser.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, nil
}
