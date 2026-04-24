package wechat

import (
	"context"
	"fmt"
	"sync"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/work"
	"github.com/silenceper/wechat/v2/work/config"
	"mochat-api-server/internal/model"
)

var (
	instances     = make(map[uint]*work.Work)
	instancesLock sync.RWMutex
)

func GetWorkApp(corp *model.Corp) *work.Work {
	instancesLock.RLock()
	if app, ok := instances[corp.ID]; ok {
		instancesLock.RUnlock()
		return app
	}
	instancesLock.RUnlock()

	instancesLock.Lock()
	defer instancesLock.Unlock()

	if app, ok := instances[corp.ID]; ok {
		return app
	}

	wc := wechat.NewWechat()
	wc.SetCache(cache.NewMemory())

	cfg := &config.Config{
		CorpID:         corp.WxCorpid,
		CorpSecret:     corp.EmployeeSecret,
		Token:          corp.Token,
		EncodingAESKey: corp.EncodingAesKey,
	}

	app := wc.GetWork(cfg)
	instances[corp.ID] = app
	return app
}

func GetContactApp(corp *model.Corp) *work.Work {
	cfg := &config.Config{
		CorpID:         corp.WxCorpid,
		CorpSecret:     corp.ContactSecret,
		Token:          corp.Token,
		EncodingAESKey: corp.EncodingAesKey,
	}

	wc := wechat.NewWechat()
	wc.SetCache(cache.NewMemory())

	return wc.GetWork(cfg)
}

func RefreshApp(corpID uint) {
	instancesLock.Lock()
	defer instancesLock.Unlock()
	delete(instances, corpID)
}

func UploadTempMedia(ctx context.Context, corp *model.Corp, mediaType string, filePath string) (string, error) {
	app := GetWorkApp(corp)
	material := app.GetMaterial()
	result, err := material.UploadTempFile(filePath, mediaType)
	if err != nil {
		return "", fmt.Errorf("upload temp media failed: %w", err)
	}
	return result.MediaID, nil
}
