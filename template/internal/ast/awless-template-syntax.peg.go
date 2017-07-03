package ast

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleScript
	ruleStatement
	ruleAction
	ruleEntity
	ruleDeclaration
	ruleValueExpr
	ruleCmdExpr
	ruleParams
	ruleParam
	ruleIdentifier
	ruleNoRefValue
	ruleValue
	ruleCustomTypedValue
	ruleOtherParamValue
	ruleDoubleQuotedValue
	ruleSingleQuotedValue
	ruleCSVValue
	ruleCidrValue
	ruleIpValue
	ruleIntRangeValue
	ruleRefValue
	ruleRefsListValue
	ruleAliasValue
	ruleHoleValue
	ruleComment
	ruleSingleQuote
	ruleDoubleQuote
	ruleWhiteSpacing
	ruleMustWhiteSpacing
	ruleEqual
	ruleBlankLine
	ruleWhitespace
	ruleEndOfLine
	ruleEndOfFile
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
)

var rul3s = [...]string{
	"Unknown",
	"Script",
	"Statement",
	"Action",
	"Entity",
	"Declaration",
	"ValueExpr",
	"CmdExpr",
	"Params",
	"Param",
	"Identifier",
	"NoRefValue",
	"Value",
	"CustomTypedValue",
	"OtherParamValue",
	"DoubleQuotedValue",
	"SingleQuotedValue",
	"CSVValue",
	"CidrValue",
	"IpValue",
	"IntRangeValue",
	"RefValue",
	"RefsListValue",
	"AliasValue",
	"HoleValue",
	"Comment",
	"SingleQuote",
	"DoubleQuote",
	"WhiteSpacing",
	"MustWhiteSpacing",
	"Equal",
	"BlankLine",
	"Whitespace",
	"EndOfLine",
	"EndOfFile",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Peg struct {
	*AST

	Buffer string
	buffer []rune
	rules  [56]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Peg) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Peg) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Peg
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *Peg) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Peg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.addDeclarationIdentifier(text)
		case ruleAction1:
			p.addValue()
		case ruleAction2:
			p.LineDone()
		case ruleAction3:
			p.addAction(text)
		case ruleAction4:
			p.addEntity(text)
		case ruleAction5:
			p.LineDone()
		case ruleAction6:
			p.addParamKey(text)
		case ruleAction7:
			p.addParamHoleValue(text)
		case ruleAction8:
			p.addAliasParam(text)
		case ruleAction9:
			p.addStringValue(text)
		case ruleAction10:
			p.addStringValue(text)
		case ruleAction11:
			p.addParamValue(text)
		case ruleAction12:
			p.addParamRefValue(text)
		case ruleAction13:
			p.addParamRefsListValue(text)
		case ruleAction14:
			p.addParamCidrValue(text)
		case ruleAction15:
			p.addParamIpValue(text)
		case ruleAction16:
			p.addCsvValue(text)
		case ruleAction17:
			p.addParamValue(text)
		case ruleAction18:
			p.LineDone()
		case ruleAction19:
			p.LineDone()

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *Peg) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Script <- <((BlankLine* Statement BlankLine*)+ WhiteSpacing EndOfFile)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
			l4:
				{
					position5, tokenIndex5 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l5
					}
					goto l4
				l5:
					position, tokenIndex = position5, tokenIndex5
				}
				{
					position6 := position
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
					{
						position7, tokenIndex7 := position, tokenIndex
						if !_rules[ruleCmdExpr]() {
							goto l8
						}
						goto l7
					l8:
						position, tokenIndex = position7, tokenIndex7
						{
							position10 := position
							{
								position11 := position
								if !_rules[ruleIdentifier]() {
									goto l9
								}
								add(rulePegText, position11)
							}
							{
								add(ruleAction0, position)
							}
							if !_rules[ruleEqual]() {
								goto l9
							}
							{
								position13, tokenIndex13 := position, tokenIndex
								if !_rules[ruleCmdExpr]() {
									goto l14
								}
								goto l13
							l14:
								position, tokenIndex = position13, tokenIndex13
								{
									position15 := position
									{
										add(ruleAction1, position)
									}
									if !_rules[ruleNoRefValue]() {
										goto l9
									}
									{
										add(ruleAction2, position)
									}
									add(ruleValueExpr, position15)
								}
							}
						l13:
							add(ruleDeclaration, position10)
						}
						goto l7
					l9:
						position, tokenIndex = position7, tokenIndex7
						{
							position18 := position
							{
								position19, tokenIndex19 := position, tokenIndex
								if buffer[position] != rune('#') {
									goto l20
								}
								position++
							l21:
								{
									position22, tokenIndex22 := position, tokenIndex
									{
										position23, tokenIndex23 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l23
										}
										goto l22
									l23:
										position, tokenIndex = position23, tokenIndex23
									}
									if !matchDot() {
										goto l22
									}
									goto l21
								l22:
									position, tokenIndex = position22, tokenIndex22
								}
								goto l19
							l20:
								position, tokenIndex = position19, tokenIndex19
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
							l24:
								{
									position25, tokenIndex25 := position, tokenIndex
									{
										position26, tokenIndex26 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l26
										}
										goto l25
									l26:
										position, tokenIndex = position26, tokenIndex26
									}
									if !matchDot() {
										goto l25
									}
									goto l24
								l25:
									position, tokenIndex = position25, tokenIndex25
								}
								{
									add(ruleAction18, position)
								}
							}
						l19:
							add(ruleComment, position18)
						}
					}
				l7:
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
				l28:
					{
						position29, tokenIndex29 := position, tokenIndex
						if !_rules[ruleEndOfLine]() {
							goto l29
						}
						goto l28
					l29:
						position, tokenIndex = position29, tokenIndex29
					}
					add(ruleStatement, position6)
				}
			l30:
				{
					position31, tokenIndex31 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex = position31, tokenIndex31
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
				l32:
					{
						position33, tokenIndex33 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l33
						}
						goto l32
					l33:
						position, tokenIndex = position33, tokenIndex33
					}
					{
						position34 := position
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
						{
							position35, tokenIndex35 := position, tokenIndex
							if !_rules[ruleCmdExpr]() {
								goto l36
							}
							goto l35
						l36:
							position, tokenIndex = position35, tokenIndex35
							{
								position38 := position
								{
									position39 := position
									if !_rules[ruleIdentifier]() {
										goto l37
									}
									add(rulePegText, position39)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l37
								}
								{
									position41, tokenIndex41 := position, tokenIndex
									if !_rules[ruleCmdExpr]() {
										goto l42
									}
									goto l41
								l42:
									position, tokenIndex = position41, tokenIndex41
									{
										position43 := position
										{
											add(ruleAction1, position)
										}
										if !_rules[ruleNoRefValue]() {
											goto l37
										}
										{
											add(ruleAction2, position)
										}
										add(ruleValueExpr, position43)
									}
								}
							l41:
								add(ruleDeclaration, position38)
							}
							goto l35
						l37:
							position, tokenIndex = position35, tokenIndex35
							{
								position46 := position
								{
									position47, tokenIndex47 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l48
									}
									position++
								l49:
									{
										position50, tokenIndex50 := position, tokenIndex
										{
											position51, tokenIndex51 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l51
											}
											goto l50
										l51:
											position, tokenIndex = position51, tokenIndex51
										}
										if !matchDot() {
											goto l50
										}
										goto l49
									l50:
										position, tokenIndex = position50, tokenIndex50
									}
									goto l47
								l48:
									position, tokenIndex = position47, tokenIndex47
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
								l52:
									{
										position53, tokenIndex53 := position, tokenIndex
										{
											position54, tokenIndex54 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l54
											}
											goto l53
										l54:
											position, tokenIndex = position54, tokenIndex54
										}
										if !matchDot() {
											goto l53
										}
										goto l52
									l53:
										position, tokenIndex = position53, tokenIndex53
									}
									{
										add(ruleAction18, position)
									}
								}
							l47:
								add(ruleComment, position46)
							}
						}
					l35:
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
					l56:
						{
							position57, tokenIndex57 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l57
							}
							goto l56
						l57:
							position, tokenIndex = position57, tokenIndex57
						}
						add(ruleStatement, position34)
					}
				l58:
					{
						position59, tokenIndex59 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l59
						}
						goto l58
					l59:
						position, tokenIndex = position59, tokenIndex59
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l0
				}
				{
					position60 := position
					{
						position61, tokenIndex61 := position, tokenIndex
						if !matchDot() {
							goto l61
						}
						goto l0
					l61:
						position, tokenIndex = position61, tokenIndex61
					}
					add(ruleEndOfFile, position60)
				}
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Statement <- <(WhiteSpacing (CmdExpr / Declaration / Comment) WhiteSpacing EndOfLine*)> */
		nil,
		/* 2 Action <- <[a-z]+> */
		nil,
		/* 3 Entity <- <([a-z] / [0-9])+> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal (CmdExpr / ValueExpr))> */
		nil,
		/* 5 ValueExpr <- <(Action1 NoRefValue Action2)> */
		nil,
		/* 6 CmdExpr <- <(<Action> Action3 MustWhiteSpacing <Entity> Action4 (MustWhiteSpacing Params)? Action5)> */
		func() bool {
			position67, tokenIndex67 := position, tokenIndex
			{
				position68 := position
				{
					position69 := position
					{
						position70 := position
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l67
						}
						position++
					l71:
						{
							position72, tokenIndex72 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l72
							}
							position++
							goto l71
						l72:
							position, tokenIndex = position72, tokenIndex72
						}
						add(ruleAction, position70)
					}
					add(rulePegText, position69)
				}
				{
					add(ruleAction3, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l67
				}
				{
					position74 := position
					{
						position75 := position
						{
							position78, tokenIndex78 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex = position78, tokenIndex78
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l67
							}
							position++
						}
					l78:
					l76:
						{
							position77, tokenIndex77 := position, tokenIndex
							{
								position80, tokenIndex80 := position, tokenIndex
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l81
								}
								position++
								goto l80
							l81:
								position, tokenIndex = position80, tokenIndex80
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l77
								}
								position++
							}
						l80:
							goto l76
						l77:
							position, tokenIndex = position77, tokenIndex77
						}
						add(ruleEntity, position75)
					}
					add(rulePegText, position74)
				}
				{
					add(ruleAction4, position)
				}
				{
					position83, tokenIndex83 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l83
					}
					{
						position85 := position
						{
							position88 := position
							{
								position89 := position
								if !_rules[ruleIdentifier]() {
									goto l83
								}
								add(rulePegText, position89)
							}
							{
								add(ruleAction6, position)
							}
							if !_rules[ruleEqual]() {
								goto l83
							}
							{
								position91 := position
								{
									switch buffer[position] {
									case '[':
										{
											position93 := position
											if buffer[position] != rune('[') {
												goto l83
											}
											position++
											{
												position94 := position
												if !_rules[ruleWhiteSpacing]() {
													goto l83
												}
												if !_rules[ruleRefValue]() {
													goto l83
												}
											l95:
												{
													position96, tokenIndex96 := position, tokenIndex
													if !_rules[ruleRefValue]() {
														goto l96
													}
													goto l95
												l96:
													position, tokenIndex = position96, tokenIndex96
												}
											l97:
												{
													position98, tokenIndex98 := position, tokenIndex
													if buffer[position] != rune(',') {
														goto l98
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l98
													}
													if !_rules[ruleRefValue]() {
														goto l98
													}
													goto l97
												l98:
													position, tokenIndex = position98, tokenIndex98
												}
												add(rulePegText, position94)
											}
											if buffer[position] != rune(']') {
												goto l83
											}
											position++
											add(ruleRefsListValue, position93)
										}
										{
											add(ruleAction13, position)
										}
										break
									case '$':
										if !_rules[ruleRefValue]() {
											goto l83
										}
										{
											add(ruleAction12, position)
										}
										break
									default:
										if !_rules[ruleNoRefValue]() {
											goto l83
										}
										break
									}
								}

								add(ruleValue, position91)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l83
							}
							add(ruleParam, position88)
						}
					l86:
						{
							position87, tokenIndex87 := position, tokenIndex
							{
								position101 := position
								{
									position102 := position
									if !_rules[ruleIdentifier]() {
										goto l87
									}
									add(rulePegText, position102)
								}
								{
									add(ruleAction6, position)
								}
								if !_rules[ruleEqual]() {
									goto l87
								}
								{
									position104 := position
									{
										switch buffer[position] {
										case '[':
											{
												position106 := position
												if buffer[position] != rune('[') {
													goto l87
												}
												position++
												{
													position107 := position
													if !_rules[ruleWhiteSpacing]() {
														goto l87
													}
													if !_rules[ruleRefValue]() {
														goto l87
													}
												l108:
													{
														position109, tokenIndex109 := position, tokenIndex
														if !_rules[ruleRefValue]() {
															goto l109
														}
														goto l108
													l109:
														position, tokenIndex = position109, tokenIndex109
													}
												l110:
													{
														position111, tokenIndex111 := position, tokenIndex
														if buffer[position] != rune(',') {
															goto l111
														}
														position++
														if !_rules[ruleWhiteSpacing]() {
															goto l111
														}
														if !_rules[ruleRefValue]() {
															goto l111
														}
														goto l110
													l111:
														position, tokenIndex = position111, tokenIndex111
													}
													add(rulePegText, position107)
												}
												if buffer[position] != rune(']') {
													goto l87
												}
												position++
												add(ruleRefsListValue, position106)
											}
											{
												add(ruleAction13, position)
											}
											break
										case '$':
											if !_rules[ruleRefValue]() {
												goto l87
											}
											{
												add(ruleAction12, position)
											}
											break
										default:
											if !_rules[ruleNoRefValue]() {
												goto l87
											}
											break
										}
									}

									add(ruleValue, position104)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l87
								}
								add(ruleParam, position101)
							}
							goto l86
						l87:
							position, tokenIndex = position87, tokenIndex87
						}
						add(ruleParams, position85)
					}
					goto l84
				l83:
					position, tokenIndex = position83, tokenIndex83
				}
			l84:
				{
					add(ruleAction5, position)
				}
				add(ruleCmdExpr, position68)
			}
			return true
		l67:
			position, tokenIndex = position67, tokenIndex67
			return false
		},
		/* 7 Params <- <Param+> */
		nil,
		/* 8 Param <- <(<Identifier> Action6 Equal Value WhiteSpacing)> */
		nil,
		/* 9 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position117, tokenIndex117 := position, tokenIndex
			{
				position118 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l117
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l117
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l117
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l117
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l117
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l117
						}
						position++
						break
					}
				}

			l119:
				{
					position120, tokenIndex120 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l120
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l120
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l120
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l120
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l120
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l120
							}
							position++
							break
						}
					}

					goto l119
				l120:
					position, tokenIndex = position120, tokenIndex120
				}
				add(ruleIdentifier, position118)
			}
			return true
		l117:
			position, tokenIndex = position117, tokenIndex117
			return false
		},
		/* 10 NoRefValue <- <((AliasValue Action8) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / ((&('\'') (SingleQuote <SingleQuotedValue> Action10 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action9 DoubleQuote)) | (&('{') (HoleValue Action7)) | (&('*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<OtherParamValue> Action11))))> */
		func() bool {
			position123, tokenIndex123 := position, tokenIndex
			{
				position124 := position
				{
					position125, tokenIndex125 := position, tokenIndex
					{
						position127 := position
						{
							position128, tokenIndex128 := position, tokenIndex
							if buffer[position] != rune('@') {
								goto l129
							}
							position++
							{
								position130 := position
								if !_rules[ruleOtherParamValue]() {
									goto l129
								}
								add(rulePegText, position130)
							}
							goto l128
						l129:
							position, tokenIndex = position128, tokenIndex128
							if buffer[position] != rune('@') {
								goto l131
							}
							position++
							if !_rules[ruleDoubleQuote]() {
								goto l131
							}
							{
								position132 := position
								if !_rules[ruleDoubleQuotedValue]() {
									goto l131
								}
								add(rulePegText, position132)
							}
							if !_rules[ruleDoubleQuote]() {
								goto l131
							}
							goto l128
						l131:
							position, tokenIndex = position128, tokenIndex128
							if buffer[position] != rune('@') {
								goto l126
							}
							position++
							if !_rules[ruleSingleQuote]() {
								goto l126
							}
							{
								position133 := position
								if !_rules[ruleSingleQuotedValue]() {
									goto l126
								}
								add(rulePegText, position133)
							}
							if !_rules[ruleSingleQuote]() {
								goto l126
							}
						}
					l128:
						add(ruleAliasValue, position127)
					}
					{
						add(ruleAction8, position)
					}
					goto l125
				l126:
					position, tokenIndex = position125, tokenIndex125
					if !_rules[ruleDoubleQuote]() {
						goto l135
					}
					if !_rules[ruleCustomTypedValue]() {
						goto l135
					}
					if !_rules[ruleDoubleQuote]() {
						goto l135
					}
					goto l125
				l135:
					position, tokenIndex = position125, tokenIndex125
					if !_rules[ruleSingleQuote]() {
						goto l136
					}
					if !_rules[ruleCustomTypedValue]() {
						goto l136
					}
					if !_rules[ruleSingleQuote]() {
						goto l136
					}
					goto l125
				l136:
					position, tokenIndex = position125, tokenIndex125
					if !_rules[ruleCustomTypedValue]() {
						goto l137
					}
					goto l125
				l137:
					position, tokenIndex = position125, tokenIndex125
					{
						switch buffer[position] {
						case '\'':
							if !_rules[ruleSingleQuote]() {
								goto l123
							}
							{
								position139 := position
								if !_rules[ruleSingleQuotedValue]() {
									goto l123
								}
								add(rulePegText, position139)
							}
							{
								add(ruleAction10, position)
							}
							if !_rules[ruleSingleQuote]() {
								goto l123
							}
							break
						case '"':
							if !_rules[ruleDoubleQuote]() {
								goto l123
							}
							{
								position141 := position
								if !_rules[ruleDoubleQuotedValue]() {
									goto l123
								}
								add(rulePegText, position141)
							}
							{
								add(ruleAction9, position)
							}
							if !_rules[ruleDoubleQuote]() {
								goto l123
							}
							break
						case '{':
							{
								position143 := position
								if buffer[position] != rune('{') {
									goto l123
								}
								position++
								if !_rules[ruleWhiteSpacing]() {
									goto l123
								}
								{
									position144 := position
									if !_rules[ruleIdentifier]() {
										goto l123
									}
									add(rulePegText, position144)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l123
								}
								if buffer[position] != rune('}') {
									goto l123
								}
								position++
								add(ruleHoleValue, position143)
							}
							{
								add(ruleAction7, position)
							}
							break
						default:
							{
								position146 := position
								if !_rules[ruleOtherParamValue]() {
									goto l123
								}
								add(rulePegText, position146)
							}
							{
								add(ruleAction11, position)
							}
							break
						}
					}

				}
			l125:
				add(ruleNoRefValue, position124)
			}
			return true
		l123:
			position, tokenIndex = position123, tokenIndex123
			return false
		},
		/* 11 Value <- <((&('[') (RefsListValue Action13)) | (&('$') (RefValue Action12)) | (&('"' | '\'' | '*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '{' | '~') NoRefValue))> */
		nil,
		/* 12 CustomTypedValue <- <((<CidrValue> Action14) / (<IpValue> Action15) / (<CSVValue> Action16) / (<IntRangeValue> Action17))> */
		func() bool {
			position149, tokenIndex149 := position, tokenIndex
			{
				position150 := position
				{
					position151, tokenIndex151 := position, tokenIndex
					{
						position153 := position
						{
							position154 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l152
							}
							position++
						l155:
							{
								position156, tokenIndex156 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l156
								}
								position++
								goto l155
							l156:
								position, tokenIndex = position156, tokenIndex156
							}
							if buffer[position] != rune('.') {
								goto l152
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l152
							}
							position++
						l157:
							{
								position158, tokenIndex158 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l158
								}
								position++
								goto l157
							l158:
								position, tokenIndex = position158, tokenIndex158
							}
							if buffer[position] != rune('.') {
								goto l152
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l152
							}
							position++
						l159:
							{
								position160, tokenIndex160 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l160
								}
								position++
								goto l159
							l160:
								position, tokenIndex = position160, tokenIndex160
							}
							if buffer[position] != rune('.') {
								goto l152
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l152
							}
							position++
						l161:
							{
								position162, tokenIndex162 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l162
								}
								position++
								goto l161
							l162:
								position, tokenIndex = position162, tokenIndex162
							}
							if buffer[position] != rune('/') {
								goto l152
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l152
							}
							position++
						l163:
							{
								position164, tokenIndex164 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l164
								}
								position++
								goto l163
							l164:
								position, tokenIndex = position164, tokenIndex164
							}
							add(ruleCidrValue, position154)
						}
						add(rulePegText, position153)
					}
					{
						add(ruleAction14, position)
					}
					goto l151
				l152:
					position, tokenIndex = position151, tokenIndex151
					{
						position167 := position
						{
							position168 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l166
							}
							position++
						l169:
							{
								position170, tokenIndex170 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l170
								}
								position++
								goto l169
							l170:
								position, tokenIndex = position170, tokenIndex170
							}
							if buffer[position] != rune('.') {
								goto l166
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l166
							}
							position++
						l171:
							{
								position172, tokenIndex172 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l172
								}
								position++
								goto l171
							l172:
								position, tokenIndex = position172, tokenIndex172
							}
							if buffer[position] != rune('.') {
								goto l166
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l166
							}
							position++
						l173:
							{
								position174, tokenIndex174 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l174
								}
								position++
								goto l173
							l174:
								position, tokenIndex = position174, tokenIndex174
							}
							if buffer[position] != rune('.') {
								goto l166
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l166
							}
							position++
						l175:
							{
								position176, tokenIndex176 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l176
								}
								position++
								goto l175
							l176:
								position, tokenIndex = position176, tokenIndex176
							}
							add(ruleIpValue, position168)
						}
						add(rulePegText, position167)
					}
					{
						add(ruleAction15, position)
					}
					goto l151
				l166:
					position, tokenIndex = position151, tokenIndex151
					{
						position179 := position
						{
							position180 := position
							if !_rules[ruleOtherParamValue]() {
								goto l178
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l178
							}
							if buffer[position] != rune(',') {
								goto l178
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l178
							}
						l181:
							{
								position182, tokenIndex182 := position, tokenIndex
								if !_rules[ruleOtherParamValue]() {
									goto l182
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l182
								}
								if buffer[position] != rune(',') {
									goto l182
								}
								position++
								if !_rules[ruleWhiteSpacing]() {
									goto l182
								}
								goto l181
							l182:
								position, tokenIndex = position182, tokenIndex182
							}
							if !_rules[ruleOtherParamValue]() {
								goto l178
							}
							add(ruleCSVValue, position180)
						}
						add(rulePegText, position179)
					}
					{
						add(ruleAction16, position)
					}
					goto l151
				l178:
					position, tokenIndex = position151, tokenIndex151
					{
						position184 := position
						{
							position185 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l149
							}
							position++
						l186:
							{
								position187, tokenIndex187 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l187
								}
								position++
								goto l186
							l187:
								position, tokenIndex = position187, tokenIndex187
							}
							if buffer[position] != rune('-') {
								goto l149
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l149
							}
							position++
						l188:
							{
								position189, tokenIndex189 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l189
								}
								position++
								goto l188
							l189:
								position, tokenIndex = position189, tokenIndex189
							}
							add(ruleIntRangeValue, position185)
						}
						add(rulePegText, position184)
					}
					{
						add(ruleAction17, position)
					}
				}
			l151:
				add(ruleCustomTypedValue, position150)
			}
			return true
		l149:
			position, tokenIndex = position149, tokenIndex149
			return false
		},
		/* 13 OtherParamValue <- <((&('*') '*') | (&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position191, tokenIndex191 := position, tokenIndex
			{
				position192 := position
				{
					switch buffer[position] {
					case '*':
						if buffer[position] != rune('*') {
							goto l191
						}
						position++
						break
					case '>':
						if buffer[position] != rune('>') {
							goto l191
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l191
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l191
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l191
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l191
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l191
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l191
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l191
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l191
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l191
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l191
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l191
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l191
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l191
						}
						position++
						break
					}
				}

			l193:
				{
					position194, tokenIndex194 := position, tokenIndex
					{
						switch buffer[position] {
						case '*':
							if buffer[position] != rune('*') {
								goto l194
							}
							position++
							break
						case '>':
							if buffer[position] != rune('>') {
								goto l194
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l194
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l194
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l194
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l194
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l194
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l194
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l194
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l194
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l194
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l194
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l194
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l194
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l194
							}
							position++
							break
						}
					}

					goto l193
				l194:
					position, tokenIndex = position194, tokenIndex194
				}
				add(ruleOtherParamValue, position192)
			}
			return true
		l191:
			position, tokenIndex = position191, tokenIndex191
			return false
		},
		/* 14 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position198 := position
			l199:
				{
					position200, tokenIndex200 := position, tokenIndex
					{
						position201, tokenIndex201 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l201
						}
						position++
						goto l200
					l201:
						position, tokenIndex = position201, tokenIndex201
					}
					if !matchDot() {
						goto l200
					}
					goto l199
				l200:
					position, tokenIndex = position200, tokenIndex200
				}
				add(ruleDoubleQuotedValue, position198)
			}
			return true
		},
		/* 15 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position203 := position
			l204:
				{
					position205, tokenIndex205 := position, tokenIndex
					{
						position206, tokenIndex206 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l206
						}
						position++
						goto l205
					l206:
						position, tokenIndex = position206, tokenIndex206
					}
					if !matchDot() {
						goto l205
					}
					goto l204
				l205:
					position, tokenIndex = position205, tokenIndex205
				}
				add(ruleSingleQuotedValue, position203)
			}
			return true
		},
		/* 16 CSVValue <- <((OtherParamValue WhiteSpacing ',' WhiteSpacing)+ OtherParamValue)> */
		nil,
		/* 17 CidrValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+ '/' [0-9]+)> */
		nil,
		/* 18 IpValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+)> */
		nil,
		/* 19 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 20 RefValue <- <('$' <Identifier>)> */
		func() bool {
			position211, tokenIndex211 := position, tokenIndex
			{
				position212 := position
				if buffer[position] != rune('$') {
					goto l211
				}
				position++
				{
					position213 := position
					if !_rules[ruleIdentifier]() {
						goto l211
					}
					add(rulePegText, position213)
				}
				add(ruleRefValue, position212)
			}
			return true
		l211:
			position, tokenIndex = position211, tokenIndex211
			return false
		},
		/* 21 RefsListValue <- <('[' <(WhiteSpacing RefValue+ (',' WhiteSpacing RefValue)*)> ']')> */
		nil,
		/* 22 AliasValue <- <(('@' <OtherParamValue>) / ('@' DoubleQuote <DoubleQuotedValue> DoubleQuote) / ('@' SingleQuote <SingleQuotedValue> SingleQuote))> */
		nil,
		/* 23 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 24 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action18))> */
		nil,
		/* 25 SingleQuote <- <'\''> */
		func() bool {
			position218, tokenIndex218 := position, tokenIndex
			{
				position219 := position
				if buffer[position] != rune('\'') {
					goto l218
				}
				position++
				add(ruleSingleQuote, position219)
			}
			return true
		l218:
			position, tokenIndex = position218, tokenIndex218
			return false
		},
		/* 26 DoubleQuote <- <'"'> */
		func() bool {
			position220, tokenIndex220 := position, tokenIndex
			{
				position221 := position
				if buffer[position] != rune('"') {
					goto l220
				}
				position++
				add(ruleDoubleQuote, position221)
			}
			return true
		l220:
			position, tokenIndex = position220, tokenIndex220
			return false
		},
		/* 27 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position223 := position
			l224:
				{
					position225, tokenIndex225 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l225
					}
					goto l224
				l225:
					position, tokenIndex = position225, tokenIndex225
				}
				add(ruleWhiteSpacing, position223)
			}
			return true
		},
		/* 28 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position226, tokenIndex226 := position, tokenIndex
			{
				position227 := position
				if !_rules[ruleWhitespace]() {
					goto l226
				}
			l228:
				{
					position229, tokenIndex229 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l229
					}
					goto l228
				l229:
					position, tokenIndex = position229, tokenIndex229
				}
				add(ruleMustWhiteSpacing, position227)
			}
			return true
		l226:
			position, tokenIndex = position226, tokenIndex226
			return false
		},
		/* 29 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position230, tokenIndex230 := position, tokenIndex
			{
				position231 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l230
				}
				if buffer[position] != rune('=') {
					goto l230
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l230
				}
				add(ruleEqual, position231)
			}
			return true
		l230:
			position, tokenIndex = position230, tokenIndex230
			return false
		},
		/* 30 BlankLine <- <(WhiteSpacing EndOfLine Action19)> */
		func() bool {
			position232, tokenIndex232 := position, tokenIndex
			{
				position233 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l232
				}
				if !_rules[ruleEndOfLine]() {
					goto l232
				}
				{
					add(ruleAction19, position)
				}
				add(ruleBlankLine, position233)
			}
			return true
		l232:
			position, tokenIndex = position232, tokenIndex232
			return false
		},
		/* 31 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position235, tokenIndex235 := position, tokenIndex
			{
				position236 := position
				{
					position237, tokenIndex237 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex = position237, tokenIndex237
					if buffer[position] != rune('\t') {
						goto l235
					}
					position++
				}
			l237:
				add(ruleWhitespace, position236)
			}
			return true
		l235:
			position, tokenIndex = position235, tokenIndex235
			return false
		},
		/* 32 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position239, tokenIndex239 := position, tokenIndex
			{
				position240 := position
				{
					position241, tokenIndex241 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l242
					}
					position++
					if buffer[position] != rune('\n') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex = position241, tokenIndex241
					if buffer[position] != rune('\n') {
						goto l243
					}
					position++
					goto l241
				l243:
					position, tokenIndex = position241, tokenIndex241
					if buffer[position] != rune('\r') {
						goto l239
					}
					position++
				}
			l241:
				add(ruleEndOfLine, position240)
			}
			return true
		l239:
			position, tokenIndex = position239, tokenIndex239
			return false
		},
		/* 33 EndOfFile <- <!.> */
		nil,
		nil,
		/* 36 Action0 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 37 Action1 <- <{ p.addValue() }> */
		nil,
		/* 38 Action2 <- <{ p.LineDone() }> */
		nil,
		/* 39 Action3 <- <{ p.addAction(text) }> */
		nil,
		/* 40 Action4 <- <{ p.addEntity(text) }> */
		nil,
		/* 41 Action5 <- <{ p.LineDone() }> */
		nil,
		/* 42 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 43 Action7 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 44 Action8 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 45 Action9 <- <{ p.addStringValue(text) }> */
		nil,
		/* 46 Action10 <- <{ p.addStringValue(text) }> */
		nil,
		/* 47 Action11 <- <{ p.addParamValue(text) }> */
		nil,
		/* 48 Action12 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 49 Action13 <- <{ p.addParamRefsListValue(text) }> */
		nil,
		/* 50 Action14 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 51 Action15 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 52 Action16 <- <{p.addCsvValue(text)}> */
		nil,
		/* 53 Action17 <- <{ p.addParamValue(text) }> */
		nil,
		/* 54 Action18 <- <{ p.LineDone() }> */
		nil,
		/* 55 Action19 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
