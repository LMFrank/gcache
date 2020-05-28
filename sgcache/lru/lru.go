package lru

import "container/list"

// LRU cache：由哈希表和双向链表组成
type Cache struct {
	maxBytes  int64 // 允许使用的最大内存
	nbytes    int64 // 当前已使用的内存
	ll        *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

// 节点的数据类型
type entry struct {
	key   string
	value Value
}

// 值可以为实现了Value接口的任意类型
type Value interface {
	Len() int // 返回值占用内存的大小
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// 寻找键值，先在字典内寻找key值，再寻找链表
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 淘汰最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// 新增或修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// key存在
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// key不存在
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 记录添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
