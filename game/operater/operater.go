package operater

import (
	"fmt"
	"lucascript/charset"
	"lucascript/function"
	"lucascript/game/context"
	"lucascript/paramter"
	"lucascript/script"
)

// Operater 需定制指令
type Operater interface {
	MESSAGE(code *script.CodeLine) string
}

// LucaOperater 通用指令
type LucaOperater interface {
	UNDEFINE(code *script.CodeLine, opname string) string
	EQU(ctx *context.Context) function.HandlerFunc
	ADD(code *script.CodeLine) string
	IFN(ctx *context.Context) function.HandlerFunc
	IFY(code *script.CodeLine) function.HandlerFunc
	GOTO(code *script.CodeLine) function.HandlerFunc
	JUMP(code *script.CodeLine) function.HandlerFunc
	FARCALL(code *script.CodeLine) function.HandlerFunc
}

// LucaOperate 通用指令
type LucaOperate struct {
	ExprCharset charset.Charset
	TextCharset charset.Charset
	LabelMap    map[int]int
}

func (g *LucaOperate) UNDEFINE(code *script.CodeLine, opcode string) string {
	if len(opcode) == 0 {
		opcode = ToString("%X", code.Opcode)
	}
	list, end := AllToUint16(code.CodeBytes)
	str := ""
	for _, num := range list {
		str += fmt.Sprintf(", %d", num)
	}
	if end > 0 {
		str += fmt.Sprintf(", 0x%X", code.CodeBytes[end])
	}
	if len(str) >= 2 {
		str = str[2:]
	}
	return ToString(`%s (%s)`, opcode, str)
}
func (g *LucaOperate) EQU(ctx *context.Context) function.HandlerFunc {
	code := ctx.Code()
	var key paramter.LUint16
	var value paramter.LUint16

	next := GetParam(code.CodeBytes, &key)
	GetParam(code.CodeBytes, &value, next)

	fun := function.EQU{}
	return func() int {
		// 这里是执行 与虚拟机逻辑有关的代码
		ctx.Variable.Set(ToString("#%d", key.Data), int(value.Data))
		// 这里执行与游戏相关代码，内部与虚拟机无关联
		fun.Call([]paramter.Paramter{&key, &value})
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
		return 0
	}
}
func (g *LucaOperate) ADD(code *script.CodeLine) string {
	opcode := "add"
	key := ToUint16(code.CodeBytes[0:2])
	exprStr, _ := DecodeString(code.CodeBytes, 2, 0, g.ExprCharset)
	return ToString(`%d:%s (#%d, %s)`, code.Pos, opcode, key, exprStr)
}

func (g *LucaOperate) IFN(ctx *context.Context) function.HandlerFunc {
	code := ctx.Code()
	var jumpPos paramter.LUint32
	var exprStr paramter.LString
	next := GetParam(code.CodeBytes, &exprStr, 0, 0, g.ExprCharset)
	GetParam(code.CodeBytes, &jumpPos, next, 4)

	fun := function.IFN{}
	return func() int {
		// 这里是执行 与虚拟机逻辑有关的代码
		eip := 0
		res := true // res:=expr(ifExprStr)
		if !res {
			eip = 0
		}
		// 这里执行与游戏相关代码，内部与虚拟机无关联
		fun.Call([]paramter.Paramter{&exprStr, &jumpPos})
		ctx.ChanEIP <- 0
		return eip
	}
}

func (g *LucaOperate) FARCALL(code *script.CodeLine) function.HandlerFunc {
	params := make([]paramter.Paramter, 0, 3)
	var index paramter.LUint16
	var fileStr paramter.LString
	var jumpPos paramter.LUint32

	next := GetParam(code.CodeBytes, &index)
	params = append(params, &index)

	next = GetParam(code.CodeBytes, &fileStr, next, 0, g.ExprCharset)
	params = append(params, &fileStr)

	GetParam(code.CodeBytes, &jumpPos, next)
	params = append(params, &jumpPos)

	fun := function.FARCALL{}
	return func() int {
		// 这里是执行内容
		return fun.Call(params)
	}
}

func (g *LucaOperate) IFY(code *script.CodeLine) function.HandlerFunc {
	params := make([]paramter.Paramter, 0, 2)
	var jumpPos paramter.LUint32
	var exprStr paramter.LString
	next := GetParam(code.CodeBytes, &exprStr, 0, 0, g.ExprCharset)
	params = append(params, &exprStr)

	GetParam(code.CodeBytes, &jumpPos, next, 4)
	params = append(params, &jumpPos)
	fun := function.IFN{}
	return func() int {
		// 这里是执行内容
		return fun.Call(params)
	}
}

func (g *LucaOperate) GOTO(code *script.CodeLine) function.HandlerFunc {
	params := make([]paramter.Paramter, 0, 1)
	var jumpPos paramter.LUint32
	GetParam(code.CodeBytes, &jumpPos)
	params = append(params, &jumpPos)
	fun := function.GOTO{}
	return func() int {
		// 这里是执行内容
		return fun.Call(params)
	}
}

func (g *LucaOperate) JUMP(code *script.CodeLine) function.HandlerFunc {

	params := make([]paramter.Paramter, 0, 2)
	var jumpPos paramter.LUint32
	var fileStr paramter.LString
	next := GetParam(code.CodeBytes, &fileStr, 0, 0, g.ExprCharset)
	params = append(params, &fileStr)

	GetParam(code.CodeBytes, &jumpPos, next, 4)
	params = append(params, &jumpPos)
	fun := function.JUMP{}
	return func() int {
		// 这里是执行内容
		return fun.Call(params)
	}

}
