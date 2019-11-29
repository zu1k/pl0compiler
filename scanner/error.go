/**
 * @auther:  zu1k
 * @date:    2019/11/10
 */
package scanner

import (
	"log"
)

var error = [...]string{
	/*  0 */ "",
	/*  1 */ "找到‘:=’，但是期望的是‘=’。",
	/*  2 */ "‘=’后面必须接一个数字。",
	/*  3 */ "标识符后面必须是‘=’。",
	/*  4 */ "在‘const’, ‘var’, 或‘procedure’后面必须是一个标识符。",
	/*  5 */ "缺少‘,’或‘;’。",
	/*  6 */ "过程名错误，未找到该过程，或者连同名的常量变量都没有",
	/*  7 */ "期待一个语句。",
	/*  8 */ "语句后面是错误的符号。",
	/*  9 */ "期待‘.’。",
	/* 10 */ "期待‘.’。",
	/* 11 */ "未声明的标识符。",
	/* 12 */ "非法声明。",
	/* 13 */ "期待赋值号‘:=’。",
	/* 14 */ "在‘call’后面必须接一个标识符(过程)。",
	/* 15 */ "常量和变量(非过程标识符)不能被call调用。",
	/* 16 */ "'then' expected.",
	/* 17 */ "';' or 'end' expected.",
	/* 18 */ "'do' expected.",
	/* 19 */ "Incorrect symbol.",
	/* 20 */ "Relative operators expected.",
	/* 21 */ "Procedure identifier can not be in an expression.",
	/* 22 */ "Missing ')'.",
	/* 23 */ "The symbol can not be followed by a factor.",
	/* 24 */ "The symbol can not be as the beginning of an expression.",
	/* 25 */ "不能改变常量的值",
	/* 26 */ "变量未定义，无法赋值",
	/* 27 */ "call后面应该是一个过程类型的标识符，这里连标识符都不是",
	/* 28 */ "不能输出一个过程",
	/* 29 */ "",
	/* 30 */ "",
	/* 31 */ "数太大了。",
	/* 32 */ "层次太深了。",
}

func (s *Scanner) Error(t int) {
	log.Printf("Error %d: %s\tON Line: %d", t, error[t], s.line)
}
