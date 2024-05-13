package services

import (
	"errors"
	"testing"

	"github.com/mmcdole/gofeed"
)

type MockFeedParser struct {
	Feed *gofeed.Feed
	Err  error
}

func (m *MockFeedParser) ParseURL(url string) (*gofeed.Feed, error) {
	return m.Feed, m.Err
}

func TestFetchFeed(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		mockFeed  *gofeed.Feed
		mockError error
		wantErr   bool
		wantTitle string
	}{
		{
			name:      "successful feed fetch",
			url:       "http://example.com/feed",
			mockFeed:  &gofeed.Feed{Title: "Example Feed"},
			mockError: nil,
			wantErr:   false,
			wantTitle: "Example Feed",
		},
		{
			name:      "failed feed fetch due to network error",
			url:       "http://example.com/feed",
			mockFeed:  nil,
			mockError: errors.New("network error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := &MockFeedParser{Feed: tt.mockFeed, Err: tt.mockError}
			gotFeed, err := fetchFeed(tt.url, mockParser)

			if (err != nil) != tt.wantErr {
				t.Errorf("fetchFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && gotFeed.Title != tt.wantTitle {
				t.Errorf("fetchFeed() got title = %v, want %v", gotFeed.Title, tt.wantTitle)
			}
		})
	}
}
