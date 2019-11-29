/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package selectset

import "pl0compiler/token"

var (
	StatementSelect = []token.Token{
		token.BEGINSYM,
		token.CALLSYM,
		token.IFSYM,
		token.WHILESYM,
	}
)
