package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"net/http"
	"testing"
)

func Sync(ctx context.Context, name string, url string) (err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	if _, err = client.Object.Put(ctx, name, resp.Body, nil); err != nil {
		return
	}
	return
}

func TestSyncEditorJS(t *testing.T) {
	ctx := context.TODO()
	source := [][]string{
		{"editorjs/editorjs.js", "@editorjs/editorjs"},
		{"editorjs/paragraph.js", "@editorjs/paragraph"},
		{"editorjs/header.js", "@editorjs/header"},
		{"editorjs/delimiter.js", "@editorjs/delimiter"},
		{"editorjs/underline.js", "@editorjs/underline"},
		{"editorjs/list.js", "@editorjs/list"},
		{"editorjs/nested-list.js", "@editorjs/nested-list"},
		{"editorjs/checklist.js", "@editorjs/checklist"},
		{"editorjs/table.js", "@editorjs/table"},
	}

	urls := make([]string, 0)
	for _, v := range source {
		err := Sync(ctx,
			fmt.Sprintf(`/assets/%s`, v[0]),
			fmt.Sprintf(`https://cdn.jsdelivr.net/npm/%s`, v[1]),
		)
		assert.Nil(t, err)
		urls = append(urls, fmt.Sprintf(`https://cdn.kainonly.com/assets/%s`, v[0]))
	}

	credential := common.NewCredential(
		values.COS.AccessKeyID,
		values.COS.AccessKeySecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	manage, _ := cdn.NewClient(credential, "", cpf)
	request := cdn.NewPurgeUrlsCacheRequest()
	request.Urls = common.StringPtrs(urls)
	_, err := manage.PurgeUrlsCache(request)
	assert.Nil(t, err)
}
