/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn

import (
	sts "strings"
)

// FORMATTER INTERFACE

// This function returns the bytes containing the canonical format for the
// specified grammar including the POSIX standard EOF marker.
func FormatGrammar(grammar GrammarLike) []byte {
	var agent = &formatter{}
	VisitGrammar(agent, grammar)
	return []byte(agent.getResult())
}

// FORMATTER IMPLEMENTATION

// This type defines the structure and methods for a canonical formatter agent.
type formatter struct {
	depth  int
	result sts.Builder
}

// PRIVATE METHODS

// This public method increments the depth of the traversal by one.
func (v *formatter) IncrementDepth() {
	v.depth++
}

// This public method decrements the depth of the traversal by one.
func (v *formatter) DecrementDepth() {
	v.depth--
}

// This public method is called for each character token.
func (v *formatter) AtCharacter(character Character) {
	v.appendString(string(character))
}

// This public method is called for each comment token.
func (v *formatter) AtComment(comment Comment) {
	v.appendNewline()
	v.appendString(string(comment))
}

// This public method is called for each intrinsic token.
func (v *formatter) AtIntrinsic(intrinsic Intrinsic) {
	v.appendString(string(intrinsic))
}

// This public method is called for each name token.
func (v *formatter) AtName(name Name) {
	v.appendString(string(name))
}

// This public method is called for each note token.
func (v *formatter) AtNote(note Note) {
	v.appendString("  ")
	v.appendString(string(note))
}

// This public method is called for each number token.
func (v *formatter) AtNumber(number Number) {
	v.appendString(string(number))
}

// This public method is called for each string token.
func (v *formatter) AtString(string_ String) {
	v.appendString(string(string_))
}

// This public method is called for each symbol token.
func (v *formatter) AtSymbol(symbol Symbol, isMultilined bool) {
	v.appendString(string(symbol))
	v.appendString(":")
	if !isMultilined {
		v.appendString(" ")
	}
}

// This public method is called before each alternative in an expression.
func (v *formatter) BeforeAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool) {
	if isMultilined {
		v.appendNewline()
		if slot == 0 {
			v.appendString("  ") // Indent additional two spaces to align with subsequent alternatives.
		}
	}
	if slot > 0 {
		v.appendString("| ")
	}
}

// This public method is called after each alternative in an expression.
func (v *formatter) AfterAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool) {
	if !isMultilined && slot < size {
		v.appendString(" ")
	}
}

// This public method is called before each definition.
func (v *formatter) BeforeDefinition(definition DefinitionLike) {
}

// This public method is called after each definition.
func (v *formatter) AfterDefinition(definition DefinitionLike) {
}

// This public method is called before each element.
func (v *formatter) BeforeElement(element Element) {
}

// This public method is called after each element.
func (v *formatter) AfterElement(element Element) {
}

// This public method is called before each exactly N grouping.
func (v *formatter) BeforeExactlyN(exactlyN ExactlyNLike, n Number) {
	v.appendString("(")
}

// This public method is called after each exactly N grouping.
func (v *formatter) AfterExactlyN(exactlyN ExactlyNLike, n Number) {
	v.appendString(")")
	if len(n) > 0 {
		v.AtNumber(n)
	}
}

// This public method is called before each expression.
func (v *formatter) BeforeExpression(expression ExpressionLike) {
}

// This public method is called after each expression.
func (v *formatter) AfterExpression(expression ExpressionLike) {
}

// This public method is called before each factor in an alternative.
func (v *formatter) BeforeFactor(factor Factor, slot int, size int) {
	if slot > 0 {
		v.appendString(" ")
	}
}

// This public method is called after each factor in an alternative.
func (v *formatter) AfterFactor(factor Factor, slot int, size int) {
}

// This public method is called before the grammar.
func (v *formatter) BeforeGrammar(grammar GrammarLike) {
}

// This public method is called after the grammar.
func (v *formatter) AfterGrammar(grammar GrammarLike) {
}

// This public method is called before each grouping.
func (v *formatter) BeforeGrouping(grouping Grouping) {
}

// This public method is called after each grouping.
func (v *formatter) AfterGrouping(grouping Grouping) {
}

// This public method is called before each inverse factor.
func (v *formatter) BeforeInverse(inverse InverseLike) {
	v.appendString("~")
}

// This public method is called after each inverse factor.
func (v *formatter) AfterInverse(inverse InverseLike) {
}

// This public method is called before each one or more grouping.
func (v *formatter) BeforeOneOrMore(oneOrMore OneOrMoreLike) {
	v.appendString("<")
}

// This public method is called after each one or more grouping.
func (v *formatter) AfterOneOrMore(oneOrMore OneOrMoreLike) {
	v.appendString(">")
}

// This public method is called before each character range.
func (v *formatter) BeforeRange(range_ RangeLike) {
}

// This public method is called between the two two characters in a range.
func (v *formatter) BetweenCharacters(first Character, last Character) {
	v.appendString("..")
}

// This public method is called after each character range.
func (v *formatter) AfterRange(range_ RangeLike) {
}

// This public method is called before each statement in a grammar.
func (v *formatter) BeforeStatement(statement StatementLike, slot int, size int) {
}

// This public method is called after each statement in a grammar.
func (v *formatter) AfterStatement(statement StatementLike, slot int, size int) {
	v.appendNewline()
	v.appendNewline()
}

// This public method is called before each zero or more grouping.
func (v *formatter) BeforeZeroOrMore(zeroOrMore ZeroOrMoreLike) {
	v.appendString("{")
}

// This public method is called after each zero or more grouping.
func (v *formatter) AfterZeroOrMore(zeroOrMore ZeroOrMoreLike) {
	v.appendString("}")
}

// This public method is called before each zero or one grouping.
func (v *formatter) BeforeZeroOrOne(zeroOrOne ZeroOrOneLike) {
	v.appendString("[")
}

// This public method is called after each zero or one grouping.
func (v *formatter) AfterZeroOrOne(zeroOrOne ZeroOrOneLike) {
	v.appendString("]")
}

// PRIVATE METHODS

// This private method appends a properly indented newline to the result.
func (v *formatter) appendNewline() {
	var separator = "\n"
	for level := 0; level < v.depth; level++ {
		separator += "    "
	}
	v.result.WriteString(separator)
}

// This private method appends the specified string to the result.
func (v *formatter) appendString(s string) {
	v.result.WriteString(s)
}

// This private method returns the canonically formatted string result.
func (v *formatter) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}
