package sgcache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	// 使用GetterFunc将匿名函数转换为接口Getter
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
