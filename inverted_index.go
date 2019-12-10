package awesomeProject

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
)

type key_file struct {
	Key       string
	File_list []string
}

type inverted_list struct {
	k    int
	list []key_file
	db *leveldb.DB
}

type key_file_list struct {
	k    int
	list []key_file
}

func fileToKey(file_name string, file_key []string) *key_file_list {
	key_file_list1 := key_file_list{}
	for i := 0; i < len(file_key); i++ {
		file_name_list := [] string{file_name}
		key_file1 := key_file{file_key[i], file_name_list}
		key_file_list1.list = append(key_file_list1.list, key_file1)
	}
	key_file_list1.k = len(key_file_list1.list)
	return &key_file_list1
}

func (index *inverted_list) swap(i, j int) *inverted_list {
	key_file1 := index.list[i]
	index.list[i] = index.list[j]
	index.list[j] = key_file1
	return index
}

func (index *inverted_list) insert(file key_file) *inverted_list {
	flag := 0
	var p int
	if len(index.list) == 0 {
		index.list = append(index.list, file)
	} else {
		for i := 0; i < len(index.list); i++ {
			if strings.Compare(file.Key, index.list[i].Key) == 0 {
				flag = 1
				for j := 0; j < len(file.File_list); j++ {
					index.list[i].File_list = append(index.list[i].File_list, file.File_list[j])
				}
				p = i
			}
		}
		if flag == 0 {
			index.list = append(index.list, file)
			p = len(index.list) - 1
		}
	}
	index.k = len(index.list)
	list_data,_ := rlp.EncodeToBytes(index.list[p].File_list)
	index.db.Put([]byte(file.Key), list_data, nil)
	return index
}

func (index *inverted_list) list_sort() *inverted_list {
	for j := index.k; j > 0; j--{
		for i := 0; i < j - 1; i++ {
			if len(index.list[i].File_list) < len(index.list[i+1].File_list){
				index.swap(i, i + 1)
			}
		}
	}
	return index
}

func (index *inverted_list) renewList() {
	iter := index.db.NewIterator(nil, nil)
	for iter.Next() {
		files := []string{}
		rlp.DecodeBytes(iter.Value(), &files)
		key_file1 := key_file{Key:string(iter.Key()), File_list: files}
		index.list = append(index.list, key_file1)
		index.k = len(index.list)
	}
	iter.Release()
	index.list_sort()
}

func (index *inverted_list) searchKey(key string) key_file {
	var keyfile1 key_file
	var files []string
	data1, _ := index.db.Get([]byte(key), nil)
	err :=rlp.DecodeBytes(data1, &files)
	if err != nil{
		fmt.Println(err)
	}
	keyfile1.File_list = files
	return keyfile1
}













