package bufferstring

//	"fmt"
//	"unicode/utf8"

const (
	BUFFER_SIZE = 255 // размер буфера по умолчанию
)

// создать буфер buffSize объема. выделение дополнительного объема будет производиться в объеме addSpace
func New(buffSize int, addSpace int) *BufferString {
	this := new(BufferString)
	this.init(buffSize, addSpace)
	return this
}

// создать буфер на основании строки. выделение дополнительного объема будет производиться в объеме addSpace
func NewFromString(s string, addSpace int) *BufferString {
	this := new(BufferString)
	this.initFromString(s, addSpace)
	return this
}
