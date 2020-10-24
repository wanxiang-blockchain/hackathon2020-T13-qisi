package main

import (
	"fisco/check"
	"flag"
)

func main() {
	module := flag.String("module", "SettledChain", "选择测试的模块, 默认为test")

	flag.Parse()

	//if *module == "test" {
	//	check.Test()
	//}

	//if *module == "atom" {
	//	check.Atom()
	//}

	if *module == "SettledChain" {
		check.Settled()
	}
}
