package awesomeProject

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func TestIndexCreate(T *testing.T) {
	db_invert, _ := leveldb.OpenFile("path/to/db_invert", nil)
	defer db_invert.Close()
	db_mst, _ := leveldb.OpenFile("path/to/db_mst", nil)
	defer db_mst.Close()
	index := inverted_list{db:db_invert}
	index.renewList()
	fmt.Println(index)
	mst := MST{rootHash:[32]byte{116, 149, 10, 31, 146, 60, 8, 82, 112, 94, 241, 181, 65, 77, 232, 141, 75, 18, 212, 154, 187, 124, 114, 15, 28, 4, 186, 222, 235, 3, 71, 111}, db:db_mst}
	mst.reNewMst()
	mst.printMst()
	createIndex(file{Name:"tom",Keys: []string{"block","chain"}},88,&index,&mst)
	mst.printMst()
	fmt.Println(mst.search([]string{"篮球"}))
}