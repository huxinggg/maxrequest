package maxrequest

import "net/http"

type RetryNotice struct {
	Num    int64  `json:"num"`
	Point  string `json:"point"`
	Result string `json:"result"`
}

type GoResults struct {
	Resp *http.Response
	Body []byte
	Err  error
}
