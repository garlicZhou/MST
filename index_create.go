package awesomeProject

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

type file struct {
	Name string
	Keys []string
}

func createIndex(f1 file, blockNumber uint,in *inverted_list,t *MST) *MST {
	key_file_list1 := fileToKey(f1.Name, f1.Keys)
	index := inverted_list{db:in.db}
	for i := range key_file_list1.list {
		index.insert(key_file_list1.list[i])
	}
	index.list_sort()
	f1.keysSort(in.db)
	t.root_insert(index_info{key: f1.Keys, pos: blockNumber})
	return t
}

func (file1 *file) keysSort(db *leveldb.DB) {
	var key_file_pre key_file
	var key_file_next key_file
	for j := len(file1.Keys); j > 0; j-- {
		for i := 0;i < j - 1;i++ {
			data, _ := db.Get([]byte(file1.Keys[i]), nil)
            key_file_pre.Key = file1.Keys[i]
			rlp.DecodeBytes(data,&key_file_pre.File_list)
			key_file_next.Key = file1.Keys[i]
			rlp.DecodeBytes(data,&key_file_next.File_list)
            if len(key_file_pre.Key) < len(key_file_next.Key) {
            	file1.swap(i, i + 1)
			}
		}
	}
}

func (file *file) swap(i, j int) {
	flag := file.Keys[i]
	file.Keys[j] = file.Keys[i]
	file.Keys[i] = flag
}
