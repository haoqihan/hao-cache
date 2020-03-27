package lru

import (
	"container/list"
	"fmt"
)

type Value interface {
	Len() int
}

type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存，
	nBytes    int64                         // nBytes 当前已使用的内存
	ll        *list.List                    // 使用GO语言标准库实现的双向链表
	cache     map[string]*list.Element      // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为nil
}

// entry 是双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// New 实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找功能
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 从字典中找到对应的双向链表节点
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // 将该节点移动到队尾
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 取到队首节点，从链表中删除
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从字典中 c.cache 删除该节点的映射关系
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新当前所用的内存 c.nBytes
		if c.OnEvicted != nil {                                // 如果回调函数 OnEvicted 不为nil，则调用回调函数
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加或修改
func (c *Cache) Add(key string, value Value) {
	// 如果键存在，则更新对应节点的值，并将该节点移到队尾
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 如果不存在 则是新增，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 更新c.nBytes 如果超过了设定的最大值 c.maxBytes则移除最少访问节点
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len 获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
func (c Cache)Tree()  {
	fmt.Println("--------------------------------------------------")
	for k,v := range c.cache{
		fmt.Println(k,v)
	}
}
