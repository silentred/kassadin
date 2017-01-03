package service

import (
	"encoding/json"
	"fmt"
	"net/url"

	redis "gopkg.in/redis.v5"
)

import "github.com/silentred/template/util"

import "strings"

var (
	countries = map[string]string{
		"":   "全球",
		"CN": "中国",
		"US": "美国",
		"TW": "台湾",
		"JP": "日本",
		"TH": "泰国",
		"RU": "俄罗斯",
		"FR": "法国",
		"DE": "德国",
		"GB": "英国",
		"SG": "新加坡",
		"TG": "多哥",
	}
)

const (
	AppInfoCacheKeyFormat = "app:%s:%s"

	AdToken = "1001lpy5"
)

type ItunesService interface {
	GenerateAdLink(bundleID, country, showUID string) (string, int64, error)
}

type ItunesSV struct {
	token    string
	redisCli *redis.Client
}

type AppInfo struct {
	TrackID      int64  `json:"trackId"`
	TrackViewUrl string `json:"trackViewUrl"`
}

func NewItunesSV(token string, cli *redis.Client) *ItunesSV {
	return &ItunesSV{token, cli}
}

func (itune *ItunesSV) searchAllCountryByBundleID(bundleID string, country string) (AppInfo, error) {
	var app AppInfo
	var err error

	c := GetMemCache()
	key := fmt.Sprintf(AppInfoCacheKeyFormat, bundleID, country)
	ret := util.TryCache(c, key, func() interface{} {
		app, err = itune.searchByBundleID(bundleID, country)
		if err == nil {
			return app
		}

		for cty := range countries {
			if cty != country {
				app, err = itune.searchByBundleID(bundleID, cty)
				if err == nil {
					return app
				}
			}
		}
		return nil
	})

	if app, ok := ret.(AppInfo); ok && ret != nil {
		return app, nil
	}

	return app, fmt.Errorf("cannot find app info by bundleID:%s, country:%s", bundleID, country)
}

func (itune *ItunesSV) searchByBundleID(bundleID string, country string) (AppInfo, error) {
	var url string
	var app AppInfo
	var err error

	if len(country) == 0 {
		url = "https://itunes.apple.com/lookup"
	} else {
		url = fmt.Sprintf("https://itunes.apple.com/%s/lookup", country)
	}
	query := map[string]string{
		"bundleId": bundleID,
	}
	config := util.NewReqeustConfig(query, nil, 0, nil, nil)
	body, _, err := util.HTTPGet(url, config)
	if err != nil {
		return app, err
	}

	result := struct {
		ResultCount int       `json:"resultCount"`
		Results     []AppInfo `json:"results"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return app, err
	}

	if len(result.Results) > 0 {
		app = result.Results[0]
		//itune.saveCache(bundleID, country, app)
		return app, nil
	}

	return app, fmt.Errorf("no results got from itunus api query: %v", query)
}

func (itune *ItunesSV) saveCache(bundleID string, country string, app AppInfo) error {
	var err error
	key := fmt.Sprintf(AppInfoCacheKeyFormat, bundleID, country)
	res, err := json.Marshal(app)
	if err != nil {
		return err
	}

	if itune.redisCli != nil {
		err = itune.redisCli.Set(key, util.String(res), 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (itune *ItunesSV) getCache(bundleID string, country string) (AppInfo, error) {
	var app AppInfo

	if itune.redisCli != nil {
		key := fmt.Sprintf(AppInfoCacheKeyFormat, bundleID, country)
		res, err := itune.redisCli.Get(key).Result()
		if err != nil {
			return app, err
		}

		err = json.Unmarshal(util.Slice(res), &app)
		if err != nil {
			return app, err
		}

		return app, nil
	}

	return app, fmt.Errorf("miss cache %s, %s", bundleID, country)
}

func (itune *ItunesSV) GenerateAdLink(bundleID, country, showUID string) (string, int64, error) {
	var userID int64
	var urlStr string

	app, err := itune.searchAllCountryByBundleID(bundleID, country)
	if err != nil {
		return "", 0, err
	}

	_, err = fmt.Sscanf(showUID, "u%d", &userID)
	if err != nil {
		return "", 0, err
	}

	index := strings.Index(app.TrackViewUrl, "?")
	if index > 0 {
		urlStr = app.TrackViewUrl[:index]
	} else {
		urlStr = app.TrackViewUrl
	}

	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", 0, err
	}

	ct := fmt.Sprintf("%s:%d", showUID, app.TrackID)
	appID := fmt.Sprintf("%d", app.TrackID)
	appendQuery := map[string]string{
		"mt":     "8",
		"uo":     "4",
		"ct":     ct,
		"app_id": appID,
	}
	query := urlObj.Query()
	for key, val := range appendQuery {
		query.Add(key, val)
	}
	urlObj.RawQuery = query.Encode()

	return urlObj.String(), app.TrackID, nil
}
