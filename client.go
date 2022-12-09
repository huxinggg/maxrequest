package maxrequest

import (
	"sync"
)

func Go(rs ...MaxRequestAttr) (results []GoResults) {
	results = make([]GoResults, len(rs))
	var wg sync.WaitGroup
	wg.Add(len(rs))
	for index, v := range rs {
		go func(vv MaxRequestAttr, i int) {
			var result GoResults
			result.Resp, result.Body, result.Err = vv.Result(nil)
			results[i] = result
			wg.Done()
		}(v, index)
	}
	wg.Wait()
	return
}
