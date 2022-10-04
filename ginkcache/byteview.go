package ginkcache

// 缓存值的抽象，保存了不可更改的bytes
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

// 返回一份已复制的数据的byte切片
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 以string返回数据
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
