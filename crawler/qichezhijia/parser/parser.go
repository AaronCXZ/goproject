package parser

import (
	"crawler/engine"
	"regexp"
)

var (
	cityListRe = regexp.MustCompile(`<div class="cont-name"><a  href="(/[0-9]+/)" target="_blank" class="font-14-b">([^<]+)</a></div>`)
)

func ParserBrandList(contents []byte) engine.ParseResult {
	match := cityListRe.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, m := range match {
		result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, engine.Request{Url: "https://k.autohome.com.cn" + string(m[1]), ParserFunc: ParseProfile})
	}
	return result
}
