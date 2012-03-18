package openid

import (
  "net/http"
)

// Interface that simplifies testing.
type httpGetter interface {
  Get(uri string, headers map[string]string) (resp *http.Response, err error)
}

type defaultGetter struct{}

var urlGetter = &defaultGetter{}

func (*defaultGetter) Get(uri string, headers map[string]string) (resp *http.Response, err error) {
  request, err := http.NewRequest("GET", uri, nil)
  if err != nil {
    return
  }
  for h, v := range headers {
    request.Header.Add(h, v)
  }
  client := &http.Client{}
  return client.Do(request)
}
