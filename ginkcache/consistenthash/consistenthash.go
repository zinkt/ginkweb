package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	// hash函数
	hashFunc Hash
	// 虚拟节点倍数
	replicas int
	// hash环（已排序）
	keys []int
	// 虚拟节点与真实节点的映射表：map[虚拟节点的哈希值]真实节点的名称
	hashMap map[int]string
}

// 创建新Map，允许自定义虚拟节点倍数和哈希函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hashFunc: fn,
		hashMap:  make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// 向Map中新增节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对每个真实节点，创建m.replicas个虚拟节点，并使其hash在字典中映射到真实节点的名字
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hashFunc([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 根据key从Map中获取最近的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))
	// binary search
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
