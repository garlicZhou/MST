package awesomeProject

type file struct {
	Name string
	Keys []string
}

func createIndex(f1 file, blockNumber uint) *MST {
	key_file_list1 := fileToKey(f1.Name, f1.Keys)
	index := inverted_list{}
	MST1 := MST{}
   for i := range key_file_list1.list {
   	index.insert(key_file_list1.list[i])
	}
   index.list_sort()
   MST1.root_insert(index_info{key:f1.Keys, pos:blockNumber})
   return &MST1
}

func (file1 *file) keysSort() {


}

func (file *file) swap()   {

}

