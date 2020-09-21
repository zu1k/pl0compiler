/**
 * @auther:  zu1k
 * @date:    2019/11/12
 */
package selectset

import "github.com/zu1k/pl0complier/token"

var (
	DeclareSelect = []token.Token{
		token.CONSTSYM,
		token.VARSYM,
		token.PROCSYM,
	}
)
