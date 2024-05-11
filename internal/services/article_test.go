package services

import (
	"errors"
	"testing"

	goose "github.com/advancedlogic/GoOse"
)

type MockGooseExtractor struct {
	Article *goose.Article
	Err     error
}

func (m *MockGooseExtractor) ExtractFromURL(url string) (*goose.Article, error) {
	return m.Article, m.Err
}

func TestExtractArticle(t *testing.T) {
	mockArticle := &goose.Article{Title: "Test Article"}
	tests := []struct {
		name        string
		url         string
		mockArticle *goose.Article
		mockError   error
		wantErr     bool
		wantTitle   string
	}{
		{
			name:        "successful extraction",
			url:         "http://example.com",
			mockArticle: mockArticle,
			mockError:   nil,
			wantErr:     false,
			wantTitle:   "Test Article",
		},
		{
			name:      "failed extraction",
			url:       "http://example.com",
			mockError: errors.New("failed to extract"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExtractor := &MockGooseExtractor{Article: tt.mockArticle, Err: tt.mockError}
			service := &ArticleService{gooseExtractor: mockExtractor}
			article, err := service.ExtractArticle(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && article != nil && article.Title != tt.wantTitle {
				t.Errorf("ExtractArticle() got = %v, want %v", article.Title, tt.wantTitle)
			}
		})
	}
}
