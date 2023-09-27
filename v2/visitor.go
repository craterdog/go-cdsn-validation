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
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	ref "reflect"
)

// VISITOR INTERFACE

// This function applies the specified agent to each node in the specified
// grammar.
func VisitGrammar(agent Specialized, grammar GrammarLike) {
	var v = &visitor{agent, 0}
	v.agent.BeforeGrammar(grammar)
	v.visitGrammar(grammar)
	v.agent.AfterGrammar(grammar)
}

// This interface defines the methods that are supported by specialized agents.
type Specialized interface {
	IncrementDepth()
	DecrementDepth()
	AtCharacter(character Character)
	BetweenCharacters(first Character, last Character)
	AtComment(comment Comment)
	AtIntrinsic(intrinsic Intrinsic)
	AtName(name Name)
	AtNote(note Note)
	AtNumber(number Number)
	AtString(string_ String)
	AtSymbol(symbol Symbol, isMultiline bool)
	BeforeAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool)
	AfterAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool)
	BeforeDefinition(definition DefinitionLike)
	AfterDefinition(definition DefinitionLike)
	BeforeElement(element Element)
	AfterElement(element Element)
	BeforeExactlyN(exactlyN ExactlyNLike, n Number)
	AfterExactlyN(exactlyN ExactlyNLike, n Number)
	BeforeExpression(expression ExpressionLike)
	AfterExpression(expression ExpressionLike)
	BeforeFactor(factor Factor, slot int, size int)
	AfterFactor(factor Factor, slot int, size int)
	BeforeGrammar(grammar GrammarLike)
	AfterGrammar(grammar GrammarLike)
	BeforeGrouping(grouping Grouping)
	AfterGrouping(grouping Grouping)
	BeforeInverse(inverse InverseLike)
	AfterInverse(inverse InverseLike)
	BeforeOneOrMore(oneOrMore OneOrMoreLike)
	AfterOneOrMore(oneOrMore OneOrMoreLike)
	BeforeRange(range_ RangeLike)
	AfterRange(range_ RangeLike)
	BeforeStatement(statement StatementLike, slot int, size int)
	AfterStatement(statement StatementLike, slot int, size int)
	BeforeZeroOrMore(zeroOrMore ZeroOrMoreLike)
	AfterZeroOrMore(zeroOrMore ZeroOrMoreLike)
	BeforeZeroOrOne(zeroOrOne ZeroOrOneLike)
	AfterZeroOrOne(zeroOrOne ZeroOrOneLike)
}

// VISITOR IMPLEMENTATION

// This type defines the structure and methods for a visitor.
type visitor struct {
	agent Specialized
	depth int
}

// This private method visits the specified alternative.
func (v *visitor) visitAlternative(alternative AlternativeLike) {
	var factors = alternative.GetFactors()
	var size = factors.GetSize()
	var iterator = col.Iterator(factors)
	for iterator.HasNext() {
		var slot = iterator.GetSlot()
		var factor = iterator.GetNext()
		v.agent.BeforeFactor(factor, slot, size)
		v.visitFactor(factor)
		slot++
		v.agent.AfterFactor(factor, slot, size)
	}
	var note = alternative.GetNote()
	if len(note) > 0 {
		v.agent.AtNote(note)
	}
}

// This private method visits the specified definition.
func (v *visitor) visitDefinition(definition DefinitionLike) {
	var symbol = definition.GetSymbol()
	var expression = definition.GetExpression()
	v.agent.AtSymbol(symbol, expression.IsMultilined())
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}

// This private method visits the specified element.
func (v *visitor) visitElement(element Element) {
	switch actual := element.(type) {
	case Intrinsic:
		v.agent.AtIntrinsic(actual)
	case String:
		v.agent.AtString(actual)
	case Number:
		v.agent.AtNumber(actual)
	case Name:
		v.agent.AtName(actual)
	default:
		panic(fmt.Sprintf("Attempted to visit:\n    element: %v\n    type: %t\n", actual, element))
	}
}

// This private method visits the specified exactly N grouping.
func (v *visitor) visitExactlyN(group ExactlyNLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}

// This private method visits the specified expression.
func (v *visitor) visitExpression(expression ExpressionLike) {
	var isMultilined = expression.IsMultilined()
	var alternatives = expression.GetAlternatives()
	var size = alternatives.GetSize()
	var iterator = col.Iterator(alternatives)
	v.agent.IncrementDepth()
	for iterator.HasNext() {
		var slot = iterator.GetSlot()
		var alternative = iterator.GetNext()
		v.agent.BeforeAlternative(alternative, slot, size, isMultilined)
		v.visitAlternative(alternative)
		slot++
		v.agent.AfterAlternative(alternative, slot, size, isMultilined)
	}
	v.agent.DecrementDepth()
}

// This private method visits the specified factor.
func (v *visitor) visitFactor(factor Factor) {
	if ref.ValueOf(factor).Kind() == ref.String {
		v.agent.BeforeElement(factor)
		v.visitElement(factor)
		v.agent.AfterElement(factor)
		return
	}
	switch actual := factor.(type) {
	case *range_:
		v.agent.BeforeRange(actual)
		v.visitRange(actual)
		v.agent.AfterRange(actual)
	case *inverse:
		v.agent.BeforeInverse(actual)
		v.visitInverse(actual)
		v.agent.AfterInverse(actual)
	default:
		v.agent.BeforeGrouping(actual)
		v.visitGrouping(actual)
		v.agent.AfterGrouping(actual)
	}
}

// This private method visits the specified grammar.
func (v *visitor) visitGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	var size = statements.GetSize()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var slot = iterator.GetSlot()
		var statement = iterator.GetNext()
		v.agent.BeforeStatement(statement, slot, size)
		v.visitStatement(statement)
		slot++
		v.agent.AfterStatement(statement, slot, size)
	}
}

// This private method visits the specified grouping.
func (v *visitor) visitGrouping(grouping Grouping) {
	switch actual := grouping.(type) {
	case *exactlyN:
		var n = actual.GetN()
		v.agent.BeforeExactlyN(actual, n)
		v.visitExactlyN(actual)
		v.agent.AfterExactlyN(actual, n)
	case *zeroOrOne:
		v.agent.BeforeZeroOrOne(actual)
		v.visitZeroOrOne(actual)
		v.agent.AfterZeroOrOne(actual)
	case *zeroOrMore:
		v.agent.BeforeZeroOrMore(actual)
		v.visitZeroOrMore(actual)
		v.agent.AfterZeroOrMore(actual)
	case *oneOrMore:
		v.agent.BeforeOneOrMore(actual)
		v.visitOneOrMore(actual)
		v.agent.AfterOneOrMore(actual)
	default:
		panic(fmt.Sprintf("Attempted to visit:\n    grouping: %v\n    type: %t\n", actual, grouping))
	}
}

// This private method visits the specified inverse.
func (v *visitor) visitInverse(inverse InverseLike) {
	var factor = inverse.GetFactor()
	v.agent.BeforeFactor(factor, 0, 0)
	v.visitFactor(factor)
	v.agent.AfterFactor(factor, 0, 0)
}

// This private method visits the specified one or more grouping.
func (v *visitor) visitOneOrMore(group OneOrMoreLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}

// This private method visits the specified range.
func (v *visitor) visitRange(range_ RangeLike) {
	var first = range_.GetFirstCharacter()
	v.agent.AtCharacter(first)
	var last = range_.GetLastCharacter()
	if len(last) > 0 {
		v.agent.BetweenCharacters(first, last)
		v.agent.AtCharacter(last)
	}
}

// This private method visits the specified statement.
func (v *visitor) visitStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.agent.AtComment(comment)
	} else {
		var definition = statement.GetDefinition()
		v.agent.BeforeDefinition(definition)
		v.visitDefinition(definition)
		v.agent.AfterDefinition(definition)
	}
}

// This private method visits the specified zero or more grouping.
func (v *visitor) visitZeroOrMore(group ZeroOrMoreLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}

// This private method visits the specified zero or one grouping.
func (v *visitor) visitZeroOrOne(group ZeroOrOneLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}
