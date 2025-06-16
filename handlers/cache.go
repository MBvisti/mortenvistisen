package handlers

import (
	"log/slog"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/maypok86/otter"
	"github.com/mbvisti/mortenvistisen/views"
)

const (
	landingPageCacheKey    = "landingPage"
	articlePageCacheKey    = "articlePage--"
	newsletterPageCacheKey = "newsletterPage--"
)

type CacheManager struct {
	pageCache       otter.Cache[string, templ.Component]
	articleCache    otter.Cache[string, views.ArticlePageProps]
	newsletterCache otter.Cache[string, views.NewsletterPageProps]
}

func NewCacheManager() (*CacheManager, error) {
	pageCacheBuilder, err := otter.NewBuilder[string, templ.Component](20)
	if err != nil {
		return nil, err
	}

	pageCache, err := pageCacheBuilder.WithTTL(48 * time.Hour).Build()
	if err != nil {
		return nil, err
	}

	articleCacheBuilder, err := otter.NewBuilder[string, views.ArticlePageProps](100)
	if err != nil {
		return nil, err
	}

	articleCache, err := articleCacheBuilder.WithTTL(48 * time.Hour).Build()
	if err != nil {
		return nil, err
	}

	newsletterCacheBuilder, err := otter.NewBuilder[string, views.NewsletterPageProps](100)
	if err != nil {
		return nil, err
	}

	newsletterCache, err := newsletterCacheBuilder.WithTTL(48 * time.Hour).Build()
	if err != nil {
		return nil, err
	}

	return &CacheManager{
		pageCache:       pageCache,
		articleCache:    articleCache,
		newsletterCache: newsletterCache,
	}, nil
}

func (cm *CacheManager) GetPageCache() otter.Cache[string, templ.Component] {
	return cm.pageCache
}

func (cm *CacheManager) GetArticleCache() otter.Cache[string, views.ArticlePageProps] {
	return cm.articleCache
}

func (cm *CacheManager) GetNewsletterCache() otter.Cache[string, views.NewsletterPageProps] {
	return cm.newsletterCache
}

func (cm *CacheManager) InvalidateLandingPage() {
	cm.pageCache.Delete(landingPageCacheKey)
	slog.Info("Invalidated landing page cache")
}

func (cm *CacheManager) InvalidateArticle(slug string) {
	key := articlePageCacheKey + slug
	cm.articleCache.Delete(key)
	slog.Info("Invalidated article cache", "slug", slug)
}

func (cm *CacheManager) InvalidateAllArticles() {
	keys := make([]string, 0)
	cm.articleCache.Range(func(key string, _ views.ArticlePageProps) bool {
		if strings.HasPrefix(key, articlePageCacheKey) {
			keys = append(keys, key)
		}
		return true
	})

	for _, key := range keys {
		cm.articleCache.Delete(key)
	}

	slog.Info("Invalidated all article caches", "count", len(keys))
}

func (cm *CacheManager) InvalidateNewsletter(slug string) {
	key := newsletterPageCacheKey + slug
	cm.newsletterCache.Delete(key)
	slog.Info("Invalidated newsletter cache", "slug", slug)
}

func (cm *CacheManager) InvalidateAllNewsletters() {
	keys := make([]string, 0)
	cm.newsletterCache.Range(func(key string, _ views.NewsletterPageProps) bool {
		if strings.HasPrefix(key, newsletterPageCacheKey) {
			keys = append(keys, key)
		}
		return true
	})

	for _, key := range keys {
		cm.newsletterCache.Delete(key)
	}

	slog.Info("Invalidated all newsletter caches", "count", len(keys))
}
