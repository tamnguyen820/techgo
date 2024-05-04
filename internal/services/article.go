package services

import (
	"sync"

	goose "github.com/advancedlogic/GoOse"
)

type ArticleService struct {
	gooseExtractor goose.Goose
}

var (
	articleService     *ArticleService
	articleServiceOnce sync.Once
)

// Singleton
func NewArticleService() *ArticleService {
	articleServiceOnce.Do(func() {
		gooseExtractor := goose.New()
		articleService = &ArticleService{gooseExtractor}
	})
	return articleService
}

func (service *ArticleService) ExtractArticle(url string) (*goose.Article, error) {
	article, err := service.gooseExtractor.ExtractFromURL(url)
	if err != nil {
		return nil, err
	}
	return article, nil
}
