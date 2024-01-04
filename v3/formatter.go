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
	sts "strings"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type formatterClass_ struct {
	// This class does not define any class constants.
}

// Private Class Namespace Reference

var formatterClass = &formatterClass_{
	// This class does not initialize any class constants.
}

// Public Class Namespace Access

func FormatterClass() FormatterClassLike {
	return formatterClass
}

// Public Class Constructors

func (c *formatterClass_) Default() FormatterLike {
	var formatter = &formatter_{}
	return formatter
}

// CLASS INSTANCES

// Private Class Type Definition

type formatter_ struct {
	depth  int
	result sts.Builder
}

// Public Interface

func (v *formatter_) FormatDefinition(definition DefinitionLike) string {
	v.formatDefinition(definition)
	return v.getResult()
}

func (v *formatter_) FormatDocument(document DocumentLike) string {
	v.formatDocument(document)
	return v.getResult()
}

// Private Interface

func (v *formatter_) appendNewline() {
	var separator = "\n"
	for level := 0; level < v.depth; level++ {
		separator += "    "
	}
	v.appendString(separator)
}

func (v *formatter_) appendString(s string) {
	v.result.WriteString(s)
}

func (v *formatter_) formatAlternative(alternative AlternativeLike) {
	var factors = alternative.GetFactors()
	var iterator = factors.GetIterator()
	var factor = iterator.GetNext()
	v.formatFactor(factor)
	for iterator.HasNext() {
		v.appendString(" ")
		factor = iterator.GetNext()
		v.formatFactor(factor)
	}
	var note = alternative.GetNote()
	if len(note) > 0 {
		v.appendString("  ")
		v.appendString(note)
	}
}

func (v *formatter_) formatAssertion(assertion AssertionLike) {
	var element = assertion.GetElement()
	var glyph = assertion.GetGlyph()
	var precedence = assertion.GetPrecedence()
	switch {
	case element != nil:
		v.formatElement(element)
	case glyph != nil:
		v.formatGlyph(glyph)
	case precedence != nil:
		v.formatPrecedence(precedence)
	default:
		panic("Attempted to format an empty assertion.")
	}
}

func (v *formatter_) formatCardinality(cardinality CardinalityLike) {
	var constraint = cardinality.GetConstraint()
	var first = constraint.GetFirst()
	var last = constraint.GetLast()
	switch {
	case first == "1" && last == "1":
		// This is the default case so do nothing.
	case first == "0" && last == "1":
		v.appendString("?")
	case first == "0" && len(last) == 0:
		v.appendString("*")
	case first == "1" && len(last) == 0:
		v.appendString("+")
	case len(first) > 0:
		v.formatConstraint(constraint)
	default:
		panic("Attempted to format an invalid cardinality.")
	}
}

func (v *formatter_) formatConstraint(constraint ConstraintLike) {
	var first = constraint.GetFirst()
	var last = constraint.GetLast()
	v.appendString("{")
	v.appendString(first)
	if first != last {
		v.appendString("..")
		if len(last) > 0 {
			v.appendString(last)
		}
	}
	v.appendString("}")
}

func (v *formatter_) formatDefinition(definition DefinitionLike) {
	var symbol = definition.GetSymbol()
	v.appendString(symbol)
	v.appendString(":")
	var expression = definition.GetExpression()
	if !expression.IsMultilined() {
		v.appendString(" ")
	}
	v.formatExpression(expression)
}

func (v *formatter_) formatDocument(document DocumentLike) {
	var grammar = document.GetGrammar()
	v.formatGrammar(grammar)
}

func (v *formatter_) formatElement(element ElementLike) {
	var intrinsic = element.GetIntrinsic()
	var name = element.GetName()
	var literal = element.GetLiteral()
	switch {
	case len(intrinsic) > 0:
		v.appendString(intrinsic)
	case len(name) > 0:
		v.appendString(name)
	case len(literal) > 0:
		v.appendString(literal)
	default:
		panic("Attempted to format an empty element.")
	}
}

func (v *formatter_) formatExpression(expression ExpressionLike) {
	var alternative AlternativeLike
	var alternatives = expression.GetAlternatives()
	var iterator = alternatives.GetIterator()
	if expression.IsMultilined() {
		v.depth++
		v.appendNewline()
		for iterator.HasNext() {
			alternative = iterator.GetNext()
			v.formatAlternative(alternative)
			if iterator.GetSlot() == alternatives.GetSize() {
				// The last newline must be out-dented by one.
				v.depth--
			}
			v.appendNewline()
		}
	} else {
		alternative = iterator.GetNext()
		v.formatAlternative(alternative)
		for iterator.HasNext() {
			v.appendString(" | ")
			alternative = iterator.GetNext()
			v.formatAlternative(alternative)
		}
	}
}

func (v *formatter_) formatFactor(factor FactorLike) {
	var predicate = factor.GetPredicate()
	v.formatPredicate(predicate)
	var cardinality = factor.GetCardinality()
	if cardinality != nil {
		v.formatCardinality(cardinality)
	}
}

func (v *formatter_) formatGlyph(glyph GlyphLike) {
	var first = glyph.GetFirst()
	v.appendString(first)
	var last = glyph.GetLast()
	if len(last) > 0 {
		v.appendString("..")
		v.appendString(last)
	}
}

func (v *formatter_) formatGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	var iterator = statements.GetIterator()
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		// Prepend a newline before all comments unless the first line is a comment.
		if iterator.GetSlot() > 1 && len(statement.GetComment()) > 0 {
			v.appendNewline()
		}
		v.formatStatement(statement)
		v.appendNewline()
	}
}

func (v *formatter_) formatPrecedence(precedence PrecedenceLike) {
	v.appendString("(")
	var expression = precedence.GetExpression()
	v.formatExpression(expression)
	v.appendString(")")
}

func (v *formatter_) formatPredicate(predicate PredicateLike) {
	var assertion = predicate.GetAssertion()
	var isInverted = predicate.IsInverted()
	if isInverted {
		v.appendString("~")
	}
	v.formatAssertion(assertion)
}

func (v *formatter_) formatStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.appendString(comment)
	} else {
		var definition = statement.GetDefinition()
		if definition == nil {
			panic("A statement must have either a comment or definition.")
		}
		v.formatDefinition(definition)
	}
}

func (v *formatter_) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}
