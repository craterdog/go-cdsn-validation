/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn

import (
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	sts "strings"
	uni "unicode"
)

// PARSER INTERFACE

// This function parses the specified document source retrieved from a POSIX
// compliant file and returns the corresponding CDSN grammar that was used
// to generate the document using the CDSN formatting capabilities.
// A POSIX compliant file must end with an EOF marker.
func ParseDocument(source []byte) GrammarLike {
	var ok bool
	var token *Token
	var grammar GrammarLike
	var tokens = make(chan Token, 256)
	ScanTokens(source, tokens) // Starts scanning in a separate go routine.
	var p = &parser{
		symbols: col.Catalog[SYMBOL, DefinitionLike](),
		source:  source,
		next:    col.StackWithCapacity[*Token](4),
		tokens:  tokens,
	}
	grammar, token, ok = p.parseGrammar()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("statement",
			"$grammar",
			"$statement")
		panic(message)
	}
	var iterator = col.Iterator[col.Binding[SYMBOL, DefinitionLike]](p.symbols)
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var symbol = association.GetKey()
		var definition = association.GetValue()
		if definition == nil {
			panic(fmt.Sprintf("Missing a definition for symbol: %v\n", symbol))
		}
	}
	return grammar
}

// PARSER IMPLEMENTATION

// This map captures the syntax expressions for Crater Dog Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar_ = map[string]string{
	"$grammar":     `<statement> EOF  ! Terminated with an end-of-file marker.`,
	"$statement":   `COMMENT | definition`,
	"$definition":  `SYMBOL ":" expression  ! This works for both tokens and rules.`,
	"$expression":  `alternative {"|" alternative}`,
	"$alternative": `<factor> [NOTE]`,
	"$factor":      `element | range | inverse | grouping`,
	"$element":     `INTRINSIC | STRING | NUMBER | NAME`,
	"$range":       `CHARACTER [".." CHARACTER]  ! A range of CHARACTERs is inclusive.`,
	"$inverse":     `"~" factor`,
	"$grouping":    `exactlyN | zeroOrOne | zeroOrMore | oneOrMore`,
	"$exactlyN":    `"(" expression ")" [NUMBER]  ! The default is exactly one.`,
	"$zeroOrOne":   `"[" expression "]"`,
	"$zeroOrMore":  `"{" expression "}"`,
	"$oneOrMore":   `"<" expression ">"`,
}

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar_[symbol])
	}
	return message
}

// This type defines the structure and methods for the parser agent.
type parser struct {
	symbols        col.CatalogLike[SYMBOL, DefinitionLike]
	source         []byte
	next           col.StackLike[*Token] // The stack of the retrieved tokens that have been put back.
	tokens         chan Token            // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 *Token                // The previous four tokens that have been retrieved.
	isToken        bool                  // Whether or not the current definition is a token definition.
}

// This method puts back the current token onto the token stream so that it can
// be retrieved by another parsing method.
func (v *parser) backupOne() {
	v.next.AddValue(v.p1)
	v.p1, v.p2, v.p3, v.p4 = v.p2, v.p3, v.p4, nil
}

// This method returns an error message containing the context for a parsing
// error.
func (v *parser) formatError(token *Token) string {
	var message = fmt.Sprintf("An unexpected token was received by the parser: %v\n", token)
	var line = token.Line
	var lines = sts.Split(string(v.source), EOL)

	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + EOL
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + EOL

	message += " \033[32m>>>─"
	var count = 0
	for count < token.Position {
		message += "─"
		count++
	}
	message += "⌃\033[36m\n"

	if line < len(lines) {
		message += fmt.Sprintf("%04d: ", line+1) + string(lines[line]) + EOL
	}
	message += "\033[0m\n"

	return message
}

// This method attempts to read the next token from the token stream and return
// it.
func (v *parser) nextToken() *Token {
	var next *Token
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF or error token.")
		}
		next = &token
		if next.Type == TokenERROR {
			var message = v.formatError(next)
			panic(message)
		}
	} else {
		next = v.next.RemoveTop()
	}
	v.p4, v.p3, v.p2, v.p1 = v.p3, v.p2, v.p1, next
	return next
}

// This method attempts to parse an alternative. It returns the alternative and
// whether or not the alternative was successfully parsed.
func (v *parser) parseAlternative() (AlternativeLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var factors = col.List[Factor]()
	var note NOTE
	var alternative AlternativeLike
	factor, token, ok = v.parseFactor()
	if !ok {
		// An alternative must have at least one factor.
		return alternative, token, false
	}
	for {
		factors.AddValue(factor)
		factor, token, ok = v.parseFactor()
		if !ok {
			// No more factors.
			break
		}
	}
	note, _, _ = v.parseNOTE() // The note is optional.
	alternative = Alternative(factors, note)
	return alternative, token, true
}

// This method attempts to parse a character. It returns the character and
// whether or not a character was successfully parsed.
func (v *parser) parseCHARACTER() (CHARACTER, *Token, bool) {
	var character CHARACTER
	var token = v.nextToken()
	if token.Type != TokenCHARACTER {
		v.backupOne()
		return character, token, false
	}
	character = CHARACTER(token.Value)
	return character, token, true
}

// This method attempts to parse a comment. It returns the comment and whether
// or not a comment was successfully parsed.
func (v *parser) parseCOMMENT() (COMMENT, *Token, bool) {
	var comment COMMENT
	var token = v.nextToken()
	if token.Type != TokenCOMMENT {
		v.backupOne()
		return comment, token, false
	}
	comment = COMMENT(token.Value)
	return comment, token, true
}

// This method attempts to parse a definition. It returns the definition and
// whether or not the definition was successfully parsed.
func (v *parser) parseDefinition() (DefinitionLike, *Token, bool) {
	var ok bool
	var token *Token
	var symbol SYMBOL
	var expression ExpressionLike
	var definition DefinitionLike
	symbol, token, ok = v.parseSYMBOL()
	if !ok {
		// This is not a definition.
		return definition, token, false
	}
	v.isToken = uni.IsUpper(rune(symbol[1]))
	_, token, ok = v.parseLITERAL(":")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(":",
			"$definition",
			"$SYMBOL",
			"$expression")
		panic(message)
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$definition",
			"$SYMBOL",
			"$expression")
		panic(message)
	}
	definition = Definition(symbol, expression)
	v.symbols.SetValue(symbol, definition)
	return definition, token, true
}

// This method attempts to parse an element. It returns the element and whether
// or not the element was successfully parsed.
func (v *parser) parseElement() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseINTRINSIC()
	if !ok {
		factor, token, ok = v.parseSTRING()
	}
	if !ok {
		factor, token, ok = v.parseNUMBER()
	}
	if !ok {
		factor, token, ok = v.parseNAME()
	}
	return factor, token, ok
}

// This method attempts to parse an exactly N grouping. It returns the exactly
// N grouping and whether or not the exactly N grouping was successfully parsed.
func (v *parser) parseExactlyN() (ExactlyNLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var exactlyN ExactlyNLike
	_, token, ok = v.parseLITERAL("(")
	if !ok {
		// This is not an exactly N grouping.
		return exactlyN, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$exactlyN",
			"$expression")
		panic(message)
	}
	expression.SetMultilined(false)
	_, token, ok = v.parseLITERAL(")")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(")",
			"$exactlyN",
			"$expression")
		panic(message)
	}
	var n, _, _ = v.parseNUMBER() // The number is optional.
	exactlyN = ExactlyN(expression, n)
	return exactlyN, token, true
}

// This method attempts to parse an expression. It returns the expression and
// whether or not the expression was successfully parsed.
func (v *parser) parseExpression() (ExpressionLike, *Token, bool) {
	var ok bool
	var token *Token
	var alternative AlternativeLike
	var alternatives = col.List[AlternativeLike]()
	var expression ExpressionLike
	alternative, token, ok = v.parseAlternative()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("alternative",
			"$expression",
			"$alternative")
		panic(message)
	}
	for {
		alternatives.AddValue(alternative)
		_, _, ok = v.parseLITERAL("|")
		if !ok {
			// No more alternatives.
			break
		}
		alternative, token, ok = v.parseAlternative()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("alternative",
				"$expression",
				"$alternative")
			panic(message)
		}
	}
	expression = Expression(alternatives)
	return expression, token, true
}

// This method attempts to parse the end-of-file (EOF) marker. It returns
// the token and whether or not an EOF marker was found. Note that the POSIX
// standard requires that the last byte in a file be an end-of-line (EOL)
// character.
func (v *parser) parseEOF() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOF {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// This method attempts to parse a factor. It returns the factor and whether or
// not the factor was successfully parsed.
func (v *parser) parseFactor() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseElement()
	if !ok {
		factor, token, ok = v.parseRange()
	}
	if !ok {
		factor, token, ok = v.parseInverse()
	}
	if !ok {
		factor, token, ok = v.parseGrouping()
	}
	return factor, token, ok
}

// This method attempts to parse a grammar. It returns the grammar and whether
// or not the grammar was successfully parsed.
func (v *parser) parseGrammar() (GrammarLike, *Token, bool) {
	var ok bool
	var token *Token
	var statement StatementLike
	var statements = col.List[StatementLike]()
	var grammar GrammarLike
	statement, token, ok = v.parseStatement()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("statement",
			"$grammar",
			"$statement")
		panic(message)
	}
	for {
		statements.AddValue(statement)
		statement, _, ok = v.parseStatement()
		if !ok {
			// No more statements.
			break
		}
	}
	_, token, ok = v.parseEOF()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("EOF",
			"$grammar",
			"$statement")
		panic(message)
	}
	grammar = Grammar(statements)
	return grammar, token, true
}

// This method attempts to parse a grouping. It returns the grouping and whether
// or not the grouping was successfully parsed.
func (v *parser) parseGrouping() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseExactlyN()
	if !ok {
		factor, token, ok = v.parseZeroOrOne()
	}
	if !ok {
		factor, token, ok = v.parseZeroOrMore()
	}
	if !ok {
		factor, token, ok = v.parseOneOrMore()
	}
	return factor, token, ok
}

// This method attempts to parse an intrinsic. It returns the intrinsic and
// whether or not the intrinsic was successfully parsed.
func (v *parser) parseINTRINSIC() (INTRINSIC, *Token, bool) {
	var intrinsic INTRINSIC
	var token = v.nextToken()
	if token.Type != TokenINTRINSIC {
		v.backupOne()
		return intrinsic, token, false
	}
	intrinsic = INTRINSIC(token.Value)
	return intrinsic, token, true
}

// This method attempts to parse an inverse. It returns the inverse and
// whether or not the inverse was successfully parsed.
func (v *parser) parseInverse() (InverseLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var inverse InverseLike
	_, token, ok = v.parseLITERAL("~")
	if !ok {
		// This is not an inverse.
		return inverse, token, false
	}
	factor, token, ok = v.parseFactor()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("factor",
			"$inverse",
			"$factor")
		panic(message)
	}
	inverse = Inverse(factor)
	return inverse, token, true
}

// This method attempts to parse the specified literal. It returns
// the token and whether or not the literal was found.
func (v *parser) parseLITERAL(literal string) (string, *Token, bool) {
	var token = v.nextToken()
	if token.Type == TokenEOF || token.Value != literal {
		v.backupOne()
		return literal, token, false
	}
	return literal, token, true
}

// This method attempts to parse a name token. It returns the token and
// whether or not a name token was found.
func (v *parser) parseNAME() (NAME, *Token, bool) {
	var name NAME
	var token = v.nextToken()
	if token.Type != TokenNAME {
		v.backupOne()
		return name, token, false
	}
	if v.isToken && uni.IsLower(rune(token.Value[0])) {
		panic(fmt.Sprintf("A token definition contains a rulename: %v\n", token.Value))
	}
	name = NAME(token.Value)
	var symbol = SYMBOL("$" + token.Value)
	var definition = v.symbols.GetValue(symbol)
	v.symbols.SetValue(symbol, definition)
	return name, token, true
}

// This method attempts to parse a note. It returns the note and whether or not
// the note was successfully parsed.
func (v *parser) parseNOTE() (NOTE, *Token, bool) {
	var note NOTE
	var token = v.nextToken()
	if token.Type != TokenNOTE {
		v.backupOne()
		return note, token, false
	}
	note = NOTE(token.Value)
	return note, token, true
}

// This method attempts to parse a number. It returns the number and
// whether or not a number was successfully parsed.
func (v *parser) parseNUMBER() (NUMBER, *Token, bool) {
	var number NUMBER
	var token = v.nextToken()
	if token.Type != TokenNUMBER {
		v.backupOne()
		return number, token, false
	}
	number = NUMBER(token.Value)
	return number, token, true
}

// This method attempts to parse an one or more grouping. It returns the one or
// more grouping and whether or not the one or more grouping was successfully parsed.
func (v *parser) parseOneOrMore() (OneOrMoreLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var oneOrMore OneOrMoreLike
	_, token, ok = v.parseLITERAL("<")
	if !ok {
		// This is not an one or more grouping.
		return oneOrMore, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$oneOrMore",
			"$expression")
		panic(message)
	}
	expression.SetMultilined(false)
	_, token, ok = v.parseLITERAL(">")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(">",
			"$oneOrMore",
			"$expression")
		panic(message)
	}
	oneOrMore = OneOrMore(expression)
	return oneOrMore, token, true
}

// This method attempts to parse a range. It returns the range and whether or
// not the range was successfully parsed.
func (v *parser) parseRange() (RangeLike, *Token, bool) {
	var ok bool
	var token *Token
	var first CHARACTER
	var last CHARACTER
	var range_ RangeLike
	first, token, ok = v.parseCHARACTER()
	if !ok {
		// This is not a range.
		return range_, token, false
	}
	_, _, ok = v.parseLITERAL("..")
	if ok {
		last, token, ok = v.parseCHARACTER()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("CHARACTER",
				"$range",
				"$CHARACTER")
			panic(message)
		}
	}
	range_ = Range(first, last)
	return range_, token, true
}

// This method attempts to parse a statement. It returns the statement and
// whether or not the statement was successfully parsed.
func (v *parser) parseStatement() (StatementLike, *Token, bool) {
	var ok bool
	var token *Token
	var comment COMMENT
	var definition DefinitionLike
	var statement StatementLike
	comment, _, ok = v.parseCOMMENT()
	if !ok {
		definition, token, ok = v.parseDefinition()
		if !ok {
			// This is not a statement.
			return statement, token, false
		}
	}
	statement = Statement(comment, definition)
	return statement, token, true
}

// This method attempts to parse a string. It returns the string and whether
// or not the string was successfully parsed.
func (v *parser) parseSTRING() (STRING, *Token, bool) {
	var string_ STRING
	var token = v.nextToken()
	if token.Type != TokenSTRING {
		v.backupOne()
		return string_, token, false
	}
	string_ = STRING(token.Value)
	return string_, token, true
}

// This method attempts to parse a symbol. It returns the symbol and
// whether or not the symbol was successfully parsed.
func (v *parser) parseSYMBOL() (SYMBOL, *Token, bool) {
	var symbol SYMBOL
	var token = v.nextToken()
	if token.Type != TokenSYMBOL {
		v.backupOne()
		return symbol, token, false
	}
	symbol = SYMBOL(token.Value)
	return symbol, token, true
}

// This method attempts to parse an zero or more grouping. It returns the zero or
// more grouping and whether or not the zero or more grouping was successfully parsed.
func (v *parser) parseZeroOrMore() (ZeroOrMoreLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var zeroOrMore ZeroOrMoreLike
	_, token, ok = v.parseLITERAL("{")
	if !ok {
		// This is not an zero or more grouping.
		return zeroOrMore, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$zeroOrMore",
			"$expression")
		panic(message)
	}
	expression.SetMultilined(false)
	_, token, ok = v.parseLITERAL("}")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("}",
			"$zeroOrMore",
			"$expression")
		panic(message)
	}
	zeroOrMore = ZeroOrMore(expression)
	return zeroOrMore, token, true
}

// This method attempts to parse an zero or more grouping. It returns the zero or
// more grouping and whether or not the zero or more grouping was successfully parsed.
func (v *parser) parseZeroOrOne() (ZeroOrOneLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var zeroOrOne ZeroOrOneLike
	_, token, ok = v.parseLITERAL("[")
	if !ok {
		// This is not an zero or more grouping.
		return zeroOrOne, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$zeroOrOne",
			"$expression")
		panic(message)
	}
	expression.SetMultilined(false)
	_, token, ok = v.parseLITERAL("]")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("]",
			"$zeroOrOne",
			"$expression")
		panic(message)
	}
	zeroOrOne = ZeroOrOne(expression)
	return zeroOrOne, token, true
}
