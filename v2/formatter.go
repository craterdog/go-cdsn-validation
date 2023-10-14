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
	col "github.com/craterdog/go-collection-framework/v2"
	sts "strings"
)

// FORMATTER INTERFACE

// This function returns the bytes containing the canonical format for the
// specified document including the POSIX standard EOF marker.
func FormatDocument(document DocumentLike) []byte {
	var v = &formatter{}
	v.formatDocument(document)
	return []byte(v.getResult())
}

// FORMATTER IMPLEMENTATION

// This type defines the structure and methods for a canonical formatter agent.
type formatter struct {
	depth  int
	result sts.Builder
}

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

// This private method formats the specified alternative.
func (v *formatter) formatAlternative(alternative AlternativeLike) {
	var predicates = alternative.GetPredicates()
	var iterator = col.Iterator(predicates)
	var predicate = iterator.GetNext()
	v.formatPredicate(predicate)
	for iterator.HasNext() {
		predicate = iterator.GetNext()
		v.appendString(" ")
		v.formatPredicate(predicate)
	}
	var note = alternative.GetNOTE()
	if len(note) > 0 {
		v.appendString("  ")
		v.appendString(string(note))
	}
}

// This private method formats the specified definition.
func (v *formatter) formatDefinition(definition DefinitionLike) {
	var symbol = definition.GetSYMBOL()
	v.appendString(string(symbol))
	v.appendString(":")
	var expression = definition.GetExpression()
	if !expression.IsAnnotated() {
		v.appendString(" ")
	}
	v.formatExpression(expression)
}

// This private method formats the specified document.
func (v *formatter) formatDocument(document DocumentLike) {
	var statements = document.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.formatStatement(statement)
	}
}

// This private method formats the specified element.
func (v *formatter) formatElement(element ElementLike) {
	var intrinsic = element.GetINTRINSIC()
	var name = element.GetNAME()
	var literal = element.GetLITERAL()
	switch {
	case len(intrinsic) > 0:
		v.appendString(string(intrinsic))
	case len(name) > 0:
		v.appendString(string(name))
	case len(literal) > 0:
		v.appendString(string(literal))
	}
}

// This private method formats the specified expression.
func (v *formatter) formatExpression(expression ExpressionLike) {
	v.depth++
	if expression.IsAnnotated() {
		v.appendNewline()
		v.appendString("  ") // Indent additional two spaces to align with subsequent alternatives.
	}
	var alternatives = expression.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	var alternative = iterator.GetNext()
	v.formatAlternative(alternative)
	for iterator.HasNext() {
		alternative = iterator.GetNext()
		if expression.IsAnnotated() {
			v.appendNewline()
		} else {
			v.appendString(" ")
		}
		v.appendString("| ")
		v.formatAlternative(alternative)
	}
	v.depth--
}

// This private method formats the specified factor.
func (v *formatter) formatFactor(factor FactorLike) {
	var precedence = factor.GetPrecedence()
	var element = factor.GetElement()
	switch {
	case precedence != nil:
		v.formatPrecedence(precedence)
	case element != nil:
		v.formatElement(element)
	}
}

// This private method formats the specified precedence.
func (v *formatter) formatPrecedence(definition PrecedenceLike) {
	v.appendString("(")
	var expression = definition.GetExpression()
	v.formatExpression(expression)
	v.appendString(")")
}

// This private method formats the specified predicate.
func (v *formatter) formatPredicate(predicate PredicateLike) {
	var glyph = predicate.GetGlyph()
	var repetition = predicate.GetRepetition()
	var factor = predicate.GetFactor()
	switch {
	case glyph != nil:
		v.formatGlyph(glyph)
	case repetition != nil:
		v.formatRepetition(repetition)
	case factor != nil:
		v.formatFactor(factor)
	}
}

// This private method formats the specified glyph.
func (v *formatter) formatGlyph(glyph GlyphLike) {
	var first = glyph.GetFirstCHARACTER()
	v.appendString(string(first))
	var last = glyph.GetLastCHARACTER()
	if len(last) > 0 {
		v.appendString("..")
		v.appendString(string(last))
	}
}

// This private method formats the specified repetition.
func (v *formatter) formatRepetition(repetition RepetitionLike) {
	var constraint = repetition.GetCONSTRAINT()
	v.appendString(string(constraint))
	var factor = repetition.GetFactor()
	v.formatFactor(factor)
}

// This private method formats the specified statement.
func (v *formatter) formatStatement(statement Statement) {
	switch actual := statement.(type) {
	case *definition:
		v.formatDefinition(actual)
	case COMMENT:
		v.appendNewline()
		v.appendString(string(actual))
	}
	v.appendNewline()
	v.appendNewline()
}

// This private method returns the canonically formatted string result.
func (v *formatter) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}
