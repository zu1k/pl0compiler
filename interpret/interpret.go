/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package interpret

import (
	"log"
	"pl0compiler/asm"
	"pl0compiler/parser/fct"
)

var s [200]int

func base(b int, l int) (c int) {
	c = b
	for l > 0 {
		c = s[c]
		l--
	}
	return
}

func Interpret(asms []asm.Asm) {
	var p, b, t int // 程序寄存器PC、基地址寄存器、栈顶寄存器
	var i asm.Asm   // 指令寄存器

	log.Printf("Start to run")
	t = 0
	b = 1
	p = 0

	s[1] = 0
	s[2] = 0
	s[3] = 0

	for {
		i = asms[p]
		p++

		switch i.Fct {
		case fct.Lit:
			t++
			s[t] = i.A
		case fct.Opr: //运算
			switch i.A {
			case 0: //返回指令
				t = b - 1
				p = s[t+3]
				b = s[t+2]
			case 1: //负号
				s[t] = -s[t]
			case 2: //加法
				t--
				s[t] = s[t] + s[t+1]
			case 3: //减法
				t--
				s[t] = s[t] - s[t+1]
			case 4: // 乘法
				t--
				s[t] = s[t] * s[t+1]
			case 5: // 除法
				t--
				s[t] = s[t] / s[t+1]
			case 6: // odd
				s[t] = s[t] % 2
			case 7:

			case 8: // ==
				t--
				bb := s[t] == s[t+1]
				if bb {
					s[t] = 1
				} else {
					s[t] = 0
				}
			case 9: // !=
				t--
				bb := s[t] != s[t+1]
				if bb {
					s[t] = 1
				} else {
					s[t] = 0
				}
			case 10: //<
				t--
				bb := s[t] < s[t+1]
				if bb {
					s[t] = 1
				} else {
					s[t] = 0
				}
			case 11: //>
				t--
				bb := s[t] > s[t+1]
				if bb {
					s[t] = 1
				} else {
					s[t] = 0
				}
			case 12: //<=
				t--
				bb := s[t] <= s[t+1]
				if bb {
					s[t] = 1
				} else {
					s[t] = 0
				}
			case 13: //>=
				t--
				bb := s[t] >= s[t+1]
				if bb {
					s[t] = 1
				} else {
					s[t] = 0
				}
			case 14: //read
				t++
				s[t] = read()
			case 15: //write
				write(s[t])
				t--
			}
		case fct.Lod: // 调用变量值指令
			t = t + 1
			s[t] = s[base(b, i.L)+i.A]
		case fct.Sto: // 将值存入变量指令
			s[base(b, i.L)+i.A] = s[t]
			//log.Printf("%10d\n", s[t])
			t = t - 1
		case fct.Cal: // 过程调用，产生新的块标记
			s[t+1] = base(b, i.L)
			s[t+2] = b
			s[t+3] = p // 记录返回地址等参数
			b = t + 1
			p = i.A
		case fct.Int: // 开内存空间
			t = t + i.A
		case fct.Jmp: // 无条件跳转指令
			p = i.A
		case fct.Jpc: // 栈顶为0跳转
			if s[t] == 0 {
				p = i.A
			}
			t--
		}

		if p == 0 {
			break
		}
	}
	log.Printf("Completed\n")
}
