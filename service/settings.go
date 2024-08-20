package service

import (
	"context"

	"jcourse_go/repository"
)

type ISiteSettings interface {
	GetSiteName(ctx context.Context) string
	SetSiteName(ctx context.Context, siteName string) string
}

type SiteSettings struct{}

func (s *SiteSettings) SetSiteName(ctx context.Context, siteName string) string {
	return ""
}

func (s *SiteSettings) GetSiteName(ctx context.Context) string {
	val, err := repository.GetSiteSetting(ctx, "site_name")
	if err != nil {
		return ""
	}
	return val
}

func GetSiteSettings() ISiteSettings {
	return &SiteSettings{}
}
