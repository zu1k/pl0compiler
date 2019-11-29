/**
 * @auther:  zu1k
 * @date:    2019/11/2
 */
package token

import (
	"strconv"
	"unicode"
)

// Token is the set of lexical tokens of the Go programming language.
type Token int

/*
 * Sym类型的枚举类型
 * The list of tokens.
 */
const (
	BADTOKEN = iota // 无效的单词

	literal_beg
	IDENTSYM  // 标识符
	NUMBERSYM // 数
	literal_end

	operator_beg
	// Operators
	PLUSSYM  // +
	MINUSYM  // -
	MULSYM   // *
	SLASHSYM // /

	relation_optr_beg
	EQLSYM // =
	NEQSYM // #
	LESSYM // <
	LEQSYM // <=
	GTRSYM // >
	GEQSYM // >=
	relation_optr_end

	LPARENTSYM   // (
	RPARENTSYM   // )
	COMMASYM     // ,
	SEMICOLOMSYM // ;
	PERIODSYM    // .
	BECOMESSYM   // :=
	operator_end

	keyword_beg
	// Keywords
	BEGINSYM // begin
	ENDSYM   // end
	IFSYM    // if
	ELSESYM  // else
	THENSYM  // then
	WHILESYM // while
	DOSYM    // do
	CALLSYM  // call
	CONSTSYM // const
	VARSYM   // var
	PROCSYM  // procedure
	ODDSYM   // odd
	WRITESYM // write
	READSYM  // read
	keyword_end

	EOFSYM // EOF
)

var tokens = [...]string{
	CONSTSYM: "const",
	VARSYM:   "var",
	PROCSYM:  "procedure",
	CALLSYM:  "call",
	BEGINSYM: "begin",
	ENDSYM:   "end",
	IFSYM:    "if",
	THENSYM:  "then",
	ELSESYM:  "else",
	WHILESYM: "while",
	DOSYM:    "do",
	READSYM:  "read",
	WRITESYM: "write",
	ODDSYM:   "odd",

	PLUSSYM:    "+",
	MINUSYM:    "-",
	MULSYM:     "*",
	SLASHSYM:   "/",
	LPARENTSYM: "(",
	RPARENTSYM: ")",
	COMMASYM:   ",",
	PERIODSYM:  ".",

	EQLSYM: "=",
	LESSYM: "<",
	GTRSYM: ">",
	NEQSYM: "#",
	LEQSYM: "<=",
	GEQSYM: ">=",

	SEMICOLOMSYM: ";",
	BECOMESSYM:   ":=",

	BADTOKEN:  "无效字符",
	NUMBERSYM: "数字",
	IDENTSYM:  "变量标识符",

	EOFSYM: "文档已结束",
}

var Tokens = tokens

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = Token(i) //TODO ?what?
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENTSYM
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

func (tok Token) IsIdent() bool     { return tok == IDENTSYM }
func (tok Token) IsNumber() bool    { return tok == NUMBERSYM }
func (tok Token) IsBecome() bool    { return tok == BECOMESSYM }
func (tok Token) IsCall() bool      { return tok == CALLSYM }
func (tok Token) IsThen() bool      { return tok == THENSYM }
func (tok Token) IsSemicolom() bool { return tok == SEMICOLOMSYM }
func (tok Token) IsEnd() bool       { return tok == ENDSYM }
func (tok Token) IsDo() bool        { return tok == DOSYM }
func (tok Token) IsConst() bool     { return tok == CONSTSYM }
func (tok Token) IsComma() bool     { return tok == COMMASYM }
func (tok Token) IsVar() bool       { return tok == VARSYM }
func (tok Token) IsProc() bool      { return tok == PROCSYM }
func (tok Token) IsLparent() bool   { return tok == LPARENTSYM }
func (tok Token) IsRparent() bool   { return tok == RPARENTSYM }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

func (tok Token) IsRelationOpr() bool { return relation_optr_beg < tok && tok < relation_optr_end }

// IsKeyword reports whether name is a PL/0 keyword, such as "const" or "read".
//
func IsKeyword(name string) bool {
	// TODO: opt: use a perfect hash function instead of a global map.
	_, ok := keywords[name]
	return ok
}

// IsIdentifier reports whether name is a Go identifier, that is, a non-empty
// string made up of letters, digits, and underscores, where the first character
// is not a digit. Keywords are not identifiers.
//
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}

/*
 * PL/0 保留字 - type
 */
var ReversedWordMap = map[string]int{
	"const":     CONSTSYM,
	"var":       VARSYM,
	"procedure": PROCSYM,
	"call":      CALLSYM,
	"begin":     BEGINSYM,
	"end":       ENDSYM,
	"if":        IFSYM,
	"then":      THENSYM,
	"else":      ELSESYM,
	"while":     WHILESYM,
	"do":        DOSYM,
	"read":      READSYM,
	"write":     WRITESYM,
	"odd":       ODDSYM,
}

/*
 * 特殊字符列表
 */
var SpecialSymbol = [...]string{"+", "-", "*", "/", "(", ")", "=", ",", ".", "<", ">", ";", ":"}

/*
 * 特殊符号 - type
 */
var SpecialSymbolMap = map[string]Token{
	"+":  PLUSSYM,
	"-":  MINUSYM,
	"*":  MULSYM,
	"/":  SLASHSYM,
	"(":  LPARENTSYM,
	")":  RPARENTSYM,
	"=":  EQLSYM,
	",":  COMMASYM,
	".":  PERIODSYM,
	"<":  LESSYM,
	">":  GTRSYM,
	"#":  NEQSYM,
	"<=": LEQSYM,
	">=": GEQSYM,
	";":  SEMICOLOMSYM,
	":=": BECOMESSYM,
}

func OptParse(opt string) Token {
	if tok, is_opt := SpecialSymbolMap[opt]; is_opt {
		return tok
	}
	return BADTOKEN
}
