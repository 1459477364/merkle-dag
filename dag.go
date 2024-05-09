package merkledag

import (
	"crypto/sha256"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, hp HashPool) []byte {
	if node.Type() == FILE {
		// 处理文件节点
		file := node.(File)
		contentHash := sha256.Sum256(file.Bytes())
		store.Put(contentHash[:], file.Bytes())
		return contentHash[:]
	} else {
		// 处理目录节点
		dir := node.(Dir)
		it := dir.It()
		var childHashes []byte
		for it.Next() {
			childNode := it.Node()
			childHash := Add(store, childNode, hp)
			childHashes = append(childHashes, childHash...)
		}
		// 计算目录的Merkle Root
		dirHash := sha256.Sum256(childHashes)
		store.Put(dirHash[:], childHashes)
		return dirHash[:]
	}
}
