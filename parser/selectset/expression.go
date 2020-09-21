/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package selectset

import "github.com/zu1k/pl0complier/token"

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
