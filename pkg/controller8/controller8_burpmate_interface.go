package controller8

import "deifzar/num8/pkg/model8"

type Controller8BurpmateInterface interface {
	GetSitemapByURLPrefix(urlPrefix string) (*model8.Sitemap8, error)
	GetSitemapFilteredOut(urlPrefix string, contenttype []string, statuscode []string) (*model8.Sitemap8, error)
	SendSitemap(*model8.Sitemap8)
}
