package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// 一致性哈希算法的主数据结构
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashmap  map[int]string // 虚拟节点与真实节点的映射表，键为虚拟节点的哈希值，值为真实节点
}

func New(replicas int, fn hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashmap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 传入一个或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 每个节点创建m.replicas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashmap[hash] = key
		}
	}
	sort.Ints(m.keys) // 哈希值排序
}

// 通过key获得hash环上最近的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key))) // 计算key的hash值
	// 顺时针找到第一个匹配的虚拟节点的idx，从m.keys获得对应的hash值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 通过hashmap映射真实节点
	return m.hashmap[m.keys[idx%len(m.keys)]]
}
