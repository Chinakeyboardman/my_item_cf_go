package main

import "my_item_cf_go/plugin/item_cf/cf_lib"

func main() {
	cf := cf_lib.GetItemCF()
	cf.DoEvaluate()
}
