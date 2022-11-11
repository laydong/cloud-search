package gcalx

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var sdkLog = true

const (
	timeFormat = "2006-01-02 15:04:05"
)

var b3Headers = []string{"x-request-id", "x-b3-traceid", "x-b3-spanid", "x-b3-parentspanid", "x-b3-sampled", "x-b3-flags", "x-ot-span-context", "x-huayu-traffic-tag"}

type RequestLog struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	IP      []string          `json:"ip"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Query   interface{}       `json:"query"`
	Body    interface{}       `json:"body"`
}

// Add params to a url string.
func addParams(url_ string, params url.Values) string {
	if len(params) == 0 {
		return url_
	}

	if !strings.Contains(url_, "?") {
		url_ += "?"
	}

	if strings.HasSuffix(url_, "?") || strings.HasSuffix(url_, "&") {
		url_ += params.Encode()
	} else {
		url_ += "&" + params.Encode()
	}

	return url_
}

// Add a file to a multipart writer.
func addFormFile(writer *multipart.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)

	return err
}

// Option Convert options with string keys to desired format.
func Option(o map[string]interface{}) map[int]interface{} {
	rst := make(map[int]interface{})
	for k, v := range o {
		k := "OPT_" + strings.ToUpper(k)
		if num, ok := CONST[k]; ok {
			rst[num] = v
		}
	}

	return rst
}

// Merge options(latter ones have higher priority)
func mergeOptions(options ...map[int]interface{}) map[int]interface{} {
	rst := make(map[int]interface{})

	for _, m := range options {
		for k, v := range m {
			rst[k] = v
		}
	}

	return rst
}

// Merge headers(latter ones have higher priority)
func mergeHeaders(headers ...map[string]string) map[string]string {
	rst := make(map[string]string)

	for _, m := range headers {
		for k, v := range m {
			rst[k] = v
		}
	}

	return rst
}

// Does the params contain a file?
func checkParamFile(params url.Values) bool {
	for k, _ := range params {
		if k[0] == '@' {
			return true
		}
	}

	return false
}

// Is opt in options?
func hasOption(opt int, options []int) bool {
	for _, v := range options {
		if opt == v {
			return true
		}
	}

	return false
}

// Map is a mixed structure with options and headers
type Map map[interface{}]interface{}

// Parse the Map, return options and headers
func parseMap(m Map) (map[int]interface{}, map[string]string) {
	var options = make(map[int]interface{})
	var headers = make(map[string]string)

	if m == nil {
		return options, headers
	}

	for k, v := range m {
		// integer is option
		if kInt, ok := k.(int); ok {
			// don't need to validate
			options[kInt] = v
		} else if kString, ok := k.(string); ok {
			kStringUpper := strings.ToUpper(kString)
			if kInt, ok := CONST[kStringUpper]; ok {
				options[kInt] = v
			} else {
				// it should be header, but we still need to validate it's type
				if vString, ok := v.(string); ok {
					headers[kString] = vString
				}
			}
		}
	}

	return options, headers
}

func toUrlValues(v interface{}) url.Values {
	switch t := v.(type) {
	case url.Values:
		return t
	case map[string][]string:
		return url.Values(t)
	case map[string]string:
		rst := make(url.Values)
		for k, v := range t {
			rst.Add(k, v)
		}
		return rst
	case nil:
		return make(url.Values)
	default:
		panic("Invalid value")
	}
}

func checkParamsType(v interface{}) int {
	switch v.(type) {
	case url.Values, map[string][]string, map[string]string:
		return 1
	case []byte, string, *bytes.Reader:
		return 2
	case nil:
		return 0
	default:
		return 3
	}
}

func toReader(v interface{}) *bytes.Reader {
	switch t := v.(type) {
	case []byte:
		return bytes.NewReader(t)
	case string:
		return bytes.NewReader([]byte(t))
	case *bytes.Reader:
		return t
	case nil:
		return bytes.NewReader(nil)
	default:
		panic("Invalid value")
	}
}

func copyRequest(r *http.Request) *http.Request {
	req := &http.Request{}
	req = r
	return req
}

func SetSdkLog(v bool) {
	sdkLog = v
}
