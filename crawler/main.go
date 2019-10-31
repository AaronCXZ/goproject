package main

import (
	"crawler/engine"
	"crawler/persist"
	"crawler/qichezhijia/parser"
	"crawler/scheduler"
)

func main() {
	itemChan, err := persist.ItemSaver()
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngins{
		Scheduler:   &scheduler.SimpleSchedule{},
		WorkerCount: 10,
		ItemChan:    itemChan,
	}

	e.Run(engine.Request{
		Url:        "https://k.autohome.com.cn/suva1",
		ParserFunc: parser.ParserBrandList,
	})
}
