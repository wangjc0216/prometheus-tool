package main

import (
	"flag"
	"fmt"
	"github.com/disiqueira/gotree"
	"github.com/prometheus/prometheus/promql/parser"
	"os"
	"reflect"
)

//func main() {
//	artist := gotree.New("Pantera")
//	album := artist.Add("Far Beyond Driven")
//	album.Add("5 minutes Alone")
//
//	fmt.Println(artist.Print())
//}

func main() {
	qs := flag.String("promql", "", "please input your promql")
	flag.Parse()
	expr, err := parser.ParseExpr(*qs)
	if err != nil {
		fmt.Println("parse error")
		os.Exit(-1)
	}
	//log("", 0, expr)

	tree := gotree.New("expr")
	walk(tree, expr)
	fmt.Println(tree.Print())
}

func walk(tree gotree.Tree, expr parser.Expr) {
	//str += field
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		//todo tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))
		tree.Add(fmt.Sprintf("Op:%s\n", itemTypes[e.Op]))
		tree.Add(fmt.Sprintf("Grouping:%v\n", e.Grouping))
		tree.Add(fmt.Sprintf("Without:%v\n", e.Without))
		tree.Add(fmt.Sprintf("PosRange:%v\n", e.PosRange))

		exprTree := tree.Add(fmt.Sprintln("Expr:"))
		walk(exprTree, e.Expr)
		paramTree := tree.Add(fmt.Sprintln("Param:"))
		walk(paramTree, e.Param)
	case *parser.BinaryExpr:
		//todo tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))
		tree.Add(fmt.Sprintf("Op:%s\n", itemTypes[e.Op]))
		tree.Add(fmt.Sprintf("ReturnBool:%v\n", e.ReturnBool))
		tree.Add(fmt.Sprintf("VectorMatching:%v\n", e.VectorMatching))

		LHSTree := tree.Add(fmt.Sprintln("LHS:"))
		walk(LHSTree, e.LHS)
		RHSTree := tree.Add(fmt.Sprintln("RHS:"))
		walk(RHSTree, e.RHS)

	case *parser.ParenExpr:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		//tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

		exprTree := tree.Add(fmt.Sprintln("Expr:"))
		walk(exprTree, e.Expr)

	case *parser.Call:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))
		//todo tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		tree.Add(fmt.Sprintf("Func:%v\n", *e.Func))
		tree.Add(fmt.Sprintf("PosRange:%v\n", e.PosRange))
		argsTree := tree.Add(fmt.Sprintln("Args:"))
		for i, expr := range e.Args {
			argTree := argsTree.Add(fmt.Sprintf("arg[%d]", i))
			walk(argTree, expr)
		}
	case *parser.SubqueryExpr:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		//todo tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		tree.Add(fmt.Sprintf("Range:%v", e.Range))
		tree.Add(fmt.Sprintf("OriginalOffset:%v", e.OriginalOffset))
		tree.Add(fmt.Sprintf("Offset:%v", e.Offset))
		tree.Add(fmt.Sprintf("Timestamp:%v", e.Timestamp))
		tree.Add(fmt.Sprintf("StartOrEnd:%v", itemTypes[e.StartOrEnd]))
		tree.Add(fmt.Sprintf("Step:%v", e.Step))
		tree.Add(fmt.Sprintf("EndPos:%v", e.EndPos))

		exprTree := tree.Add(fmt.Sprintln("Expr:"))
		walk(exprTree, e.Expr)

	case *parser.UnaryExpr:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		//fmt.Printf("%s startpos:%v\n", str, e.StartPos)
		//fmt.Printf("%s Op:%s\n", str, itemTypes[e.Op])
		exprTree := tree.Add(fmt.Sprintln("Expr:"))
		walk(exprTree, e.Expr)

	case *parser.MatrixSelector:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		//fmt.Printf("%s endpos:%v\n", str, e.EndPos)
		//fmt.Printf("%s range:%v\n", str, e.Range)
		vs := tree.Add("VectorSelector:")
		walk(vs, e.VectorSelector)

	case *parser.StepInvariantExpr:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		exprTree := tree.Add("Expr")
		walk(exprTree, e.Expr)

	case *parser.NumberLiteral:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))
		//fmt.Printf("%s val:%v\n", str, e.Val)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

	case *parser.StringLiteral:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))

		//fmt.Printf("%s val:%v\n", str, e.Val)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

	case *parser.VectorSelector:
		tree.Add(fmt.Sprintf("Type:%s\n", reflect.TypeOf(e).String()))

		tree.Add(fmt.Sprintf("[%s]:%s\n", reflect.TypeOf(e).String(), e.String()))

		//fmt.Printf("%s name:%v\n", str, e.Name)
		//fmt.Printf("%s duration:%v\n", str, e.OriginalOffset)
		//fmt.Printf("%s offset:%v\n", str, e.Offset)
		//fmt.Printf("%s timestamp:%v\n", str, e.Timestamp)
		//fmt.Printf("%s StartOrEnd:%v\n", str, itemTypes[e.StartOrEnd])
		//fmt.Printf("%s LabelMatchers:%v\n", str, e.LabelMatchers)
		//fmt.Printf("%s UnexpandedSeriesSet:%v\n", str, e.UnexpandedSeriesSet)
		//fmt.Printf("%s Series:%v\n", str, e.Series)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

	}
}

func log(field string, level int, expr parser.Expr) {
	var str string
	for i := 0; i < level; i++ {
		str += "    "
	}
	//str += field
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		//Op       ItemType // The used aggregation operation.
		//Expr     Expr     // The Vector expression over which is aggregated.
		//Param    Expr     // Parameter used by some aggregators.
		//Grouping []string // The labels by which to group the Vector.
		//Without  bool     // Whether to drop the given labels rather than keep them.
		//PosRange PositionRange
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s Op:%s\n", str, itemTypes[e.Op])
		//fmt.Printf("%s Grouping:%v\n", str, e.Grouping)
		//fmt.Printf("%s without:%v\n", str, e.Without)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)
		fmt.Printf("%s Expr:\n", str)
		log("Expr", level+1, e.Expr)
		fmt.Printf("%s Param:\n", str)
		log("Param", level+1, e.Param)
	case *parser.BinaryExpr:
		//Op       ItemType // The operation of the expression.
		//LHS, RHS Expr     // The operands on the respective sides of the operator.
		//
		//// The matching behavior for the operation if both operands are Vectors.
		//// If they are not this field is nil.
		//VectorMatching *VectorMatching
		//
		//// If a comparison operator, return 0/1 rather than filtering.
		//ReturnBool bool
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s Op:%s\n", str, itemTypes[e.Op])
		//fmt.Printf("%s returnBool:%v\n", str, e.ReturnBool)
		//fmt.Printf("%s vectorMatching:%v\n", str, e.VectorMatching)
		fmt.Printf("%s LHS:\n", str)
		log("LHS", level+1, e.LHS)
		fmt.Printf("%s RHS:\n", str)
		log("RHS", level+1, e.RHS)

	case *parser.ParenExpr:
		//Expr     Expr
		//PosRange PositionRange
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)
		fmt.Printf("%s Expr:\n", str)
		log("Expr", level+1, e.Expr)
	case *parser.Call:
		//Func *Function   // The function that was called.
		//Args Expressions // Arguments used in the call.
		//PosRange PositionRange
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s func:%v\n", str, *e.Func)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)
		for _, expr := range e.Args {
			fmt.Printf("%s Arg:\n", str)
			log("Arg", level+1, expr)
		}
	case *parser.SubqueryExpr:
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s range:%v\n", str, e.Range)
		//fmt.Printf("%s originalOffset:%v\n", str, e.OriginalOffset)
		//fmt.Printf("%s offset:%v\n", str, e.Offset)
		//fmt.Printf("%s Timestamp:%v\n", str, e.Timestamp)
		//fmt.Printf("%s StartOrEnd:%v\n", str, itemTypes[e.StartOrEnd])
		//fmt.Printf("%s Step:%v\n", str, e.Step)
		//fmt.Printf("%s endpos:%v\n", str, e.EndPos)
		fmt.Printf("%s Expr:\n", str)
		log("Expr", level+1, e.Expr)

		//Expr  Expr
		//Range time.Duration
		//// OriginalOffset is the actual offset that was set in the query.
		//// This never changes.
		//OriginalOffset time.Duration
		//// Offset is the offset used during the query execution
		//// which is calculated using the original offset, at modifier time,
		//// eval time, and subquery offsets in the AST tree.
		//Offset     time.Duration
		//Timestamp  *int64
		//StartOrEnd ItemType // Set when @ is used with start() or end()
		//Step       time.Duration
		//EndPos Pos

	case *parser.UnaryExpr:
		//Op   ItemType
		//Expr Expr
		//
		//StartPos Pos
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s startpos:%v\n", str, e.StartPos)
		//fmt.Printf("%s Op:%s\n", str, itemTypes[e.Op])
		fmt.Printf("%s Expr:\n", str)
		log("Expr", level+1, e.Expr)

	case *parser.MatrixSelector:
		//VectorSelector Expr
		//Range          time.Duration
		//EndPos Pos
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s endpos:%v\n", str, e.EndPos)
		//fmt.Printf("%s range:%v\n", str, e.Range)
		log("VectorSelector", level+1, e.VectorSelector)
	case *parser.StepInvariantExpr:

		log("Expr", level+1, e.Expr)
	case *parser.NumberLiteral:
		//Val float64
		//PosRange PositionRange
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s val:%v\n", str, e.Val)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

	case *parser.StringLiteral:
		//Val      string
		//PosRange PositionRange
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s val:%v\n", str, e.Val)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

	case *parser.VectorSelector:
		fmt.Printf("%s [%s]:%s\n", str, reflect.TypeOf(e).String(), e.String())

		//fmt.Printf("%s name:%v\n", str, e.Name)
		//fmt.Printf("%s duration:%v\n", str, e.OriginalOffset)
		//fmt.Printf("%s offset:%v\n", str, e.Offset)
		//fmt.Printf("%s timestamp:%v\n", str, e.Timestamp)
		//fmt.Printf("%s StartOrEnd:%v\n", str, itemTypes[e.StartOrEnd])
		//fmt.Printf("%s LabelMatchers:%v\n", str, e.LabelMatchers)
		//fmt.Printf("%s UnexpandedSeriesSet:%v\n", str, e.UnexpandedSeriesSet)
		//fmt.Printf("%s Series:%v\n", str, e.Series)
		//fmt.Printf("%s posrange:%v\n", str, e.PosRange)

		//Name string
		//// OriginalOffset is the actual offset that was set in the query.
		//// This never changes.
		//OriginalOffset time.Duration
		//// Offset is the offset used during the query execution
		//// which is calculated using the original offset, at modifier time,
		//// eval time, and subquery offsets in the AST tree.
		//Offset        time.Duration
		//Timestamp     *int64
		//StartOrEnd    ItemType // Set when @ is used with start() or end()
		//LabelMatchers []*labels.Matcher
		//
		//// The unexpanded seriesSet populated at query preparation time.
		//UnexpandedSeriesSet storage.SeriesSet
		//Series              []storage.Series
		//
		//PosRange PositionRange
	}
}

var itemTypes = map[parser.ItemType]string{
	57346: "EQL",
	57347: "BLANK",
	57348: "COLON",
	57349: "COMMA",
	57350: "COMMENT",
	57351: "DURATION",
	57352: "EOF",
	57353: "ERROR",
	57354: "IDENTIFIER",
	57355: "LEFT_BRACE",
	57356: "LEFT_BRACKET",
	57357: "LEFT_PAREN",
	57358: "METRIC_IDENTIFIER",
	57359: "NUMBER",
	57360: "RIGHT_BRACE",
	57361: "RIGHT_BRACKET",
	57362: "RIGHT_PAREN",
	57363: "SEMICOLON",
	57364: "SPACE",
	57365: "STRING",
	57366: "TIMES",
	57367: "operatorsStart",
	57368: "ADD",
	57369: "DIV",
	57370: "EQLC",
	57371: "EQL_REGEX",
	57372: "GTE",
	57373: "GTR",
	57374: "LAND",
	57375: "LOR",
	57376: "LSS",
	57377: "LTE",
	57378: "LUNLESS",
	57379: "MOD",
	57380: "MUL",
	57381: "NEQ",
	57382: "NEQ_REGEX",
	57383: "POW",
	57384: "SUB",
	57385: "AT",
	57386: "ATAN2",
	57387: "operatorsEnd",
	57388: "aggregatorsStart",
	57389: "AVG",
	57390: "BOTTOMK",
	57391: "COUNT",
	57392: "COUNT_VALUES",
	57393: "GROUP",
	57394: "MAX",
	57395: "MIN",
	57396: "QUANTILE",
	57397: "STDDEV",
	57398: "STDVAR",
	57399: "SUM",
	57400: "TOPK",
	57401: "aggregatorsEnd",
	57402: "keywordsStart",
	57403: "BOOL",
	57404: "BY",
	57405: "GROUP_LEFT",
	57406: "GROUP_RIGHT",
	57407: "IGNORING",
	57408: "OFFSET",
	57409: "ON",
	57410: "WITHOUT",
	57411: "keywordsEnd",
	57412: "preprocessorStart",
	57413: "START",
	57414: "END",
	57415: "preprocessorEnd",
	57416: "startSymbolsStart",
	57417: "START_METRIC",
	57418: "START_SERIES_DESCRIPTION",
	57419: "START_EXPRESSION",
	57420: "START_METRIC_SELECTOR",
	57421: "startSymbolsEnd",
}
