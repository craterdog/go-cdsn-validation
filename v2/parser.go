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
func ParseDocument(source []byte) DocumentLike {
	var ok bool
	var token *Token
	var document DocumentLike
	var tokens = make(chan Token, 256)
	ScanTokens(source, tokens) // Starts scanning in a separate go routine.
	var p = &parser{
		symbols: col.Catalog[SYMBOL, DefinitionLike](),
		source:  source,
		next:    col.StackWithCapacity[*Token](4),
		tokens:  tokens,
	}
	document, token, ok = p.parseDocument()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("statement",
			"$document",
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
	return document
}

// PARSER IMPLEMENTATION

// This map captures the syntax expressions for Crater Dog Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar = map[string]string{
	"$document":    `statement+ EOF  ! Terminated with an end-of-file marker.`,
	"$statement":   `definition | COMMENT`,
	"$definition":  `SYMBOL ":" expression  ! This works for both tokens and rules.`,
	"$expression":  `alternative ("|" alternative)*`,
	"$alternative": `predicate+ NOTE?`,
	"$predicate":   `factor cardinality?  ! The default cardinality is one.`,
	"$factor":      `element | glyph | inversion | precedence`,
	"$element":     `INTRINSIC | NAME | LITERAL`,
	"$glyph":       `CHARACTER (".." CHARACTER)?  ! The range of CHARACTERs in a glyph is inclusive.`,
	"$inversion":   `"~" factor`,
	"$precedence":  `"(" expression ")"`,
	"$cardinality": `
      LIMIT
    | "{" NUMBER (".." NUMBER?)? "}"  ! The range of NUMBERs in a cardinality is inclusive.`,
}

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar[symbol])
	}
	return message
}

// This type defines the structure and methods for the parser agent.
type parser struct {
	symbols col.CatalogLike[SYMBOL, DefinitionLike]
	source  []byte
	next    col.StackLike[*Token] // The stack of the retrieved tokens that have been put back.
	tokens  chan Token            // The queue of unread tokens coming from the scanner.
	isToken bool                  // Whether or not the current definition is a token definition.
}

// This method puts back the current token onto the token stream so that it can
// be retrieved by another parsing method.
func (v *parser) backupOne(token *Token) {
	v.next.AddValue(token)
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
	return next
}

// This method attempts to parse a new end-of-file token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseEOF() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOF {
		v.backupOne(token)
		return token, token, false
	}
	return token, token, true
}

// This method attempts to parse the specified delimiter token. It returns the
// token and whether or not the token was successfully parsed.
func (v *parser) parseDELIMITER(delimiter string) (string, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenDELIMITER || token.Value != delimiter {
		v.backupOne(token)
		return delimiter, token, false
	}
	return delimiter, token, true
}

// This method attempts to parse a new intrinsic token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseINTRINSIC() (INTRINSIC, *Token, bool) {
	var intrinsic INTRINSIC
	var token = v.nextToken()
	if token.Type != TokenINTRINSIC {
		v.backupOne(token)
		return intrinsic, token, false
	}
	intrinsic = INTRINSIC(token.Value)
	return intrinsic, token, true
}

// This method attempts to parse a new note token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseNOTE() (NOTE, *Token, bool) {
	var note NOTE
	var token = v.nextToken()
	if token.Type != TokenNOTE {
		v.backupOne(token)
		return note, token, false
	}
	note = NOTE(token.Value)
	return note, token, true
}

// This method attempts to parse a new comment token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseCOMMENT() (COMMENT, *Token, bool) {
	var comment COMMENT
	var token = v.nextToken()
	if token.Type != TokenCOMMENT {
		v.backupOne(token)
		return comment, token, false
	}
	comment = COMMENT(token.Value)
	return comment, token, true
}

// This method attempts to parse a new number token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseNUMBER() (NUMBER, *Token, bool) {
	var number NUMBER
	var token = v.nextToken()
	if token.Type != TokenNUMBER {
		v.backupOne(token)
		return number, token, false
	}
	number = NUMBER(token.Value)
	return number, token, true
}

// This method attempts to parse a new character token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseCHARACTER() (CHARACTER, *Token, bool) {
	var character CHARACTER
	var token = v.nextToken()
	if token.Type != TokenCHARACTER {
		v.backupOne(token)
		return character, token, false
	}
	character = CHARACTER(token.Value)
	return character, token, true
}

// This method attempts to parse a new literal token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseLITERAL() (LITERAL, *Token, bool) {
	var literal LITERAL
	var token = v.nextToken()
	if token.Type != TokenLITERAL {
		v.backupOne(token)
		return literal, token, false
	}
	literal = LITERAL(token.Value)
	return literal, token, true
}

// This method attempts to parse a new name token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseNAME() (NAME, *Token, bool) {
	var name NAME
	var token = v.nextToken()
	if token.Type != TokenNAME {
		v.backupOne(token)
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

// This method attempts to parse a new symbol token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseSYMBOL() (SYMBOL, *Token, bool) {
	var symbol SYMBOL
	var token = v.nextToken()
	if token.Type != TokenSYMBOL {
		v.backupOne(token)
		return symbol, token, false
	}
	symbol = SYMBOL(token.Value)
	return symbol, token, true
}

// This method attempts to parse a new constraint token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parseCONSTRAINT() (CONSTRAINT, *Token, bool) {
	var constraint CONSTRAINT
	var token = v.nextToken()
	if token.Type != TokenCONSTRAINT {
		v.backupOne(token)
		return constraint, token, false
	}
	constraint = CONSTRAINT(token.Value)
	return constraint, token, true
}

// This method attempts to parse a new document. It returns the document
// and whether or not the document was successfully parsed.
func (v *parser) parseDocument() (DocumentLike, *Token, bool) {
	var ok bool
	var token *Token
	var statement StatementLike
	var statements = col.List[StatementLike]()
	var document DocumentLike
	statement, token, ok = v.parseStatement()
	if !ok {
		// This is not a document.
		return document, token, false
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
			"$document",
			"$statement")
		panic(message)
	}
	document = Document(statements)
	return document, token, true
}

// This method attempts to parse a new statement. It returns the statement
// and whether or not the statement was successfully parsed.
func (v *parser) parseStatement() (StatementLike, *Token, bool) {
	var ok bool
	var token *Token
	var definition DefinitionLike
	var comment COMMENT
	var statement StatementLike
	definition, token, ok = v.parseDefinition()
	if !ok {
		comment, token, ok = v.parseCOMMENT()
	}
	if !ok {
		// This is not a statement.
		return statement, token, false
	}
	statement = Statement(definition, comment)
	return statement, token, true
}

// This method attempts to parse a new definition. It returns the definition
// and whether or not the definition was successfully parsed.
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
	_, token, ok = v.parseDELIMITER(":")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(":",
			"$definition",
			"$expression")
		panic(message)
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$definition",
			"$expression")
		panic(message)
	}
	definition = Definition(symbol, expression)
	v.symbols.SetValue(symbol, definition)
	return definition, token, true
}

// This method attempts to parse a new expression. It returns the expression
// and whether or not the expression was successfully parsed.
func (v *parser) parseExpression() (ExpressionLike, *Token, bool) {
	var ok bool
	var token *Token
	var alternative AlternativeLike
	var alternatives = col.List[AlternativeLike]()
	var expression ExpressionLike
	alternative, token, ok = v.parseAlternative()
	if !ok {
		// An expression must have at least one alternative.
		return expression, token, false
	}
	for {
		alternatives.AddValue(alternative)
		_, _, ok = v.parseDELIMITER("|")
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

// This method attempts to parse a new alternative. It returns the alternative
// and whether or not the alternative was successfully parsed.
func (v *parser) parseAlternative() (AlternativeLike, *Token, bool) {
	var ok bool
	var token *Token
	var predicate PredicateLike
	var predicates = col.List[PredicateLike]()
	var note NOTE
	var alternative AlternativeLike
	predicate, token, ok = v.parsePredicate()
	if !ok {
		// An alternative must have at least one predicate.
		return alternative, token, false
	}
	for {
		predicates.AddValue(predicate)
		predicate, token, ok = v.parsePredicate()
		if !ok {
			// No more predicates.
			break
		}
	}
	note, _, _ = v.parseNOTE() // The note is optional.
	alternative = Alternative(predicates, note)
	return alternative, token, true
}

// This method attempts to parse a new predicate. It returns the predicate
// and whether or not the predicate was successfully parsed.
func (v *parser) parsePredicate() (PredicateLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor FactorLike
	var cardinality CardinalityLike
	var predicate PredicateLike
	factor, token, ok = v.parseFactor()
	if !ok {
		// This is not a predicate.
		return predicate, token, false
	}
	cardinality, token, _ = v.parseCardinality()
	predicate = Predicate(factor, cardinality)
	return predicate, token, true
}

// This method attempts to parse a new factor. It returns the factor
// and whether or not the factor was successfully parsed.
func (v *parser) parseFactor() (FactorLike, *Token, bool) {
	var ok bool
	var token *Token
	var element ElementLike
	var glyph GlyphLike
	var inversion InversionLike
	var precedence PrecedenceLike
	var factor FactorLike
	element, token, ok = v.parseElement()
	if !ok {
		glyph, token, ok = v.parseGlyph()
	}
	if !ok {
		inversion, token, ok = v.parseInversion()
	}
	if !ok {
		precedence, token, ok = v.parsePrecedence()
	}
	if !ok {
		// This is not a factor.
		return factor, token, false
	}
	factor = Factor(element, glyph, inversion, precedence)
	return factor, token, true
}

// This method attempts to parse an element. It returns the element
// and whether or not the element was successfully parsed.
func (v *parser) parseElement() (ElementLike, *Token, bool) {
	var ok bool
	var token *Token
	var intrinsic INTRINSIC
	var name NAME
	var literal LITERAL
	var element ElementLike
	intrinsic, token, ok = v.parseINTRINSIC()
	if !ok {
		name, token, ok = v.parseNAME()
	}
	if !ok {
		literal, token, ok = v.parseLITERAL()
	}
	if !ok {
		// This is not an element.
		return element, token, false
	}
	element = Element(intrinsic, name, literal)
	return element, token, true
}

// This method attempts to parse a new glyph. It returns the glyph
// and whether or not the glyph was successfully parsed.
func (v *parser) parseGlyph() (GlyphLike, *Token, bool) {
	var ok bool
	var token *Token
	var first CHARACTER
	var last CHARACTER
	var glyph GlyphLike
	first, token, ok = v.parseCHARACTER()
	if !ok {
		// This is not a glyph.
		return glyph, token, false
	}
	_, _, ok = v.parseDELIMITER("..")
	if ok {
		last, token, ok = v.parseCHARACTER()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("CHARACTER",
				"$glyph")
			panic(message)
		}
	}
	glyph = Glyph(first, last)
	return glyph, token, true
}

// This method attempts to parse a new inversion. It returns the inversion
// and whether or not the inversion was successfully parsed.
func (v *parser) parseInversion() (InversionLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor FactorLike
	var inversion InversionLike
	_, token, ok = v.parseDELIMITER("~")
	if !ok {
		// This is not a inversion.
		return inversion, token, false
	}
	factor, token, ok = v.parseFactor()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("factor",
			"$inversion",
			"$factor")
		panic(message)
	}
	inversion = Inversion(factor)
	return inversion, token, true
}

// This method attempts to parse a new precedence. It returns the precedence
// and whether or not the precedence was successfully parsed.
func (v *parser) parsePrecedence() (PrecedenceLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var precedence PrecedenceLike
	_, token, ok = v.parseDELIMITER("(")
	if !ok {
		// This is not a precedence.
		return precedence, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$precedence",
			"$expression")
		panic(message)
	}
	expression.SetAnnotated(false)
	_, token, ok = v.parseDELIMITER(")")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(")",
			"$precedence",
			"$expression")
		panic(message)
	}
	precedence = Precedence(expression)
	return precedence, token, true
}

// This method attempts to parse a new cardinality. It returns the cardinality
// and whether or not the cardinality was successfully parsed.
func (v *parser) parseCardinality() (CardinalityLike, *Token, bool) {
	var ok bool
	var token *Token
	var constraint CONSTRAINT
	var first NUMBER
	var last NUMBER
	var cardinality CardinalityLike
	constraint, token, ok = v.parseCONSTRAINT()
	if !ok {
		_, token, ok = v.parseDELIMITER("{")
		if !ok {
			// This is not a cardinality.
			return cardinality, token, false
		}
		first, token, ok = v.parseNUMBER()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("NUMBER",
				"$cardinality")
			panic(message)
		}
		_, _, ok = v.parseDELIMITER("..")
		if ok {
			last, token, ok = v.parseNUMBER()
			if !ok {
				var message = v.formatError(token)
				message += generateGrammar("NUMBER",
					"$cardinality")
				panic(message)
			}
		}
		_, token, ok = v.parseDELIMITER("}")
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("}",
				"$cardinality")
			panic(message)
		}
	}
	cardinality = Cardinality(constraint, first, last)
	return cardinality, token, true
}
