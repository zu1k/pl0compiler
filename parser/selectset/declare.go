/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package selectset

import "pl0compiler/token"

var (
	DeclareSelect = []token.Token{
		token.CONSTSYM,
		token.VARSYM,
		token.PROCSYM,
	}
)
