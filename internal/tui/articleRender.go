package tui

import (
	"fmt"
	"time"
)

type ArticleInfo struct {
	Title         string
	FeedName      string
	Authors       string
	URL           string
	PublishedTime *time.Time
	CleanedText   string
}

func (a ArticleInfo) Render() string {
	return fmt.Sprintf("# %s\n## %s\n## %s\n## %s\n%s\n", a.FeedName, a.Title, a.Authors, a.PublishedTime, a.CleanedText)
}
