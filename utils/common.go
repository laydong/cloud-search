package utils

import (
	"codesearch/global/gcalx"
	"github.com/gin-gonic/gin"
	"io"
	"path"

	"strings"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	RequestIdKey  = "request_id" // 日志key
	OPT_TIMEOUT   = 10
)

// HttpGet http get请求
func HttpGet(c *gin.Context, urls string, head map[string]string) (resp []byte, err error) {
	client := gcalx.DefaultClient().WithITrace(c).WithOption(OPT_TIMEOUT, 60).WithHeader("Content-Type", "application/json")
	if len(head) > 0 {
		for k, v := range head {
			client = client.WithHeader(k, v)
		}
	}
	res, err := client.Get(urls, nil)
	if err != nil {
		return
	}
	defer res.Body.Close()
	resp, err = io.ReadAll(res.Body)
	return
}

// HttpPost http POST请求
func HttpPost(c *gin.Context, url string, data map[string]interface{}, head map[string]string) (body []byte, err error) {
	client := gcalx.DefaultClient().WithITrace(c).WithOption(OPT_TIMEOUT, 60).WithHeader("Content-Type", "application/json")
	if len(head) > 0 {
		for k, v := range head {
			client = client.WithHeader(k, v)
		}
	}
	resp, err := client.PostJson(url, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	return
}

// SupeString 字符串匹配
func SupeString(str string) bool {
	var esFile = [...]string{
		"application/octet-stream",
		"image/gif",
		"image/jpeg",
		"image/png",
		"image/fax",
		"image/gif",
		"image/tiff",
		"image/x-icon",
		"image/vnd.rn-realpix",
		"image/vnd.wap.wbmp",
		"video/x-ms-asf",
		"video/avi",
		"video/x-ivf",
		"video/x-mpeg",
		"video/mpeg4",
		"video/x-sgi-movie",
		"video/mpeg",
		"video/mpg",
		"video/vnd.rn-realvideo",
		"video/x-ms-wm",
		"video/x-ms-wmv",
		"video/x-ms-wmx",
		"video/x-ms-wvx",
		"application/pdf",
		"application/msword",
		"application/x-jpg",
		"application/x-jpe",
		"application/x-img",
		"application/x-msdownload",
		"application/x-ico",
		"application/vnd.ms-powerpoint",
		"application/x-ppt",
		"application/vnd.android.package-archive",
		"audio/mp1",
		"audio/mp2",
		"audio/mp3",
		"audio/rn-mpeg",
	}
	for _, v := range esFile {
		if strings.Contains(v, str) {
			return true
		}
	}
	return false
}

// FileExt 检查文件路径后缀
func FileExt(fileStr string) bool {
	fileStr = path.Ext(fileStr)
	var esFile = [...]string{
		".PNG",
		".gif",
		".jpeg",
		".png",
		".fax",
		".GIF",
		".tiff",
		".icon",
		".exe",
		".dll",
		".Zip",
		".avi",
		".mpeg",
		".mpg",
		".pdf",
		".mp4",
		".mp3",
		".docx",
		".csv",
		".xlsx",
		".md",
		".crdownload",
		".svga",
		".plist",
		".xcuserstate",
		".strings",
		".mobileprovision",
		".jks",
	}
	for _, v := range esFile {
		if strings.Contains(v, fileStr) {
			return true
		}
	}
	return false
}

// InSliceString string是否在[]string里面
func InSliceString(k string, s []string) bool {
	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false
}
