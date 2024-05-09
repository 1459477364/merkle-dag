package merkledag

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
)

// serializeNode 将Node对象序列化为字节数组
func serializeNode(node Node) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(node)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// deserializeNode 将字节数组反序列化为Node对象
func deserializeNode(data []byte) (Node, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var node Node
	err := decoder.Decode(&node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) ([]byte, error) {
	// 从KVStore中获取根节点数据
	data, err := store.Get(hash)
	if err != nil {
		return nil, err
	}

	// 反序列化根节点
	node, err := deserializeNode(data)
	if err != nil {
		return nil, err
	}

	// 根据路径查找特定文件
	pathSegments := splitPath(path)
	for _, segment := range pathSegments {
		if node.Type() == DIR {
			found := false
			dir := node.(Dir)
			it := dir.It()
			for it.Next() {
				child := it.Node()
				if child.Name() == segment {
					node = child
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("path segment not found: " + segment)
			}
		} else {
			return nil, errors.New("expected directory but found file at segment: " + segment)
		}
	}

	if node.Type() == FILE {
		file := node.(File)
		return file.Bytes(), nil
	}

	return nil, errors.New("path does not point to a file")
}

func splitPath(path string) []string {
	// Assume path is separated by "/"
	return strings.Split(path, "/")
}
