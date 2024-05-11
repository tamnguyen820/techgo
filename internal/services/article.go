package services

import (
	"sync"

	goose "github.com/advancedlogic/GoOse"
)

// GooseExtractor interface abstracts the functionality of goose.Goose for testing
type GooseExtractor interface {
	ExtractFromURL(url string) (*goose.Article, error)
}

type ArticleService struct {
	gooseExtractor GooseExtractor
}

var (
	articleService     *ArticleService
	articleServiceOnce sync.Once
)

func NewArticleService() *ArticleService {
	articleServiceOnce.Do(func() {
		gooseExtractor := goose.New()
		articleService = &ArticleService{gooseExtractor}
	})
	return articleService
}

func (service *ArticleService) ExtractArticle(url string) (*goose.Article, error) {
	return service.gooseExtractor.ExtractFromURL(url)
}
