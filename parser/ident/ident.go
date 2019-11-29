/**
 * @auther:  zu1k
 * @date:    2019/11/9
 */
package ident

type Ident int

const (
	Constant = iota
	Variable
	Proc
)

var identtable = [...]string{
	Constant: "常量",
	Variable: "变量",
	Proc:     "过程",
}

func (i Ident) String() string {
	return identtable[i]
}

func (id Ident) IsVariable() bool  { return id == Variable }
func (id Ident) IsConsttant() bool { return id == Constant }
func (id Ident) IsProcedure() bool { return id == Proc }
