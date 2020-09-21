/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package selectset

import "github.com/zu1k/pl0complier/token"

var (
	StatementSelect = []token.Token{
		token.BEGINSYM,
		token.CALLSYM,
		token.IFSYM,
		token.WHILESYM,
	}
)
