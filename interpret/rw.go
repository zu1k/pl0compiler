/**
 * @auther:  zu1k
 * @date:    2019/11/13
 */
package interpret

import (
	"fmt"
	"log"
)

func read() (i int) {
	log.Printf("请输入一个无符号整数：")
	_, err := fmt.Scanln(&i)
	for err != nil {
		log.Printf("输入的不是无符号整数，请重新输入：")
		_, err = fmt.Scanln(&i)
	}
	return
}

func write(i int) {
	log.Printf("Write: %d", i)
}
