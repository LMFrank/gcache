package sgcache

import (
	"fmt"
	"log"
	"sync"
)

// 缓存的命名空间
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据的回调
	mainCache cache  // 并发缓存
}

type Getter interface {
	Get(key string) ([]byte, error)
}

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type GetterFunc func(key string) ([]byte, error)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter GetterFunc) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 从mainCache中查找缓存，如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[SgCache] hit")
		return v, nil
	}

	// 缓存不存在，则调用load方法
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用用户的回调方法获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 将源数据添加到缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
