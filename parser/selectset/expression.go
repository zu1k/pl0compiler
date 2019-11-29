/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package selectset

import "pl0compiler/token"

var (
	ExpressionSelect = []token.Token{
		token.EQLSYM,
		token.NEQSYM,
		token.LESSYM,
		token.LEQSYM,
		token.GTRSYM,
		token.GEQSYM,
	}
)
