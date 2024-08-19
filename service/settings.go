package service

import (
	"context"

	"jcourse_go/repository"
)

type ISiteSettings interface {
	SiteName(ctx context.Context) string
}

type SiteSettings struct{}

func (s *SiteSettings) SiteName(ctx context.Context) string {
	val, err := repository.GetSiteSetting(ctx, "site_name")
	if err != nil {
		return ""
	}
	return val
}
