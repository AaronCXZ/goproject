package parser

import (
	"io/ioutil"
	"testing"
)

func TestParserBrandList(t *testing.T) {
	contents, err := ioutil.ReadFile("parser_test_data.html")
	if err != nil {
		panic(err)
	}
	result := ParserBrandList(contents)
	const resultSize = 138
	expectedUrls := []string{
		"https://k.autohome.com.cn/314/", "https://k.autohome.com.cn/3386/", "https://k.autohome.com.cn/2123/",
	}
	expectedBrands := []string{
		"本田CR-V", "宝马X2", "哈弗H6",
	}
	if len(result.Requests) != resultSize {
		t.Errorf("%d, %d", resultSize, len(result.Requests))
	}
	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("%d: %s;but %s", i, url, result.Requests[i].Url)
		}
	}

	if len(result.Items) != resultSize {
		t.Errorf("%d, %d", resultSize, len(result.Items))
	}

	for i, brand := range expectedBrands {
		if result.Items[i].Payload.(string) != brand {
			t.Errorf("%d: %s;but %s", i, brand, result.Items[i])
		}
	}
}
