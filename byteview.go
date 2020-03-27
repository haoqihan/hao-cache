package hao_cache

type ByteView struct {
	b []byte
}

// 返回其所占的内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// 防止缓存值被外部程序修改。
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 返回 string 类型
func (v ByteView) String() string {
	return string(v.b)

}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
