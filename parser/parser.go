/**
 * @auther:  zu1k
 * @date:    2019/11/8
 */
package parser

import (
	"log"
	"pl0compiler/asm"
	"pl0compiler/interpret"
	"pl0compiler/parser/fct"
	"pl0compiler/parser/ident"
	"pl0compiler/parser/selectset"
	"pl0compiler/scanner"
	"pl0compiler/token"
)

type Parser struct {
	scanner scanner.Scanner
	// Next token
	tok token.Token
	lit string
	num int

	//嵌套层数
	lev int

	//中间代码
	codes []asm.Asm
	cidx  int

	//符号表
	table []identitem
	tx    int
	id    string //最近一次识别出来的ident

	dx int
}

func (p *Parser) Interpret() {
	interpret.Interpret(p.codes)
}

func (p *Parser) Scan() {
	p.tok, p.lit, p.num = p.scanner.Scan()
	if p.tok.IsIdent() {
		p.id = p.lit
	}
}

func (p *Parser) Init(filepath string) {
	p.scanner = scanner.Scanner{}
	p.scanner.Init(filepath)
	p.Scan()
	p.lev = 0
	p.codes = make([]asm.Asm, 0)
	p.cidx = 0
	p.table = make([]identitem, 100)
	p.tx = 0
	p.dx = 0
}

const (
	MAX_LEV = 3
)

//符号表里的符号
type identitem struct {
	name  string
	kind  ident.Ident
	lev   int
	value int //变量，且只可能是无符号整数
	addr  int //变量，在数据栈里的地址，也就是dx
}

/*
 * 符号表里加一项
 * 可能是常亮  a=2
 * 可能是变量 var v
 * 可能是过程 procedure p
 */
func (p *Parser) enter(k ident.Ident) {
	p.tx++
	tmp := identitem{
		name:  p.id,
		lev:   p.lev,
		kind:  k,
		value: p.num,
		addr:  p.dx,
	}
	if k == ident.Variable {
		//变量需要在数据栈里留空间
		p.dx++
	}
	p.table[p.tx] = tmp
}

// 生成中间代码
func (p *Parser) gen(fct fct.Fct, y, z int) {
	p.codes = append(p.codes, asm.Asm{
		Fct: fct,
		L:   y,
		A:   z,
	})
	p.cidx++
}

//返回标识符在符号表的位置
//检查标识符是否已经存在
//注意不同lev的不应该相同
func (p *Parser) checkintable(lit string) (id identitem, in bool) {
	for _, v := range p.table {
		if lit == v.name {
			return v, true
		}
	}
	return identitem{}, false
}

func tokintoks(toks []token.Token, tok token.Token) bool {
	for _, v := range toks {
		if v == tok {
			return true
		}
	}
	return false
}

/*
 * 检查终结符是不是在select集中
 */
func (p *Parser) test(toks1 []token.Token, toks2 []token.Token, t int) {
	toks1 = append(toks1, toks2...)
	if !tokintoks(toks1, p.tok) {
		p.scanner.Error(t)
		log.Println("出错的当前单词：" + p.tok.String() + p.lit)
		for !tokintoks(toks1, p.tok) && p.tok != token.PERIODSYM {
			p.Scan()
		}
	}
}

/*
 * 出现了const后进入这个过程，理应是常量
 * 常量声明
 */
func (p *Parser) constdeclaration() {
	if p.tok.IsIdent() {
		p.Scan()                                                //看下一个
		if p.tok == token.EQLSYM || p.tok == token.BECOMESSYM { //等号或赋值号
			if p.tok == token.BECOMESSYM {
				//这里是容错处理，报错，但是当做等号使用
				error()
			}

			p.Scan()
			//数字就加入符号表
			if p.tok.IsNumber() {
				p.enter(ident.Constant)
				p.Scan()
			} else {
				//TODO 等号后面不是数字，类型不可接受
				error()
			}
		} else {
			//TODO 不是等号或赋值号，非常量声明部分
			error()
		}
	} else {
		error()
	}
}

/*
 * 变量声明部分
 */
func (p *Parser) vardeclaration() {
	if p.tok.IsIdent() {
		p.enter(ident.Variable)
		p.Scan()
	} else {
		//TODO 不是变量名，报错
		error()
	}
}

/*
 * 因子的产生式
 * <因子> → <标识符>|<无符号整数>|(<表达式>)
 */
// Checked
func (p *Parser) factor(toks []token.Token) {
	p.test(selectset.FactorSelect, toks, 24)
	for tokintoks(selectset.FactorSelect, p.tok) {
		switch p.tok {
		case token.IDENTSYM: //是标识符
			//可能是常量或者变量，需要查询符号表确定位置
			it, in := p.checkintable(p.id)
			if in { //标识符定义了
				switch it.kind {
				case ident.Constant:
					p.gen(fct.Lit, 0, it.value)
				case ident.Variable:
					p.gen(fct.Lod, p.lev-it.lev, it.addr)
				case ident.Proc:
					p.scanner.Error(21)
					//TODO 不允许接受过程
				}
			} else {
				//TODO 标识符未定义
				p.scanner.Error(11)
			}
			p.Scan()
		case token.NUMBERSYM: //是无符号整数
			//TODO 判断数字有没有溢出
			p.gen(fct.Lit, 0, p.num)
			p.Scan()
		case token.LPARENTSYM: //是左括号
			p.Scan()
			p.expression(append(toks, token.RPARENTSYM))
			//表达式结束后是右括号
			if p.tok == token.RPARENTSYM {
				p.Scan()
			} else {
				//TODO 缺少右括号
				p.scanner.Error(22)
			}

		}
		p.test(toks, []token.Token{token.LPARENTSYM}, 23)
	}
}

/*
 * 项
 * <项> → <因子>{<乘除运算符><因子>}
 */
// Checked
func (p *Parser) term(toks []token.Token) {
	ttoks := append(toks, []token.Token{token.MULSYM, token.SLASHSYM}...)
	//第一个是因子
	p.factor(ttoks)

	//后面可以接无限个 乘除 因子
	for p.tok == token.MULSYM || p.tok == token.SLASHSYM {
		opt := p.tok
		p.Scan()
		p.factor(ttoks)
		if opt == token.MULSYM {
			p.gen(fct.Opr, 0, 4) //乘法
		} else {
			p.gen(fct.Opr, 0, 5) //除法
		}
	}
}

/*
 * 表达式
 * <表达式> → [+|-]<项>{<加减运算符><项>}
 */
func (p *Parser) expression(toks []token.Token) {
	var addop token.Token
	ttoks := append(toks, []token.Token{token.PLUSSYM, token.MINUSYM}...)
	if p.tok == token.PLUSSYM || p.tok == token.MINUSYM { //处理可能出现的正负号
		addop = p.tok
		p.Scan()
		p.term(ttoks)
		if addop == token.MINUSYM {
			p.gen(fct.Opr, 0, 1)
		}
	} else {
		p.term(ttoks)
	}
	//后面可以重复的加减运算和项
	for p.tok == token.PLUSSYM || p.tok == token.MINUSYM {
		addop = p.tok
		p.Scan()

		p.term(ttoks)
		if addop == token.PLUSSYM {
			p.gen(fct.Opr, 0, 2) //加
		} else {
			p.gen(fct.Opr, 0, 3) //减
		}
	}
}

/*
 * 条件
 * <条件> → <表达式><关系运算符><表达式>|ood<表达式>
 */
func (p *Parser) condition(toks []token.Token) {
	if p.tok == token.ODDSYM { //odd <表达式>
		p.Scan()
		p.expression(toks)
		p.gen(fct.Opr, 0, 6)
	} else {
		//表达式
		p.expression(append(toks, selectset.ExpressionSelect...))
		if p.tok.IsRelationOpr() { //是关系运算符
			relop := p.tok
			p.Scan()

			p.expression(toks)
			switch relop {
			case token.EQLSYM:
				p.gen(fct.Opr, 0, 8) // =
			case token.NEQSYM:
				p.gen(fct.Opr, 0, 9) // #
			case token.LESSYM:
				p.gen(fct.Opr, 0, 10) // <
			case token.GTRSYM:
				p.gen(fct.Opr, 0, 11) // >
			case token.LEQSYM:
				p.gen(fct.Opr, 0, 12) // <=
			case token.GEQSYM:
				p.gen(fct.Opr, 0, 13) // >=
			}
		}
	}
}

/*
 * 语句
 * <语句> → <赋值语句>|<条件语句>|<当型循环语句>|<过程调用语句>|<读语句>|<写语句>|<复合语句>|<空>
 * <赋值语句> → <标识符>:=<表达式>
 * <复合语句> → begin<语句>{ ；<语句>}<end>
 * <条件语句> → if<条件>then<语句>
 * <过程调用语句> → call<标识符>
 * <当型循环语句> → while<条件>do<语句>
 * <读语句> → read(<标识符>{ ，<标识符>})
 * <写语句> → write(<标识符>{，<标识符>})
 */
func (p *Parser) statement(toks []token.Token) {
	switch p.tok {
	case token.IDENTSYM: //赋值语句
		// <赋值语句> → <标识符>:=<表达式>
		id, in := p.checkintable(p.id)
		if in {
			if id.kind.IsConsttant() {
				// 不能改变常量的值
				p.scanner.Error(25)
				in = false
			}
		} else {
			p.scanner.Error(26)
			//变量未定义，不能赋值
		}
		p.Scan()
		if p.tok.IsBecome() {
			p.Scan()
		} else {
			//不是赋值语句？？！
			p.scanner.Error(13)
		}
		p.expression(toks)
		if in {
			p.gen(fct.Sto, p.lev-id.lev, id.addr)
		}
	case token.CALLSYM: //过程调用语句
		// call<标识符>
		p.Scan()
		if p.tok.IsIdent() {
			id, in := p.checkintable(p.lit)
			if in {
				if id.kind.IsProcedure() {
					p.gen(fct.Cal, p.lev-id.lev, id.addr)
				} else {
					// 调用的对象不是一个过程
					p.scanner.Error(15)
				}
			} else {
				p.scanner.Error(6)
				//未找到调用的过程
			}
			p.Scan()
		} else {
			//不是ident， 不是过程调用语句
			p.scanner.Error(27)
		}
	case token.IFSYM: //条件语句
		// if<条件>then<语句>
		p.Scan()
		p.condition(append(toks, []token.Token{token.THENSYM, token.DOSYM}...))
		if p.tok.IsThen() {
			p.Scan()
		} else {
			//TODO if条件之后未找到then
			error()
		}
		cidx := p.cidx //挖坑，false集的 a 需要时then后面的语句？
		p.gen(fct.Jpc, 0, 0)
		p.statement(toks)
		p.codes[cidx].A = p.cidx //false集，也就是else的部分？没有else?
	case token.BEGINSYM: // 复合语句
		// <复合语句> → begin<语句>{ ；<语句>}<end>
		p.Scan()
		ttoks := append(toks, []token.Token{token.SEMICOLOMSYM, token.ENDSYM}...)
		p.statement(ttoks)
		for p.tok.IsSemicolom() {
			p.Scan()
			p.statement(ttoks)
		}
		if p.tok.IsEnd() {
			p.Scan()
		} else {
			//TODO 过程没有结束符号
			error()
		}
	case token.WHILESYM: //while循环
		// <当型循环语句> → while<条件>do<语句>
		cidx1 := p.cidx //判断前面，循环体结束后需要跳过来
		p.Scan()
		p.condition(append(toks, token.DOSYM))
		cidx2 := p.cidx //退出循环体的地址后面分配好代码后回填
		p.gen(fct.Jpc, 0, 0)
		if p.tok.IsDo() {
			p.Scan()
		} else {
			//TODO 缺少do，可能忘写了
			error()
		}
		p.statement(toks)
		p.gen(fct.Jmp, 0, cidx1)
		p.codes[cidx2].A = p.cidx
	//	* <读语句> → read(<标识符>{ ，<标识符>})
	//	* <写语句> → write(<标识符>{，<标识符>})
	case token.READSYM:
		p.Scan()
		if p.tok.IsLparent() {
			p.Scan()
			if p.tok.IsIdent() {
				//检查变量表，应该是一个已经定义的变量
				id, in := p.checkintable(p.lit)
				if in {
					if !id.kind.IsVariable() {
						// 不能改变常量的值
						p.scanner.Error(25)
						in = false
					}
				} else {
					p.scanner.Error(26)
					//变量未定义，不能赋值
				}
				if in {
					p.gen(fct.Opr, 0, 14)                 //读入数字放栈顶
					p.gen(fct.Sto, p.lev-id.lev, id.addr) //从栈顶放到相应位置
				}
				p.Scan()
			} else {
				//TODO 应该是一个标识符，但是这里不是
				error()
			}
		}
		for p.tok.IsComma() {
			p.Scan()
			if p.tok.IsIdent() {
				//检查变量表，应该是一个已经定义的变量
				id, in := p.checkintable(p.lit)
				if in {
					if id.kind.IsConsttant() {
						// 不能改变常量的值
						p.scanner.Error(25)
						in = false
					}
				} else {
					p.scanner.Error(26)
					//变量未定义，不能赋值
				}
				if in {
					p.gen(fct.Opr, 0, 14)                 //读入数字放栈顶
					p.gen(fct.Sto, p.lev-id.lev, id.addr) //从栈顶放到相应位置
				}
			}
			p.Scan()
		}
		if p.tok.IsRparent() {
			p.Scan()
		} else {
			error()
		}
	case token.WRITESYM:
		p.Scan()
		if p.tok.IsLparent() {
			p.Scan()
			if p.tok.IsIdent() {
				//检查变量表，应该是一个已经定义的变量
				id, in := p.checkintable(p.id)
				if in {
					if id.kind.IsProcedure() {
						// 不能读过程
						p.scanner.Error(28)
						in = false
					} else if id.kind.IsConsttant() {
						p.gen(fct.Lit, 0, id.value) //从相应位置读到栈顶
						p.gen(fct.Opr, 0, 15)       //从栈顶显示出来
						in = false
					}
				} else {
					p.scanner.Error(26)
					//变量未定义，不能赋值
				}
				if in {
					p.gen(fct.Lod, p.lev-id.lev, id.addr) //从相应位置读到栈顶
					p.gen(fct.Opr, 0, 15)                 //从栈顶显示出来
				}
				p.Scan()
			} else {
				//TODO 应该是一个标识符，但是这里不是
				error()
			}
		}
		for p.tok.IsComma() {
			p.Scan()
			if p.tok.IsIdent() {
				//检查变量表，应该是一个已经定义的变量
				id, in := p.checkintable(p.lit)
				if in {
					if id.kind.IsProcedure() {
						// 不能改变常量的值
						p.scanner.Error(28)
						in = false
					} else if id.kind.IsConsttant() {
						p.gen(fct.Lit, 0, id.value) //从相应位置读到栈顶
						p.gen(fct.Opr, 0, 15)       //从栈顶显示出来
						in = false
					}
				} else {
					p.scanner.Error(26)
					//变量未定义，不能赋值
				}
				if in {
					p.gen(fct.Lod, p.lev-id.lev, id.addr) //把数据放到栈顶
					p.gen(fct.Opr, 0, 15)                 //读入数字放栈顶
				}
			}
			p.Scan()
		}
		if p.tok.IsRparent() {
			p.Scan()
		} else {
			error()
		}
	}

	p.test(toks, []token.Token{}, 19)
}

/*
 *〈程序〉→〈分程序>.
 *〈分程序〉→ [<常量说明部分>][<变量说明部分>][<过程说明部分>]〈语句〉
 * <常量说明部分> → CONST<常量定义>{ ,<常量定义>}；
 * <变量说明部分> → VAR<标识符>{ ,<标识符>}；
 * <过和说明部分> → <过程首部><分程度>；{<过程说明部分>}
 * <过程首部> → procedure<标识符>；
 */
func (p *Parser) block(toks []token.Token) {
	p.dx = 3
	tx0 := p.tx
	p.table[p.tx].addr = p.cidx
	p.gen(fct.Jmp, 0, 0)

	if p.lev > MAX_LEV {
		//TODO 嵌套层次太大
		error()
	}
	for { //声明部分
		//常量声明部分
		// * <常量说明部分> → CONST<常量定义>{ ,<常量定义>}；
		// * <常量定义> → <标识符>=<无符号整数>
		if p.tok.IsConst() {
			p.Scan()
			for p.tok.IsIdent() {
				//第一个常量声明
				p.constdeclaration()
				//逗号打头的多个常量的声明
				for p.tok.IsComma() {
					p.Scan()
					p.constdeclaration()
				}
				//分号结束
				if p.tok.IsSemicolom() {
					p.Scan()
				} else {
					//TODO 缺少分号
					error()
				}
			}
		}

		//变量声明部分
		// * <变量说明部分> → VAR<标识符>{ ,<标识符>}；
		if p.tok.IsVar() {
			p.Scan()
			for p.tok.IsIdent() {
				p.vardeclaration()
				for p.tok.IsComma() {
					p.Scan()
					p.vardeclaration()
				}
				//分号结束
				if p.tok.IsSemicolom() {
					p.Scan()
				} else {
					//TODO 缺少分号
					error()
				}
			}
		}

		//过程声明部分
		// * <过和说明部分> → <过程首部><分程度>；{<过程说明部分>}
		// * <过程首部> → procedure<标识符>；
		for p.tok.IsProc() {
			p.Scan()
			if p.tok.IsIdent() {
				p.enter(ident.Proc)
				p.Scan()
			} else {
				//TODO proc后面需要紧跟一个标识符
				error()
			}
			if p.tok.IsSemicolom() {
				p.Scan()
			} else {
				//TODO 过程首部里面缺少分号
				error()
			}

			p.lev++
			tx1 := p.tx
			dx1 := p.dx
			p.block(append(toks, token.SEMICOLOMSYM))
			p.lev--
			p.tx = tx1
			p.dx = dx1

			if p.tok.IsSemicolom() {
				p.Scan()
				p.test(append(selectset.StatementSelect, []token.Token{token.IDENTSYM, token.PROCSYM}...), toks, 6)
			} else {
				//TODO 缺少分号
				error()
			}
		}

		p.test(append(selectset.StatementSelect, token.IDENTSYM), selectset.DeclareSelect, 7)
		if !tokintoks(selectset.DeclareSelect, p.tok) {
			break
		}
	}

	p.codes[p.table[tx0].addr].A = p.cidx
	p.table[tx0].addr = p.cidx
	//cx0 := p.cidx
	p.gen(fct.Int, 0, p.dx)
	p.statement(append(toks, []token.Token{token.SEMICOLOMSYM, token.ENDSYM}...))
	p.gen(fct.Opr, 0, 0)
	p.test(toks, []token.Token{}, 8)
}

func (p *Parser) Start() {
	p.block(append(append(selectset.DeclareSelect, selectset.StatementSelect...), token.PERIODSYM))
	log.Printf("结束\n")
}
