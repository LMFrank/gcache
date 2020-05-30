package sgcache

// 根据传入的key选择相应的节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 从对应的group查找缓存值，对应Http客户端
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
