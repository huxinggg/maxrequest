package maxrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type MaxRequestAttr interface {
	SetTimeout(sec int64) MaxRequestAttr
	Get(url string) MaxRequestAttr
	Post(url string, body interface{}) MaxRequestAttr
	PostForm(url string, body interface{}) MaxRequestAttr
	SetHeader(key, val string) MaxRequestAttr
	SetRetry(num, interval int64, notice func(retryNotice RetryNotice) bool) MaxRequestAttr
	Result(in interface{}) (resp *http.Response, body []byte, err error)
}

type attr struct {
	Url           string
	Method        string
	Headers       map[string]string
	TimeOut       int64
	LastRetryNum  int64
	RetryNum      int64
	RetryInterval int64
	PostBody      interface{}
	NoticeFunc    func(retryNotice RetryNotice) bool
}

func New() MaxRequestAttr {
	return &attr{
		TimeOut:       10,
		RetryInterval: 5,
	}
}

func (s *attr) SetTimeout(sec int64) MaxRequestAttr {
	s.TimeOut = sec
	return s
}

//设置重试
func (s *attr) SetRetry(num, interval int64, noticeFunc func(retryNotice RetryNotice) bool) MaxRequestAttr {
	s.LastRetryNum = num
	s.RetryInterval = interval
	s.NoticeFunc = noticeFunc
	return s
}

func (s *attr) SetHeader(key, val string) MaxRequestAttr {
	if s.Headers == nil {
		s.Headers = make(map[string]string)
	}
	s.Headers[key] = val
	return s
}

func (s *attr) Get(url string) MaxRequestAttr {
	s.Url = url
	s.Method = "GET"
	return s
}

func (s *attr) Post(url string, body interface{}) MaxRequestAttr {
	s.Url = url
	s.Method = "POST"
	s.PostBody = body
	return s
}

func (s *attr) PostForm(url string, body interface{}) MaxRequestAttr {
	s.Url = url
	s.Method = "POST_FORM"
	s.PostBody = body
	s.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return s
}

func (s *attr) Result(in interface{}) (resp *http.Response, body []byte, err error) {
	postByte := make([]byte, 0)
	if s.Method != "GET" {
		if s.Method == "POST_FORM" {
			s.Method = "POST"
			formValues := url.Values{}
			objT := reflect.TypeOf(s.PostBody)
			objV := reflect.ValueOf(s.PostBody)
			for i := 0; i < objT.NumField(); i++ {
				fileName, ok := objT.Field(i).Tag.Lookup("json")
				if ok {
					formValues.Set(fileName, fmt.Sprintf("%v", objV.Field(i).Interface()))
				} else {
					formValues.Set(objT.Field(i).Name, fmt.Sprintf("%v", objV.Field(i).Interface()))
				}
			}
			fmt.Println(formValues)
			postByte = []byte(formValues.Encode())
		} else {
			if reflect.TypeOf(s.PostBody).String() == "string" {
				postByte = []byte(s.PostBody.(string))
			} else if reflect.TypeOf(s.PostBody).String() == "[]uint8" {
				postByte = s.PostBody.([]uint8)
			} else {
				postByte, err = json.Marshal(s.PostBody)
				if err != nil {
					return
				}
			}
		}
	}
	request, e := http.NewRequest(s.Method, s.Url, bytes.NewReader(postByte))
	if e != nil {
		err = e
		return
	}
	client := http.Client{Timeout: time.Duration(s.TimeOut) * time.Second}
	for k, v := range s.Headers {
		request.Header.Set(k, v)
	}
	resp, err = client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if in != nil {
		err = json.Unmarshal(body, in)
		if err != nil {
			return
		}
		if s.LastRetryNum > 0 { //设置了重试
			time.Sleep(time.Second * time.Duration(s.RetryInterval))
			s.LastRetryNum--
			s.RetryNum++
			objects := reflect.ValueOf(in)
			typeOfType := objects.Elem().Type()
			for i := 0; i < objects.Elem().NumField(); i++ {
				field := objects.Elem().Field(i)
				if typeOfType.Field(i).Tag.Get("maxrequestRetry") == fmt.Sprintf("%v", field.Interface()) {
					if s.NoticeFunc != nil {
						stop := s.NoticeFunc(RetryNotice{
							Num:    s.RetryNum,
							Point:  fmt.Sprintf("%v=%s", typeOfType.Field(i).Name, typeOfType.Field(i).Tag.Get("maxrequestRetry")),
							Result: string(body),
						})
						if stop {
							s.LastRetryNum = 0
							return s.Result(in)
						}
					}
					resp, body, err = s.Result(in)
					if s.LastRetryNum == 0 {
						err = fmt.Errorf("Retry times exhausted")
						return
					}
					return
				}
			}
		} else if s.LastRetryNum < 0 { //无限重试
			s.RetryNum++
			time.Sleep(time.Second * time.Duration(s.RetryInterval))
			objects := reflect.ValueOf(in)
			typeOfType := objects.Elem().Type()
			for i := 0; i < objects.Elem().NumField(); i++ {
				field := objects.Elem().Field(i)
				if typeOfType.Field(i).Tag.Get("maxrequestRetry") == fmt.Sprintf("%v", field.Interface()) {
					if s.NoticeFunc != nil {
						stop := s.NoticeFunc(RetryNotice{
							Num:    s.RetryNum,
							Point:  fmt.Sprintf("%v=%s", typeOfType.Field(i).Name, typeOfType.Field(i).Tag.Get("maxrequestRetry")),
							Result: string(body),
						})
						if stop {
							s.LastRetryNum = 0
							return s.Result(in)
						}
					}
					resp, body, err = s.Result(in)
					return
				}
			}
		}
	}
	return
}
