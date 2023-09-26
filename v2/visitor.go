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

// AGENT INTERFACE

type AgentLike interface {
	AtCharacter(character Character, depth int)
	BetweenCharacters(first Character, last Character, depth int)
	AtComment(comment Comment, depth int)
	AtIntrinsic(intrinsic Intrinsic, depth int)
	AtName(name Name, depth int)
	AtNote(note Note, depth int)
	AtNumber(number Number, depth int)
	AtString(string_ String, depth int)
	AtSymbol(symbol Symbol, isMultiline bool, depth int)
	BeforeAlternative(alternative AlternativeLike, slot int, size int, isAnnotated bool, depth int)
	AfterAlternative(alternative AlternativeLike, slot int, size int, isAnnotated bool, depth int)
	BeforeDefinition(definition DefinitionLike, depth int)
	AfterDefinition(definition DefinitionLike, depth int)
	BeforeElement(element Element, depth int)
	AfterElement(element Element, depth int)
	BeforeExactlyN(exactlyN ExactlyNLike, n Number, depth int)
	AfterExactlyN(exactlyN ExactlyNLike, n Number, depth int)
	BeforeExpression(expression ExpressionLike, depth int)
	AfterExpression(expression ExpressionLike, depth int)
	BeforeFactor(factor Factor, slot int, size int, depth int)
	AfterFactor(factor Factor, slot int, size int, depth int)
	BeforeGrammar(grammar GrammarLike, depth int)
	AfterGrammar(grammar GrammarLike, depth int)
	BeforeGrouping(grouping Grouping, depth int)
	AfterGrouping(grouping Grouping, depth int)
	BeforeInverse(inverse InverseLike, depth int)
	AfterInverse(inverse InverseLike, depth int)
	BeforeOneOrMore(oneOrMore OneOrMoreLike, depth int)
	AfterOneOrMore(oneOrMore OneOrMoreLike, depth int)
	BeforeRange(range_ RangeLike, depth int)
	AfterRange(range_ RangeLike, depth int)
	BeforeStatement(statement StatementLike, slot int, size int, depth int)
	AfterStatement(statement StatementLike, slot int, size int, depth int)
	BeforeZeroOrMore(zeroOrMore ZeroOrMoreLike, depth int)
	AfterZeroOrMore(zeroOrMore ZeroOrMoreLike, depth int)
	BeforeZeroOrOne(zeroOrOne ZeroOrOneLike, depth int)
	AfterZeroOrOne(zeroOrOne ZeroOrOneLike, depth int)
}

// VISITOR INTERFACE

// This function applies the specified agent to each node in the specified
// grammar.
func VisitGrammar(agent AgentLike, grammar GrammarLike) {
	var v = &visitor{agent, 0}
	v.agent.BeforeGrammar(grammar, v.depth)
	v.visitGrammar(grammar)
	v.agent.AfterGrammar(grammar, v.depth)
}

// VISITOR IMPLEMENTATION

// This type defines the structure and methods for a visitor.
type visitor struct {
	agent AgentLike
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
		v.agent.BeforeFactor(factor, slot, size, v.depth)
		v.visitFactor(factor)
		slot++
		v.agent.AfterFactor(factor, slot, size, v.depth)
	}
	var note = alternative.GetNote()
	if len(note) > 0 {
		v.agent.AtNote(note, v.depth)
	}
}

// This private method visits the specified definition.
func (v *visitor) visitDefinition(definition DefinitionLike) {
	var symbol = definition.GetSymbol()
	var expression = definition.GetExpression()
	v.agent.AtSymbol(symbol, expression.IsAnnotated(), v.depth)
	v.agent.BeforeExpression(expression, v.depth)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression, v.depth)
}

// This private method visits the specified element.
func (v *visitor) visitElement(element Element) {
	switch actual := element.(type) {
	case Intrinsic:
		v.agent.AtIntrinsic(actual, v.depth)
	case String:
		v.agent.AtString(actual, v.depth)
	case Number:
		v.agent.AtNumber(actual, v.depth)
	case Name:
		v.agent.AtName(actual, v.depth)
	default:
		panic(fmt.Sprintf("Attempted to visit:\n    element: %v\n    type: %t\n", actual, element))
	}
}

// This private method visits the specified exactly N grouping.
func (v *visitor) visitExactlyN(group ExactlyNLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression, v.depth)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression, v.depth)
}

// This private method visits the specified expression.
func (v *visitor) visitExpression(expression ExpressionLike) {
	var isAnnotated = expression.IsAnnotated()
	var alternatives = expression.GetAlternatives()
	var size = alternatives.GetSize()
	var iterator = col.Iterator(alternatives)
	v.depth++
	for iterator.HasNext() {
		var slot = iterator.GetSlot()
		var alternative = iterator.GetNext()
		v.agent.BeforeAlternative(alternative, slot, size, isAnnotated, v.depth)
		v.visitAlternative(alternative)
		slot++
		v.agent.AfterAlternative(alternative, slot, size, isAnnotated, v.depth)
	}
	v.depth--
}

// This private method visits the specified factor.
func (v *visitor) visitFactor(factor Factor) {
	if ref.ValueOf(factor).Kind() == ref.String {
		v.agent.BeforeElement(factor, v.depth)
		v.visitElement(factor)
		v.agent.AfterElement(factor, v.depth)
		return
	}
	switch actual := factor.(type) {
	case *range_:
		v.agent.BeforeRange(actual, v.depth)
		v.visitRange(actual)
		v.agent.AfterRange(actual, v.depth)
	case *inverse:
		v.agent.BeforeInverse(actual, v.depth)
		v.visitInverse(actual)
		v.agent.AfterInverse(actual, v.depth)
	default:
		v.agent.BeforeGrouping(actual, v.depth)
		v.visitGrouping(actual)
		v.agent.AfterGrouping(actual, v.depth)
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
		v.agent.BeforeStatement(statement, slot, size, v.depth)
		v.visitStatement(statement)
		slot++
		v.agent.AfterStatement(statement, slot, size, v.depth)
	}
}

// This private method visits the specified grouping.
func (v *visitor) visitGrouping(grouping Grouping) {
	switch actual := grouping.(type) {
	case *exactlyN:
		var n = actual.GetN()
		v.agent.BeforeExactlyN(actual, n, v.depth)
		v.visitExactlyN(actual)
		v.agent.AfterExactlyN(actual, n, v.depth)
	case *zeroOrOne:
		v.agent.BeforeZeroOrOne(actual, v.depth)
		v.visitZeroOrOne(actual)
		v.agent.AfterZeroOrOne(actual, v.depth)
	case *zeroOrMore:
		v.agent.BeforeZeroOrMore(actual, v.depth)
		v.visitZeroOrMore(actual)
		v.agent.AfterZeroOrMore(actual, v.depth)
	case *oneOrMore:
		v.agent.BeforeOneOrMore(actual, v.depth)
		v.visitOneOrMore(actual)
		v.agent.AfterOneOrMore(actual, v.depth)
	default:
		panic(fmt.Sprintf("Attempted to visit:\n    grouping: %v\n    type: %t\n", actual, grouping))
	}
}

// This private method visits the specified inverse.
func (v *visitor) visitInverse(inverse InverseLike) {
	var factor = inverse.GetFactor()
	v.agent.BeforeFactor(factor, 0, 0, v.depth)
	v.visitFactor(factor)
	v.agent.AfterFactor(factor, 0, 0, v.depth)
}

// This private method visits the specified one or more grouping.
func (v *visitor) visitOneOrMore(group OneOrMoreLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression, v.depth)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression, v.depth)
}

// This private method visits the specified range.
func (v *visitor) visitRange(range_ RangeLike) {
	var first = range_.GetFirstCharacter()
	v.agent.AtCharacter(first, v.depth)
	var last = range_.GetLastCharacter()
	if len(last) > 0 {
		v.agent.BetweenCharacters(first, last, v.depth)
		v.agent.AtCharacter(last, v.depth)
	}
}

// This private method visits the specified statement.
func (v *visitor) visitStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.agent.AtComment(comment, v.depth)
	} else {
		var definition = statement.GetDefinition()
		v.agent.BeforeDefinition(definition, v.depth)
		v.visitDefinition(definition)
		v.agent.AfterDefinition(definition, v.depth)
	}
}

// This private method visits the specified zero or more grouping.
func (v *visitor) visitZeroOrMore(group ZeroOrMoreLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression, v.depth)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression, v.depth)
}

// This private method visits the specified zero or one grouping.
func (v *visitor) visitZeroOrOne(group ZeroOrOneLike) {
	var expression = group.GetExpression()
	v.agent.BeforeExpression(expression, v.depth)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression, v.depth)
}
