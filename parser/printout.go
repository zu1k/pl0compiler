/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package parser

import (
	"log"
)

func (p *Parser) Sout() {
	log.Println("============================ 中间代码 Begin ============================")
	log.Printf("\t\t%s\t\t%s\t\t%s\n", "功能码", "层次差", "位移量")
	for i, v := range p.codes {
		log.Printf("\t%d\t%s\t\t%d\t\t%d\n", i, v.Fct.String(), v.L, v.A)
	}
	log.Println("============================ 中间代码 End ============================\n")

	log.Println("============================ 标识符表 Begin ============================")
	log.Printf("\t%s\t%s\t\t%s\t\t%s\t\t%s\n", "名", "类型", "层次", "值", "地址")
	for _, v := range p.table {
		if v.name != "" {
			//if true {
			log.Printf("\t%s\t%s\t\t%d\t\t%d\t\t%d\n", v.name, v.kind.String(), v.lev, v.value, v.addr)
		}

	}
	log.Println("============================ 标识符表 End ============================\n")
}
