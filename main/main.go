package main

import (
	"fmt"
	"github.com/huxinggg/maxrequest"
)

type ResourceData struct {
	CanvasData []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"canvas_data"`
	Materials struct {
		EffectImg []struct {
			Code string `json:"code"`
			Img  string `json:"img"`
		} `json:"effect_img"`
		CutpartImg []struct {
			Code string `json:"code"`
			Img  string `json:"img"`
		} `json:"cutpart_img"`
	} `json:"materials"`
}

func main() {
	var r ResourceData
	maxrequest.New().Get("https://shuibaba.oss-cn-beijing.aliyuncs.com/files/5bc681871c5376fe8bf94e39a878cdd6.json").Result(&r)
	fmt.Println(r)

}
