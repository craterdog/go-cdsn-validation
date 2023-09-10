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
	Scanner(source, tokens) // Starts scanning in a separate go routine.
	var p = &parser{
		source: source,
		next:   col.StackWithCapacity[*Token](4),
		tokens: tokens,
	}
	grammar, token, ok = p.parseGrammar()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("statement",
			"$grammar",
			"$statement")
		panic(message)
	}
	return grammar
}

// PARSER IMPLEMENTATION

// This type defines the structure and methods for the parser agent.
type parser struct {
	source         []byte
	next           col.StackLike[*Token] // The stack of the retrieved tokens that have been put back.
	tokens         chan Token            // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 *Token                // The previous four tokens that have been retrieved.
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
		if next.Type == TokenError {
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
	var note Note
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
	note, _, _ = v.parseNote() // The note is optional.
	alternative = Alternative(factors, note)
	return alternative, token, true
}

// This method attempts to parse an annotation. It returns the annotation and
// whether or not the annotation was successfully parsed.
func (v *parser) parseAnnotation() (Annotation, *Token, bool) {
	var ok bool
	var token *Token
	var note Note
	var comment Comment
	var annotation Annotation
	note, _, ok = v.parseNote()
	if !ok {
		comment, token, ok = v.parseComment()
		if !ok {
			// This is not an annotation.
			return annotation, token, false
		}
		annotation = Annotation(string(comment))
	} else {
		annotation = Annotation(string(note))
	}
	return annotation, token, true
}

// This method attempts to parse a comment. It returns the comment and whether
// or not a comment was successfully parsed.
func (v *parser) parseComment() (Comment, *Token, bool) {
	var comment Comment
	var token = v.nextToken()
	if token.Type != TokenComment {
		v.backupOne()
		return comment, token, false
	}
	comment = Comment(token.Value)
	return comment, token, true
}

// This method attempts to parse a definition. It returns the definition and
// whether or not the definition was successfully parsed.
func (v *parser) parseDefinition() (DefinitionLike, *Token, bool) {
	var ok bool
	var token *Token
	var alternative AlternativeLike
	var alternatives = col.List[AlternativeLike]()
	var definition DefinitionLike
	v.parseEOL() // The EOL is optional.
	alternative, token, ok = v.parseAlternative()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("alternative",
			"$definition",
			"$alternative")
		panic(message)
	}
	for {
		alternatives.AddValue(alternative)
		v.parseEOL() // The EOL is optional.
		_, _, ok = v.parseDelimiter("|")
		if !ok {
			// No more alternatives.
			break
		}
		alternative, token, ok = v.parseAlternative()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("alternative",
				"$definition",
				"$alternative")
			panic(message)
		}
	}
	definition = Definition(alternatives)
	return definition, token, true
}

// This method attempts to parse the specified delimiter. It returns
// the token and whether or not the delimiter was found.
func (v *parser) parseDelimiter(delimiter string) (string, *Token, bool) {
	var token = v.nextToken()
	if token.Type == TokenEOF || token.Value != delimiter {
		v.backupOne()
		return delimiter, token, false
	}
	return delimiter, token, true
}

// This method attempts to parse an exact count group. It returns the
// exact count group and whether or not the exact count group was
// successfully parsed.
func (v *parser) parseExactlyN() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var definition DefinitionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("(")
	if !ok {
		// This is not a precedence group.
		return group, token, false
	}
	definition, token, ok = v.parseDefinition()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("definition",
			"$factor",
			"$definition")
		panic(message)
	}
	definition.SetMultilined(false)
	_, token, ok = v.parseDelimiter(")")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(")",
			"$factor",
			"$definition")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(definition, ExactlyN, number)
	return group, token, true
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

// This method attempts to parse the end-of-line (EOL) token. It returns
// the token and whether or not an EOF token was found.
func (v *parser) parseEOL() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOL {
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
	factor, token, ok = v.parseInverse()
	if !ok {
		factor, token, ok = v.parseZeroOrOne()
	}
	if !ok {
		factor, token, ok = v.parseZeroOrMore()
	}
	if !ok {
		factor, token, ok = v.parseOneOrMore()
	}
	if !ok {
		factor, token, ok = v.parseExactlyN()
	}
	if !ok {
		factor, token, ok = v.parseRange()
	}
	if !ok {
		factor, token, ok = v.parseRune()
	}
	if !ok {
		factor, token, ok = v.parseString()
	}
	if !ok {
		factor, token, ok = v.parseNumber()
	}
	if !ok {
		factor, token, ok = v.parseIntrinsic()
	}
	if !ok {
		factor, token, ok = v.parseTokenName()
	}
	if !ok {
		factor, token, ok = v.parseRuleName()
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

// This method attempts to parse a token name token. It returns the token and
// whether or not a token name token was found.
func (v *parser) parseTokenName() (TokenName, *Token, bool) {
	var identifier TokenName
	var token = v.nextToken()
	if token.Type != TokenTokenName {
		v.backupOne()
		return identifier, token, false
	}
	identifier = TokenName(token.Value)
	return identifier, token, true
}

// This method attempts to parse a token name token. It returns the token and
// whether or not a token name token was found.
func (v *parser) parseRuleName() (RuleName, *Token, bool) {
	var identifier RuleName
	var token = v.nextToken()
	if token.Type != TokenRuleName {
		v.backupOne()
		return identifier, token, false
	}
	identifier = RuleName(token.Value)
	return identifier, token, true
}

// This method attempts to parse a intrinsic. It returns the intrinsic and
// whether or not the intrinsic was successfully parsed.
func (v *parser) parseIntrinsic() (Intrinsic, *Token, bool) {
	var intrinsic Intrinsic
	var token = v.nextToken()
	if token.Type != TokenIntrinsic {
		v.backupOne()
		return intrinsic, token, false
	}
	intrinsic = Intrinsic(token.Value)
	return intrinsic, token, true
}

// This method attempts to parse an inverse. It returns the inverse and
// whether or not the inverse was successfully parsed.
func (v *parser) parseInverse() (InverseLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var inverse InverseLike
	_, token, ok = v.parseDelimiter("~")
	if !ok {
		// This is not an inverse.
		return inverse, token, false
	}
	factor, token, ok = v.parseFactor()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("factor",
			"$factor")
		panic(message)
	}
	inverse = Inverse(factor)
	return inverse, token, true
}

// This method attempts to parse a note. It returns the note and whether or not
// the note was successfully parsed.
func (v *parser) parseNote() (Note, *Token, bool) {
	var note Note
	var token = v.nextToken()
	if token.Type != TokenNote {
		v.backupOne()
		return note, token, false
	}
	note = Note(token.Value)
	return note, token, true
}

// This method attempts to parse a number. It returns the number and
// whether or not a number was successfully parsed.
func (v *parser) parseNumber() (Number, *Token, bool) {
	var number Number
	var token = v.nextToken()
	if token.Type != TokenNumber {
		v.backupOne()
		return number, token, false
	}
	number = Number(token.Value)
	return number, token, true
}

// This method attempts to parse a one or more group. It returns the one or
// more group and whether or not the one or more group was successfully
// parsed.
func (v *parser) parseOneOrMore() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var definition DefinitionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("<")
	if !ok {
		// This is not a one or more group.
		return group, token, false
	}
	definition, token, ok = v.parseDefinition()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("definition",
			"$factor",
			"$definition")
		panic(message)
	}
	definition.SetMultilined(false)
	_, token, ok = v.parseDelimiter(">")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(">",
			"$factor",
			"$definition")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(definition, OneOrMore, number)
	return group, token, true
}

// This method attempts to parse a production. It returns the production and
// whether or not the production was successfully parsed.
func (v *parser) parseProduction() (ProductionLike, *Token, bool) {
	var ok bool
	var token *Token
	var symbol Symbol
	var definition DefinitionLike
	var production ProductionLike
	symbol, token, ok = v.parseSymbol()
	if !ok {
		// This is not a production.
		return production, token, false
	}
	_, token, ok = v.parseDelimiter(":")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(":",
			"$production",
			"$SYMBOL",
			"$definition")
		panic(message)
	}
	definition, token, ok = v.parseDefinition()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("definition",
			"$production",
			"$SYMBOL",
			"$definition")
		panic(message)
	}
	production = Production(symbol, definition)
	return production, token, true
}

// This method attempts to parse a range. It returns the range and whether or
// not the range was successfully parsed.
func (v *parser) parseRange() (RangeLike, *Token, bool) {
	var ok bool
	var token *Token
	var first Rune
	var last Rune
	var range_ RangeLike
	first, token, ok = v.parseRune()
	if !ok {
		// This is not a range.
		return range_, token, false
	}
	_, token, ok = v.parseDelimiter("..")
	if !ok {
		// This is not a range.
		v.backupOne() // Put back the rune.
		return range_, token, false
	}
	last, token, ok = v.parseRune()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("CHARACTER",
			"$range",
			"$CHARACTER")
		panic(message)
	}
	range_ = Range(first, last)
	return range_, token, true
}

// This method attempts to parse a rune. It returns the rune and
// whether or not a rune was successfully parsed.
func (v *parser) parseRune() (Rune, *Token, bool) {
	var rune_ Rune
	var token = v.nextToken()
	if token.Type != TokenRune {
		v.backupOne()
		return rune_, token, false
	}
	rune_ = Rune(token.Value)
	return rune_, token, true
}

// This method attempts to parse a statement. It returns the statement and
// whether or not the statement was successfully parsed.
func (v *parser) parseStatement() (StatementLike, *Token, bool) {
	var ok bool
	var token *Token
	var annotation Annotation
	var production ProductionLike
	var statement StatementLike
	annotation, _, ok = v.parseAnnotation()
	if !ok {
		production, token, ok = v.parseProduction()
		if !ok {
			// This is not a statement.
			return statement, token, false
		}
	}
	for {
		_, _, ok = v.parseEOL()
		if !ok {
			// No more blank lines.
			break
		}
	}
	statement = Statement(annotation, production)
	return statement, token, true
}

// This method attempts to parse a string. It returns the string and whether
// or not the string was successfully parsed.
func (v *parser) parseString() (String, *Token, bool) {
	var string_ String
	var token = v.nextToken()
	if token.Type != TokenString {
		v.backupOne()
		return string_, token, false
	}
	string_ = String(token.Value)
	return string_, token, true
}

// This method attempts to parse a symbol. It returns the symbol and whether
// or not the symbol was successfully parsed.
func (v *parser) parseSymbol() (Symbol, *Token, bool) {
	var symbol Symbol
	var token = v.nextToken()
	if token.Type != TokenSymbol {
		v.backupOne()
		return symbol, token, false
	}
	symbol = Symbol(token.Value)
	return symbol, token, true
}

// This method attempts to parse a zero or more group. It returns the zero or
// more group and whether or not the zero or more group was successfully
// parsed.
func (v *parser) parseZeroOrMore() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var definition DefinitionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("{")
	if !ok {
		// This is not a zero or more group.
		return group, token, false
	}
	definition, token, ok = v.parseDefinition()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("definition",
			"$factor",
			"$definition")
		panic(message)
	}
	definition.SetMultilined(false)
	_, token, ok = v.parseDelimiter("}")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("}",
			"$factor",
			"$definition")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(definition, ZeroOrMore, number)
	return group, token, true
}

// This method attempts to parse a zero or one group. It returns the zero or
// one group and whether or not the zero or one group was successfully
// parsed.
func (v *parser) parseZeroOrOne() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var definition DefinitionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("[")
	if !ok {
		// This is not a zero or one group.
		return group, token, false
	}
	definition, token, ok = v.parseDefinition()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("definition",
			"$factor",
			"$definition")
		panic(message)
	}
	definition.SetMultilined(false)
	_, token, ok = v.parseDelimiter("]")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("]",
			"$factor",
			"$definition")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(definition, ZeroOrOne, number)
	return group, token, true
}

// GRAMMAR UTILITIES

// This map captures the syntax definitions for Crater Dog Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar_ = map[string]string{
	"$grammar":     `<statement> EOF  ! Terminated by an end-of-file marker.`,
	"$statement":   `(annotation | production) <EOL>`,
	"$annotation":  `NOTE | COMMENT`,
	"$production":  `SYMBOL ":" definition`,
	"$definition":  `[EOL] alternative {[EOL] "|" alternative}`,
	"$alternative": `<factor> [NOTE]`,
	"$range":       `RUNE ".." RUNE  ! A range includes the first and last RUNEs listed.`,
	"$CHARACTER":   `"'" ~"'" "'"`,
	"$LITERAL":     `'"' <~'"'> '"'`,
	"$INTRINSIC":   `"LETTER" | "DIGIT" | "EOL" | "EOF"`,
	"$SYMBOL":      `"$" IDENTIFIER`,
	"$IDENTIFIER":  `LETTER {LETTER | DIGIT}`,
	"$NUMBER":      `<DIGIT>`,
	"$NOTE":        `"! " {~EOL}`,
	"$COMMENT":     `"!>" EOL {COMMENT | ~"<!"} EOL "<!"`,
	"$factor": `
      range  ! A range of runes.
    | "~" factor  ! The inverse of the factor.
    | "[" definition "]"  ! Zero or one instances of the definition.
    | "(" definition ")" [NUMBER]  ! Exact (default one) number of instances of the definition.
    | "<" definition ">" [NUMBER]  ! Minimum (default one) number of instances of the definition.
    | "{" definition "}" [NUMBER]  ! Maximum (default unlimited) number of instances of the definition.
    | RUNE
    | STRING
    | NUMBER
    | INTRINSIC
    | IDENTIFIER
    | DELIMITER`,
}

const header = `!>
    A formal definition of Crater Dog Syntax Notation™ (CDSN) using Crater Dog
    Syntax Notation™ itself.  This language grammar consists of rule definitions
	and token definitions.

    Each rule name begins with a lowercase letter.  The rules are applied in the
	order listed. So—for example—within a factor, a range of RUNEs takes
	precedence over an individual RUNE.  The starting rule is the "$grammar" rule.

	Each token name begins with an uppercase letter.  The INTRINSIC tokens are
	environment and language specific, and are therefore left undefined. The
	tokens are also scanned in the order listed.  So an INTRINSIC token takes
	precedence over an IDENTIFIER token.
<!

`

func FormatGrammar() string {
	var builder sts.Builder
	builder.WriteString(header)
	var unsorted = make([]string, len(grammar_))
	var index = 0
	for key := range grammar_ {
		unsorted[index] = key
		index++
	}
	var keys = col.ListFromArray(unsorted)
	keys.SortValues()
	var iterator = col.Iterator[string](keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = grammar_[key]
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", key, value))
	}
	return builder.String()
}

// PRIVATE FUNCTIONS

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar_[symbol])
	}
	return message
}
