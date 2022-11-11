// Package gcal Powerful and easy to use http client
package gcalx

import "sync"

// var Defaults = DefaultClient.Defaults
// The default client for convenience
func DefaultClient() *HttpClient {
	return &HttpClient{
		reuseTransport: true,
		reuseJar:       true,
		lock:           new(sync.Mutex),
	}
}

//var Begin = DefaultClient.Begin
//var Do = DefaultClient.Do
//var Get = DefaultClient.Get
//var Delete = DefaultClient.Delete
//var Head = DefaultClient.Head
//var Post = DefaultClient.Post
//var PostJson = DefaultClient.PostJson
//var PostMultipart = DefaultClient.PostMultipart
//var Put = DefaultClient.Put
//var PutJson = DefaultClient.PutJson
//var WithOption = DefaultClient.WithOption
//var WithOptions = DefaultClient.WithOptions
//var WithHeader = DefaultClient.WithHeader
//var WithHeaders = DefaultClient.WithHeaders
//var WithITrace = DefaultClient.WithITrace
//var WithCommonHeader = DefaultClient.WithCommonHeader
