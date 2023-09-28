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
	byt "bytes"
	osx "os"
)

// COMPILER INTERFACE

// This function compiles the specified grammar into its corresponding parser.
func CompileGrammar(directory, packageName string, grammar GrammarLike) {
	var agent = Compiler(directory, packageName)
	VisitGrammar(agent, grammar)
}

type CompilerLike interface {
	Specialized
}

func Compiler(directory, packageName string) CompilerLike {
	var v = &compiler{directory: directory, packageName: packageName}
	v.initializeConfiguration()
	v.initializeScanner()
	v.initializeParser()
	v.initializeVisitor()
	return v
}

// COMPILER IMPLEMENTATION

// This type defines the structure and methods for a compiler agent.
type compiler struct {
	directory     string
	packageName   string
	scannerBuffer byt.Buffer
	parserBuffer  byt.Buffer
	visitorBuffer byt.Buffer
	depth         int
}

// PUBLIC METHODS

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

// PRIVATE METHODS

// This private method creates a new configuration (package.go) file if one
// does not already exist.
func (v *compiler) initializeConfiguration() {
	var err error
	var template []byte
	var configuration = v.directory + "/package.go"
	_, err = osx.Open(configuration)
	if err != nil {
		template, err = osx.ReadFile("./templates/package.tp")
		if err != nil {
			panic(err)
		}
		template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
		err = osx.WriteFile(configuration, template, 0666)
		if err != nil {
			panic(err)
		}
	}
}

// This private method creates the byte buffer for the generated scanner code.
func (v *compiler) initializeScanner() {
	var template, err = osx.ReadFile("./templates/scanner.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.scannerBuffer.Write(template)
}

// This private method creates the byte buffer for the generated parser code.
func (v *compiler) initializeParser() {
	var template, err = osx.ReadFile("./templates/parser.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.parserBuffer.Write(template)
}

// This private method creates the byte buffer for the generated visitor code.
func (v *compiler) initializeVisitor() {
	var template, err = osx.ReadFile("./templates/visitor.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.visitorBuffer.Write(template)
}
