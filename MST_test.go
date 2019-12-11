package awesomeProject

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)


func TestMstInsert(t *testing.T) {
	mst := new()
	db, _ := leveldb.OpenFile("path/to/db_mst", nil)
	defer db.Close()
	mst.putDb(db)
	mst.root_insert(index_info{key: []string{"篮球"},pos:10})
	mst.root_insert(index_info{key: []string{"篮球","网球"},pos:30})
	mst.root_insert(index_info{key: []string{"篮球","羽毛球","乒乓球"},pos:40})
	mst.root_insert(index_info{key: []string{"足球","网球"},pos:20})
	mst.root_insert(index_info{key: []string{"游泳","潜泳","网球"},pos:50})
	mst.root_insert(index_info{key: []string{"游泳","潜泳","台球"},pos:60})
	mst.root_insert(index_info{key: []string{"篮球","网球"},pos:70})
	mst.root_insert(index_info{key: []string{"足球","排球"},pos:80})
	mst.putRootHash()
	mst.printMst()
	fmt.Println("查询")
	keys1 := []string{"羽毛球"}
	fmt.Println(mst.search(keys1))

	mst2 := MST{rootHash:mst.rootHash, db:mst.db}
	mst2.reNewMst()
	mst2.printMst()
}
//
func TestSearchMst(t *testing.T) {
	db, _ := leveldb.OpenFile("path/to/db_mst", nil)
	defer db.Close()
	nodekv1 := nodekv{}
	data1, _ := db.Get([]byte{249, 245, 89, 133, 74, 53, 240, 79, 54, 64, 188, 241, 100, 133, 182, 208, 104, 3, 88, 204, 130, 226, 202, 51, 68, 88, 148, 70, 214, 145, 129, 79}, nil)
	err :=rlp.DecodeBytes(data1, &nodekv1)
	if err != nil{
		fmt.Println(err)
	}
	//mst := MST{db}
	//mst.printMst()

}