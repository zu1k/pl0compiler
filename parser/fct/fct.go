/**
 * @auther:  zu1k
 * @date:    2019/11/9
 */
package fct

type Fct int

const (
	Lit = iota
	Opr
	Lod
	Sto
	Cal
	Int
	Jmp
	Jpc
)

var fcttable = [...]string{
	Lit: "LIT",
	Opr: "OPR",
	Lod: "LOD",
	Sto: "STO",
	Cal: "CAL",
	Int: "INT",
	Jmp: "JMP",
	Jpc: "JPC",
}

func (f Fct) String() string {
	return fcttable[f]
}

/*  lit 0, a : load constant a
    opr 0, a : execute operation a
    lod l, a : load variable l, a
    sto l, a : store variable l, a
    cal l, a : call procedure a at level l
    Int 0, a : increment t-register by a
    jmp 0, a : jump to a
    jpc 0, a : jump conditional to a       */
