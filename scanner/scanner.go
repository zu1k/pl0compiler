package scanner

import (
	"log"
	"github.com/zu1k/pl0complier/fileio"
	"github.com/zu1k/pl0complier/token"
	"unicode"
)

const (
	MAX_TOKEN_LENGTH = 10
)

type Scanner struct {
	file fileio.File
	line int
}

func (s *Scanner) Init(filepath string) {
	s.file = fileio.File{}
	s.file.Open(filepath)
	s.line = 1
}

type Symbol struct {
	SID   token.Token
	VALUE string
	NUM   int
}

var runecache []rune
var c = ' '
var end = false

func (s *Scanner) Scan() (sid token.Token, value string, num int) {
	runecache = make([]rune, 10)
	//去除空格、制表符、换行与回车等空白符
	for {
		if !isSpace(c) || end {
			break
		}
		s.getCh()
	}
	if end {
		sid = token.EOFSYM
		return
	}

	//出现字母，可能是保留字或者标识符
	if isLetter(c) {
		off := 0
		for {
			runecache[off] = c
			off++
			if off > MAX_TOKEN_LENGTH-1 {
				//TODO 标识符长度超过10
				log.Printf("ERROR: 标识符长度超限，等待解决的地方")
				break
			}
			s.getCh()
			if !isLetter(c) && !isDecimal(c) {
				break
			}
		}
		//现在是一个标识符或者关键字了，判断一下是什么
		value = string(runecache[:off])
		sid = token.Lookup(value)
		return
	} else
	//数字打头，只可能是一个数
	if isDecimal(c) {
		n := 0
		for {
			n = n*10 + int(c) - '0'
			s.getCh()
			if !isDecimal(c) {
				break
			}
		}
		//现在是一个数
		num = n
		sid = token.NUMBERSYM
		return
	} else
	//可能会成对出现的符号
	if maybeDouble(c) {
		var opt = [2]rune{c}
		var opts string
		s.getCh()
		if c == '=' {
			s.getCh()
			opt[1] = '='
			opts = string(opt[:])
		} else {
			opts = string(opt[:1])
		}
		sid = token.OptParse(opts)
		value = opts
		return
	} else {
		//一些可能出现的孤立的符号
		switch c {
		case '=', '+', '-', '*', '/', '#', '(', ')', ';', '.', ',':
			sid = token.OptParse(string(c))
			value = string(c)
		default:
			//TODO 这里是出现了不应该出现的字符，应该处理错误
			log.Printf("这里是出现了不应该出现的字符: %s", string(c))
		}
		s.getCh()
		return
	}
}

/*
 * 读取一个字符
 */
func (s *Scanner) getCh() {
	c, end = s.file.ReadRune()
	if c == '\n' {
		s.line++
	}
	//log.Println(string(c))
}

/*
 * 判断是不是空白符
 */
func isSpace(ch rune) (spcae bool) {
	return unicode.IsSpace(ch)
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' && unicode.IsLetter(ch)
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }

func parseNumber(intrune []rune) (n int) {
	return
}

func maybeDouble(ch rune) bool {
	return ch == '<' || ch == '>' || ch == ':'
}
