package awesomeProject

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
)

const MAXLEVEL = 3

type node struct {
	parent    *node
	child     []*node
	childHash [][32]byte
	key       []string
	value     []uint
	hash      [32]byte
	isLeaf    bool
	isExtend  bool
}

type nodekv struct {
	ChildHash [][32]byte
	Key       []string
	Value     []uint
	Hash      [32]byte
	IsLeaf    bool
	IsExtend  bool
}

type index_info struct {
	key []string
	pos uint
}

type MST struct {
	root *node
	db   *leveldb.DB
}

func new() *MST {
	mst := &MST{root: &node{isLeaf: true, isExtend: false}}
	return mst
}

func (t *MST) root_insert(in index_info) {
	li := len(in.key)
	for i := 0; i < len(in.key); i++ {
		in.key = in.key[i:li:li]
		root1 := t.root
		flag := true
		for r := 0; r < len(root1.child); r++ {
			if strings.Compare(in.key[0], root1.child[r].key[0]) == 0 {
				root1.isLeaf = false
				flag = false
				root1.child[r].insert(in)
				break
			}
		}
		if flag {
			root1.child = append(root1.child, &node{parent: root1, key: in.key, value: []uint{in.pos}, isLeaf: true, isExtend: false})
			root1.isLeaf = false
			if len(root1.child) == 1 {
				root1.isExtend = true
			} else {
				root1.isExtend = false
			}
			j := len(root1.child)
			root1.child[j-1].updateHash()
		}
	}
}

func (node1 *node) insert(in index_info) {
	ln := len(node1.key)
	li := len(in.key)
	flag := true
	if li == ln {
		for k := 0; k < ln; k++ {
			if strings.Compare(node1.key[k], in.key[k]) == 0 {
				if k == ln-1 {
					if node1.isLeaf || node1.isExtend {
						node1.value = append(node1.value, in.pos)
						node1.updateHash()
					} else {
						node1.isExtend = true
						node1.value = append(node1.value, in.pos)
						node1.updateHash()
					}
				} else {
					continue
				}
			} else {
				node1.split(in, k, ln, li)
				break
			}
		}
	} else if li > ln {
		if ln == 0 {
			for r := 0; r < len(node1.child); r++ {
				if strings.Compare(in.key[0], node1.child[r].key[0]) == 0 {
					node1.isLeaf = false
					flag = false
					node1.child[r].insert(in)
					break
				}
			}
			if flag {
				node1.child = append(node1.child, &node{parent: node1, key: in.key[ln:li:li], value: []uint{in.pos}, isLeaf: true, isExtend: false})
				node1.child[len(node1.child)-1].updateHash()
			}
			return
		}
		for k := 0; k < ln; k++ {
			if strings.Compare(node1.key[k], in.key[k]) == 0 {
				if k == ln-1 {
					if node1.isLeaf {
						node1.isExtend = true
						node1.isLeaf = false
						node1.child = append(node1.child, &node{parent: node1, isLeaf: true, key: in.key[k+1 : li : li], value: []uint{in.pos}})
						node1.child[len(node1.child)-1].updateHash()
					} else {
						for r := 0; r < len(node1.child); r++ {
							if strings.Compare(in.key[ln], node1.child[r].key[0]) == 0 {
								node1.child[r].isExtend = false
								node1.child[r].insert(index_info{key: in.key[ln:li:li], pos: in.pos})
								flag = false
								break
							}
						}
						if flag {
							node1.isExtend = false
							node1.child = append(node1.child, &node{parent: node1, key: in.key[ln:li:li], value: []uint{in.pos}, isLeaf: true, isExtend: false})
							node1.child[len(node1.child)-1].updateHash()
						}
					}
				} else {
					continue
				}
			} else {
				node1.split(in, k, ln, li)
				break
			}
		}

	} else if li < ln {
		for k := 0; k < li; k++ {
			if strings.Compare(node1.key[k], in.key[k]) == 0 {
				if k == li-1 {
					node2 := node{parent: node1, child: node1.child, key: node1.key[li:ln:ln], value: node1.value, isLeaf: node1.isLeaf, isExtend: node1.isExtend}
					node1.isLeaf = false
					node1.isExtend = true
					node1.value = nil
					node1.value = append(node1.value, in.pos)
					node1.key = node1.key[0:li:li]
					node1.child = nil
					node1.child = append(node1.child, &node2)
					node1.child[len(node1.child)-1].updateHash()
				}
			} else {
				node1.split(in, k, ln, li)
				break
			}
		}
	}
}

func (node1 *node) split(in index_info, k int, ln int, li int) {
	node1.child = append(node1.child, &node{parent: node1, isLeaf: true, key: node1.key[k:ln:ln], value: node1.value})
	node1.child[len(node1.child)-1].updateHash()
	node1.child = append(node1.child, &node{parent: node1, isLeaf: true, key: in.key[k:li:li], value: []uint{in.pos}})
	node1.child[len(node1.child)-1].updateHash()
	node1.isLeaf = false
	node1.key = node1.key[0:k:ln]
	node1.value = nil
}

func (node1 *node) updateHash() {
	if len(node1.child) > 0 {
		node1.childHash = nil
		for k := range node1.child {
			node1.childHash =  append(node1.childHash, node1.child[k].hash)
		}
	}
	nodekv1 := nodekv{node1.childHash, node1.key, node1.value, node1.hash, node1.isLeaf, node1.isExtend}
	//fmt.Println(nodekv1)
	var data []byte
	if nodekv1.IsLeaf {
		data, _ = rlp.EncodeToBytes(nodekv1)
		//fmt.Println("data:", data)
		nodekv1.Hash = sha256.Sum256(data)
		//fmt.Println("hash:", nodekv1.Hash)
		node1.hash = nodekv1.Hash
		if node1.parent == nil {
			return
		} else {
			node1.parent.updateHash()
		}
	} else if nodekv1.IsExtend {
		data, _ = rlp.EncodeToBytes(nodekv1)
		for i := range nodekv1.ChildHash {
			for j := range nodekv1.ChildHash[i] {
				data = append(data, nodekv1.ChildHash[i][j])
			}
		}
		nodekv1.Hash = sha256.Sum256(data)
		node1.hash = nodekv1.Hash
		if node1.parent == nil {
			return
		} else {
			node1.parent.updateHash()
		}
	} else {
		for i := range nodekv1.ChildHash {
			for j := range nodekv1.ChildHash[i] {
				data = append(data, nodekv1.ChildHash[i][j])
			}
		}
		nodekv1.Hash = sha256.Sum256(data)
		node1.hash = nodekv1.Hash
		if node1.parent == nil {
			return
		} else {
			node1.parent.updateHash()
		}
	}
}

func (t *MST) printMst() {
	fmt.Printf("root: ")
	t.root.printNode()
	fmt.Print("\n")
	for i := 0; i < len(t.root.child); i++ {
		t.root.child[i].printNode()
	}
	fmt.Print("\n")
	for i := 0; i < len(t.root.child); i++ {
		for j := 0; j < len(t.root.child[i].child); j++ {
			t.root.child[i].child[j].printNode()
		}
		fmt.Print("  ")
	}
	fmt.Print("\n")
}

func (node1 *node) printNode() {
	fmt.Printf("parent: %p", node1.parent)
	fmt.Print(" ", "keys:", node1.key, " ", "value:", node1.value, " ", "isLeaf: ", node1.isLeaf, " ", "isExtend: ", node1.isExtend, " ", "hash: ", node1.hash, " ")
}

func (node1 *node) searchNode(keys []string) []uint {
	lk := len(keys)
	ln := len(node1.key)
	if node1.isLeaf {
		if lk > ln {
			return nil
		}
		for i := 0; i < ln; i++ {
			if keys[i] == node1.key[i] {
				if i == lk-1 {
					return node1.value
				}
			} else {
				return nil
			}
		}
	} else {
		for i := 0; i < ln; i++ {
			if keys[i] == node1.key[i] {
				if i == lk-1 {
					value_result := node1.value
					for j := 0; j < len(node1.child); j++ {
						for k := 0; k < len(node1.child[j].value); k++ {
							value_result = append(value_result, node1.child[j].value[k])
						}
					}
					return value_result
				} else if i == ln-1 {
					for j := 0; j < len(node1.child); j++ {
						if node1.child[j].key[0] == keys[i+1] {
							return node1.child[j].searchNode(keys[i+1 : lk : lk])
						}
					}
				}
			}
		}
	}
	return nil
}

func (t *MST) search(keys []string) []uint {
	for i := 0; i < len(t.root.child); i++ {
		if t.root.child[i].key[0] == keys[0] {
			return t.root.child[i].searchNode(keys)
		}
	}
	return nil
}

//func (node1 *node) reNewNode() {
//	for i := 0;i < node.
//
//}
