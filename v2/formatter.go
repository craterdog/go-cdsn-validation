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
	result sts.Builder
}

// PRIVATE METHODS

// This public method is called for each character token.
func (v *formatter) AtCharacter(character Character, depth int) {
	v.appendString(string(character))
}

// This public method is called between the two two characters in a range.
func (v *formatter) BetweenCharacters(first Character, last Character, depth int) {
	v.appendString("..")
}

// This public method is called for each comment token.
func (v *formatter) AtComment(comment Comment, depth int) {
	v.appendNewline(depth)
	v.appendString(string(comment))
}

// This public method is called for each intrinsic token.
func (v *formatter) AtIntrinsic(intrinsic Intrinsic, depth int) {
	v.appendString(string(intrinsic))
}

// This public method is called for each name token.
func (v *formatter) AtName(name Name, depth int) {
	v.appendString(string(name))
}

// This public method is called for each note token.
func (v *formatter) AtNote(note Note, depth int) {
	v.appendString("  ")
	v.appendString(string(note))
}

// This public method is called for each number token.
func (v *formatter) AtNumber(number Number, depth int) {
	v.appendString(string(number))
}

// This public method is called for each string token.
func (v *formatter) AtString(string_ String, depth int) {
	v.appendString(string(string_))
}

// This public method is called for each symbol token.
func (v *formatter) AtSymbol(symbol Symbol, isMultiline bool, depth int) {
	v.appendString(string(symbol))
	v.appendString(":")
	if !isMultiline {
		v.appendString(" ")
	}
}

// This public method is called before each alternative in an expression.
func (v *formatter) BeforeAlternative(alternative AlternativeLike, slot int,
	size int, isAnnotated bool, depth int) {
	if isAnnotated {
		v.appendNewline(depth)
		if slot == 0 {
			v.appendString("  ")
		}
	}
	if slot > 0 {
		v.appendString("| ")
	}
}

// This public method is called after each alternative in an expression.
func (v *formatter) AfterAlternative(alternative AlternativeLike, slot int,
	size int, isAnnotated bool, depth int) {
	if !isAnnotated && slot < size {
		v.appendString(" ")
	}
}

// This public method is called before each definition.
func (v *formatter) BeforeDefinition(definition DefinitionLike, depth int) {
}

// This public method is called after each definition.
func (v *formatter) AfterDefinition(definition DefinitionLike, depth int) {
}

// This public method is called before each element.
func (v *formatter) BeforeElement(element Element, depth int) {
}

// This public method is called after each element.
func (v *formatter) AfterElement(element Element, depth int) {
}

// This public method is called before each exactly N grouping.
func (v *formatter) BeforeExactlyN(exactlyN ExactlyNLike, n Number, depth int) {
	v.appendString("(")
}

// This public method is called after each exactly N grouping.
func (v *formatter) AfterExactlyN(exactlyN ExactlyNLike, n Number, depth int) {
	v.appendString(")")
	if len(n) > 0 {
		v.AtNumber(n, depth)
	}
}

func (v *formatter) BeforeExpression(expression ExpressionLike, depth int) {
}

func (v *formatter) AfterExpression(expression ExpressionLike, depth int) {
}

func (v *formatter) BeforeFactor(
	factor Factor, slot int,
	size int, depth int) {
	if slot > 0 {
		v.appendString(" ")
	}
}

func (v *formatter) AfterFactor(
	factor Factor, slot int,
	size int, depth int) {
}

func (v *formatter) BeforeGrammar(grammar GrammarLike, depth int) {
}

func (v *formatter) AfterGrammar(grammar GrammarLike, depth int) {
}

func (v *formatter) BeforeGrouping(grouping Grouping, depth int) {
}

func (v *formatter) AfterGrouping(grouping Grouping, depth int) {
}

func (v *formatter) BeforeInverse(inverse InverseLike, depth int) {
	v.appendString("~")
}

func (v *formatter) AfterInverse(inverse InverseLike, depth int) {
}

// This public method is called before each one or more grouping.
func (v *formatter) BeforeOneOrMore(oneOrMore OneOrMoreLike, depth int) {
	v.appendString("<")
}

// This public method is called after each one or more grouping.
func (v *formatter) AfterOneOrMore(oneOrMore OneOrMoreLike, depth int) {
	v.appendString(">")
}

func (v *formatter) BeforeRange(range_ RangeLike, depth int) {
}

func (v *formatter) AfterRange(range_ RangeLike, depth int) {
}

func (v *formatter) BeforeStatement(
	statement StatementLike, slot int,
	size int, depth int) {
}

func (v *formatter) AfterStatement(
	statement StatementLike, slot int,
	size int, depth int) {
	v.appendNewline(depth)
	v.appendNewline(depth)
}

// This public method is called before each zero or more grouping.
func (v *formatter) BeforeZeroOrMore(zeroOrMore ZeroOrMoreLike, depth int) {
	v.appendString("{")
}

// This public method is called after each zero or more grouping.
func (v *formatter) AfterZeroOrMore(zeroOrMore ZeroOrMoreLike, depth int) {
	v.appendString("}")
}

// This public method is called before each zero or one grouping.
func (v *formatter) BeforeZeroOrOne(zeroOrOne ZeroOrOneLike, depth int) {
	v.appendString("[")
}

// This public method is called after each zero or one grouping.
func (v *formatter) AfterZeroOrOne(zeroOrOne ZeroOrOneLike, depth int) {
	v.appendString("]")
}

// PRIVATE METHODS

// This private method appends a properly indented newline to the result.
func (v *formatter) appendNewline(depth int) {
	var separator = "\n"
	for level := 0; level < depth; level++ {
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
