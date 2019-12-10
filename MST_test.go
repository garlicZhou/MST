package awesomeProject

import (
	"fmt"
	"testing"
)


func TestMstInsert(t *testing.T) {
	mst := new()
	mst.root_insert(index_info{key: []string{"篮球"},pos:10})
	mst.root_insert(index_info{key: []string{"篮球","网球"},pos:30})
	mst.root_insert(index_info{key: []string{"篮球","羽毛球","乒乓球"},pos:40})
	mst.root_insert(index_info{key: []string{"足球","网球"},pos:20})
	mst.root_insert(index_info{key: []string{"游泳","潜泳","网球"},pos:50})
	mst.root_insert(index_info{key: []string{"游泳","潜泳","台球"},pos:60})
	mst.root_insert(index_info{key: []string{"篮球","网球"},pos:70})
	mst.root_insert(index_info{key: []string{"足球","排球"},pos:80})
	mst.printMst()
	fmt.Println("查询")
	keys1 := []string{"羽毛球"}
	fmt.Println(mst.search(keys1))


	//fmt.Println(*mst.root)
}
//
//func TestSearch(t *testing.T) {
//	keys1 := []string{"篮球"}
//	fmt.Println()
//
//}