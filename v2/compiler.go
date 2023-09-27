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

// COMPILER INTERFACE

// This function compiles the specified grammar into its corresponding parser.
func CompileGrammar(packageName string, grammar GrammarLike) {
	var agent = &compiler{packageName, 0}
	VisitGrammar(agent, grammar)
	// add #package#.go file if one does not yet exist
	// generate scanner.go file
	// generate parser.go file
}

// COMPILER IMPLEMENTATION

// This type defines the structure and methods for a compiler agent.
type compiler struct {
	packageName string
	depth       int
}

// PRIVATE METHODS

// This public method increments the depth of the traversal by one.
func (v *compiler) IncrementDepth() {
	v.depth++
}

// This public method decrements the depth of the traversal by one.
func (v *compiler) DecrementDepth() {
	v.depth--
}

// This public method is called for each character token.
func (v *compiler) AtCharacter(character Character) {
}

// This public method is called for each comment token.
func (v *compiler) AtComment(comment Comment) {
}

// This public method is called for each intrinsic token.
func (v *compiler) AtIntrinsic(intrinsic Intrinsic) {
}

// This public method is called for each name token.
func (v *compiler) AtName(name Name) {
}

// This public method is called for each note token.
func (v *compiler) AtNote(note Note) {
}

// This public method is called for each number token.
func (v *compiler) AtNumber(number Number) {
}

// This public method is called for each string token.
func (v *compiler) AtString(string_ String) {
}

// This public method is called for each symbol token.
func (v *compiler) AtSymbol(symbol Symbol, isMultilined bool) {
}

// This public method is called before each alternative in an expression.
func (v *compiler) BeforeAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool) {
}

// This public method is called after each alternative in an expression.
func (v *compiler) AfterAlternative(alternative AlternativeLike, slot int, size int, isMultilined bool) {
}

// This public method is called before each definition.
func (v *compiler) BeforeDefinition(definition DefinitionLike) {
}

// This public method is called after each definition.
func (v *compiler) AfterDefinition(definition DefinitionLike) {
}

// This public method is called before each element.
func (v *compiler) BeforeElement(element Element) {
}

// This public method is called after each element.
func (v *compiler) AfterElement(element Element) {
}

// This public method is called before each exactly N grouping.
func (v *compiler) BeforeExactlyN(exactlyN ExactlyNLike, n Number) {
}

// This public method is called after each exactly N grouping.
func (v *compiler) AfterExactlyN(exactlyN ExactlyNLike, n Number) {
}

// This public method is called before each expression.
func (v *compiler) BeforeExpression(expression ExpressionLike) {
}

// This public method is called after each expression.
func (v *compiler) AfterExpression(expression ExpressionLike) {
}

// This public method is called before each factor in an alternative.
func (v *compiler) BeforeFactor(factor Factor, slot int, size int) {
}

// This public method is called after each factor in an alternative.
func (v *compiler) AfterFactor(factor Factor, slot int, size int) {
}

// This public method is called before the grammar.
func (v *compiler) BeforeGrammar(grammar GrammarLike) {
}

// This public method is called after the grammar.
func (v *compiler) AfterGrammar(grammar GrammarLike) {
}

// This public method is called before each grouping.
func (v *compiler) BeforeGrouping(grouping Grouping) {
}

// This public method is called after each grouping.
func (v *compiler) AfterGrouping(grouping Grouping) {
}

// This public method is called before each inverse factor.
func (v *compiler) BeforeInverse(inverse InverseLike) {
}

// This public method is called after each inverse factor.
func (v *compiler) AfterInverse(inverse InverseLike) {
}

// This public method is called before each one or more grouping.
func (v *compiler) BeforeOneOrMore(oneOrMore OneOrMoreLike) {
}

// This public method is called after each one or more grouping.
func (v *compiler) AfterOneOrMore(oneOrMore OneOrMoreLike) {
}

// This public method is called before each character range.
func (v *compiler) BeforeRange(range_ RangeLike) {
}

// This public method is called between the two two characters in a range.
func (v *compiler) BetweenCharacters(first Character, last Character) {
}

// This public method is called after each character range.
func (v *compiler) AfterRange(range_ RangeLike) {
}

// This public method is called before each statement in a grammar.
func (v *compiler) BeforeStatement(statement StatementLike, slot int, size int) {
}

// This public method is called after each statement in a grammar.
func (v *compiler) AfterStatement(statement StatementLike, slot int, size int) {
}

// This public method is called before each zero or more grouping.
func (v *compiler) BeforeZeroOrMore(zeroOrMore ZeroOrMoreLike) {
}

// This public method is called after each zero or more grouping.
func (v *compiler) AfterZeroOrMore(zeroOrMore ZeroOrMoreLike) {
}

// This public method is called before each zero or one grouping.
func (v *compiler) BeforeZeroOrOne(zeroOrOne ZeroOrOneLike) {
}

// This public method is called after each zero or one grouping.
func (v *compiler) AfterZeroOrOne(zeroOrOne ZeroOrOneLike) {
}
