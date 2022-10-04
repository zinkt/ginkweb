package lru

import "container/list"

// last recent use cache
// 并发访问不安全
type Cache struct {
	// 允许使用的最大内存
	maxBytes int64
	// 当前已使用的内存
	nbytes int64

	ll *list.List

	// 值是双向链表中的结点指针
	cache map[string]*list.Element

	// 某条记录被移除时的回调函数，可为nil
	OnEvicted func(key string, value Value)
}

// 双向链表节点的数据类型
// 保存key是为了淘汰队首节点时，删除字典中的映射
type entry struct {
	key   string
	value Value
}

// 值需要实现此接口
// Len()返回占用多少bytes
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     map[string]*list.Element{},
		OnEvicted: onEvicted,
	}
}

// (约定Front为队尾)
// 查找时自动将节点移动至队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 缓存淘汰
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果键已存在，则更新值，更新占用空间，并将其移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 键不存在，则插入队尾
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 限制内存占用
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//  方便测试
func (c *Cache) Len() int {
	return c.ll.Len()
}
