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
)

// VISITOR INTERFACE

// This function applies the specified agent to each node in the specified
// document.
func VisitDocument(agent Specialized, document DocumentLike) {
	var v = &visitor{agent, 0}
	v.agent.BeforeDocument(document)
	v.visitDocument(document)
	v.agent.AfterDocument(document)
}

// This interface defines the methods that are supported by specialized agents.
type Specialized interface {
	IncrementDepth()
	DecrementDepth()
	AtCHARACTER(character CHARACTER)
	BetweenCHARACTERs(first CHARACTER, last CHARACTER)
	AtCOMMENT(comment COMMENT)
	AtCONSTRAINT(constraint CONSTRAINT)
	AtINTRINSIC(intrinsic INTRINSIC)
	AtNAME(name NAME)
	AtNOTE(note NOTE)
	AtSTRING(string_ STRING)
	AtSYMBOL(symbol SYMBOL, isMultiline bool)
	BeforeAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool)
	AfterAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool)
	BeforeDefinition(definition DefinitionLike)
	AfterDefinition(definition DefinitionLike)
	BeforeDocument(document DocumentLike)
	AfterDocument(document DocumentLike)
	BeforeElement(element Element)
	AfterElement(element Element)
	BeforeExpression(expression ExpressionLike)
	AfterExpression(expression ExpressionLike)
	BeforeFactor(factor Factor)
	AfterFactor(factor Factor)
	BeforePrecedence(precedence PrecedenceLike)
	AfterPrecedence(precedence PrecedenceLike)
	BeforePredicate(predicate Predicate, slot int, size int)
	AfterPredicate(predicate Predicate, slot int, size int)
	BeforeRange(range_ RangeLike)
	AfterRange(range_ RangeLike)
	BeforeRepetition(repetition RepetitionLike)
	AfterRepetition(repetition RepetitionLike)
	BeforeStatement(statement Statement, slot int, size int)
	AfterStatement(statement Statement, slot int, size int)
}

// VISITOR IMPLEMENTATION

// This type defines the structure and methods for a visitor.
type visitor struct {
	agent Specialized
	depth int
}

// This private method visits the specified alternative.
func (v *visitor) visitAlternative(alternative AlternativeLike) {
	var predicates = alternative.GetPredicates()
	var size = predicates.GetSize()
	var iterator = col.Iterator(predicates)
	for iterator.HasNext() {
		var slot = iterator.GetSlot()
		var predicate = iterator.GetNext()
		v.agent.BeforePredicate(predicate, slot, size)
		v.visitPredicate(predicate)
		slot++
		v.agent.AfterPredicate(predicate, slot, size)
	}
	var note = alternative.GetNOTE()
	if len(note) > 0 {
		v.agent.AtNOTE(note)
	}
}

// This private method visits the specified definition.
func (v *visitor) visitDefinition(definition DefinitionLike) {
	var symbol = definition.GetSYMBOL()
	var expression = definition.GetExpression()
	v.agent.AtSYMBOL(symbol, expression.IsMultilined())
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}

// This private method visits the specified document.
func (v *visitor) visitDocument(document DocumentLike) {
	var statements = document.GetStatements()
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

// This private method visits the specified element.
func (v *visitor) visitElement(element Element) {
	switch actual := element.(type) {
	case INTRINSIC:
		v.agent.AtINTRINSIC(actual)
	case NAME:
		v.agent.AtNAME(actual)
	case STRING:
		v.agent.AtSTRING(actual)
	default:
		panic(fmt.Sprintf("Attempted to visit:\n    element: %v\n    type: %t\n", actual, element))
	}
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
	switch actual := factor.(type) {
	case *precedence:
		v.agent.BeforePrecedence(actual)
		v.visitPrecedence(actual)
		v.agent.AfterPrecedence(actual)
	default:
		v.agent.BeforeElement(actual)
		v.visitElement(actual)
		v.agent.AfterElement(actual)
	}
}

// This private method visits the specified precedence.
func (v *visitor) visitPrecedence(definition PrecedenceLike) {
	var expression = definition.GetExpression()
	v.agent.BeforeExpression(expression)
	v.visitExpression(expression)
	v.agent.AfterExpression(expression)
}

// This private method visits the specified predicate.
func (v *visitor) visitPredicate(predicate Predicate) {
	switch actual := predicate.(type) {
	case *range_:
		v.agent.BeforeRange(actual)
		v.visitRange(actual)
		v.agent.AfterRange(actual)
	case *repetition:
		v.agent.BeforeRepetition(actual)
		v.visitRepetition(actual)
		v.agent.AfterRepetition(actual)
	default:
		v.agent.BeforeFactor(actual)
		v.visitFactor(actual)
		v.agent.AfterFactor(actual)
	}
}

// This private method visits the specified range.
func (v *visitor) visitRange(range_ RangeLike) {
	var first = range_.GetFirstCHARACTER()
	v.agent.AtCHARACTER(first)
	var last = range_.GetLastCHARACTER()
	if len(last) > 0 {
		v.agent.BetweenCHARACTERs(first, last)
		v.agent.AtCHARACTER(last)
	}
}

// This private method visits the specified repetition.
func (v *visitor) visitRepetition(repetition RepetitionLike) {
	var constraint = repetition.GetCONSTRAINT()
	v.agent.AtCONSTRAINT(constraint)
	var factor = repetition.GetFactor()
	v.agent.BeforeFactor(factor)
	v.visitFactor(factor)
	v.agent.AfterFactor(factor)
}

// This private method visits the specified statement.
func (v *visitor) visitStatement(statement Statement) {
	switch actual := statement.(type) {
	case *definition:
		v.agent.BeforeDefinition(actual)
		v.visitDefinition(actual)
		v.agent.AfterDefinition(actual)
	case COMMENT:
		v.agent.AtCOMMENT(actual)
	}
}
