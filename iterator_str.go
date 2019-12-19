package jzon

import (
	"unicode/utf16"
)

const (
	noEscape   = 0
	invalidHex = -1
)

var (
	escapeMap [charNum]byte
	hexValue  [charNum]int8
)

func init() {
	// escaped characters
	for i := 0; i < charNum; i++ {
		escapeMap[i] = noEscape
	}
	for k, v := range map[byte]byte{
		'"':  '"',
		'\\': '\\',
		'/':  '/',
		'b':  '\b',
		'f':  '\f',
		'n':  '\n',
		'r':  '\r',
		't':  '\t',
	} {
		escapeMap[k] = v
	}
	// hex values
	for i := 0; i < charNum; i++ {
		hexValue[i] = invalidHex
	}
	for c := '0'; c <= '9'; c++ {
		hexValue[c] = int8(c - '0')
	}
	for c := 'a'; c <= 'f'; c++ {
		hexValue[c] = int8(c - 'a' + 10)
	}
	for c := 'A'; c <= 'F'; c++ {
		hexValue[c] = int8(c - 'A' + 10)
	}
}

func (it *Iterator) readU4() (ret rune, err error) {
	remain := 4
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c := it.buffer[i]
			u4v := hexValue[c]
			if u4v == invalidHex {
				return 0, InvalidUnicodeCharError{c: c}
			}
			ret = ret<<4 + int32(u4v)
			if remain == 1 {
				it.head = i + 1
				return
			}
			remain--
		}
		it.head = i
		if err = it.readMore(); err != nil {
			return
		}
	}
}

func (it *Iterator) readEscapedChar(b []byte) ([]byte, error) {
	c, err := it.nextByte()
	if err != nil {
		return b, err
	}
	escaped := escapeMap[c]
	if escaped != noEscape {
		it.head += 1
		return append(b, escaped), nil
	}
	if c != 'u' {
		return b, InvalidEscapeCharError{c: c}
	}
	it.head += 1
	r, err := it.readU4()
	if err != nil {
		return b, err
	}
	if utf16.IsSurrogate(r) {
		c, err := it.nextByte()
		if err != nil {
			return b, err
		}
		if c != '\\' {
			return appendRune(b, r), nil
		}
		it.head += 1
		c, err = it.nextByte()
		if err != nil {
			return b, err
		}
		if c != 'u' {
			b = appendRune(b, r)
			escaped := escapeMap[c]
			if escaped == noEscape {
				return b, InvalidEscapeCharError{c: c}
			}
			it.head += 1
			return append(b, escaped), nil
		}
		it.head += 1
		r2, err := it.readU4()
		if err != nil {
			return b, err
		}
		combined := utf16.DecodeRune(r, r2)
		if combined == runeError {
			b = appendRune(b, r)
			return appendRune(b, r2), nil
		}
		return appendRune(b, combined), nil
	} else {
		return appendRune(b, r), nil
	}
}

// internal, call only after a '"' is consumed
func (it *Iterator) readStringAsSlice(buf []byte) (_ []byte, err error) {
	for {
		i := it.head
		for i < it.tail {
			c := it.buffer[i]
			if c == '"' {
				buf = append(buf, it.buffer[it.head:i]...)
				it.head = i + 1
				return buf, nil
			} else if c == '\\' {
				buf = append(buf, it.buffer[it.head:i]...)
				it.head = i + 1
				buf, err = it.readEscapedChar(buf)
				if err != nil {
					return buf, err
				}
				i = it.head
			} else if c < ' ' { // json.org
				return buf, InvalidStringCharError{c: c}
			} else {
				i++
			}
		}
		// i == it.tail
		buf = append(buf, it.buffer[it.head:i]...)
		it.head = i
		if err = it.readMore(); err != nil {
			return buf, err
		}
	}
}

func (it *Iterator) expectQuote() error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '"' {
		return UnexpectedByteError{exp: '"', got: c}
	}
	it.head += 1 // consume the leading '"'
	return nil
}

func (it *Iterator) ReadStringAsSlice(buf []byte) (_ []byte, err error) {
	if err = it.expectQuote(); err != nil {
		return
	}
	return it.readStringAsSlice(buf)
}

// internal, call only after a '"' is consumed
func (it *Iterator) readString() (ret string, err error) {
	buf, err := it.readStringAsSlice(it.tmpBuffer[:0])
	it.tmpBuffer = buf
	if err == nil {
		ret = string(buf)
	}
	return
}

func (it *Iterator) ReadString() (ret string, err error) {
	if err = it.expectQuote(); err != nil {
		return
	}
	return it.readString()
}

// From jsoniter
const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	maxRune   = '\U0010FFFF' // Maximum valid Unicode code point.
	runeError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
)

func appendRune(p []byte, r rune) []byte {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		p = append(p, byte(r))
		return p
	case i <= rune2Max:
		p = append(p, t2|byte(r>>6))
		p = append(p, tx|byte(r)&maskx)
		return p
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = runeError
		fallthrough
	case i <= rune3Max:
		p = append(p, t3|byte(r>>12))
		p = append(p, tx|byte(r>>6)&maskx)
		p = append(p, tx|byte(r)&maskx)
		return p
	default:
		p = append(p, t4|byte(r>>18))
		p = append(p, tx|byte(r>>12)&maskx)
		p = append(p, tx|byte(r>>6)&maskx)
		p = append(p, tx|byte(r)&maskx)
		return p
	}
}
