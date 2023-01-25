package main

import (
	"my_item_cf_go/plugin/item_cf_big_data/cf_lib"
)

func main() {
	cf := cf_lib.GetItemCF()
	cf.DoCalculate()
	//训练完成后直接显示评估结果
	cf.EvaluateData()
}
