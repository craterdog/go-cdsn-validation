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
	col "github.com/craterdog/go-collection-framework/v3"
	stc "strconv"
	uni "unicode"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type validatorClass_ struct {
	// This class does not define any class constants.
}

// Private Class Namespace Reference

var validatorClass = &validatorClass_{
	// This class does not initialize any class constants.
}

// Public Class Namespace Access

func ValidatorClass() ValidatorClassLike {
	return validatorClass
}

// Public Class Constructors

func (c *validatorClass_) Default() ValidatorLike {
	var validator = &validator_{}
	return validator
}

// CLASS INSTANCES

// Private Class Type Definition

type validator_ struct {
	definitions col.StackLike[DefinitionLike]
	inInversion bool
	isToken     bool
}

// Public Interface

func (v *validator_) ValidateDocument(document DocumentLike) {
	v.definitions = col.StackClass[DefinitionLike]().Empty()
	var grammar = document.GetGrammar()
	v.validateGrammar(grammar)
}

// Private Interface

func (v *validator_) formatError(message string) string {
	var definition = v.definitions.RemoveTop()
	message = fmt.Sprintf(
		"The definition for %v is invalid:\n%v\n",
		definition.GetSymbol(),
		message,
	)
	return message
}

func (v *validator_) validateAlternative(alternative AlternativeLike) {
	var factors = alternative.GetFactors()
	if factors == nil || factors.IsEmpty() {
		var message = v.formatError(
			"Each alternative must have at least one factor.",
		)
		panic(message)
	}
	var iterator = factors.GetIterator()
	for iterator.HasNext() {
		var factor = iterator.GetNext()
		v.validateFactor(factor)
	}
	var note = alternative.GetNote()
	if len(note) > 0 {
		v.validateNote(note)
	}
}

func (v *validator_) validateAssertion(assertion AssertionLike) {
	var element = assertion.GetElement()
	var glyph = assertion.GetGlyph()
	var precedence = assertion.GetPrecedence()
	switch {
	case element != nil && glyph == nil && precedence == nil:
		v.validateElement(element)
	case element == nil && glyph != nil && precedence == nil:
		v.validateGlyph(glyph)
	case element == nil && glyph == nil && precedence != nil:
		v.validatePrecedence(precedence)
	default:
		var message = v.formatError(
			"An assertion must contain exactly one element, glyph, or precedence.",
		)
		panic(message)
	}
}

func (v *validator_) validateCharacter(character string) {
	var matches = ScannerClass().MatchCharacter(character)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid character.",
		)
		panic(message)
	}
}

func (v *validator_) validateCardinality(cardinality CardinalityLike) {
	var constraint = cardinality.GetConstraint()
	if constraint == nil {
		var message = v.formatError(
			"A cardinality must have a constraint.",
		)
		panic(message)
	}
	v.validateConstraint(constraint)
}

func (v *validator_) validateComment(comment string) {
	var matches = ScannerClass().MatchComment(comment)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid comment.",
		)
		panic(message)
	}
}

func (v *validator_) validateConstraint(constraint ConstraintLike) {
	var first = constraint.GetFirst()
	v.validateNumber(first)
	var last = constraint.GetLast()
	if len(last) > 0 {
		v.validateNumber(last)
		var firstNumber, _ = stc.ParseInt(first, 10, 64)
		var lastNumber, _ = stc.ParseInt(last, 10, 64)
		if firstNumber > lastNumber {
			var message = v.formatError(
				"The first number in a constraint cannot be greater than the last.",
			)
			panic(message)
		}
	}
}

func (v *validator_) validateDefinition(definition DefinitionLike) {
	v.definitions.AddValue(definition)
	var symbol = definition.GetSymbol()
	v.validateSymbol(symbol)
	var expression = definition.GetExpression()
	if expression == nil {
		var message = v.formatError(
			"A definition must contain an expression.",
		)
		panic(message)
	}
	v.validateExpression(expression)
	v.inInversion = false
	var _ = v.definitions.RemoveTop()
}

func (v *validator_) validateElement(element ElementLike) {
	var intrinsic = element.GetIntrinsic()
	var name = element.GetName()
	var literal = element.GetLiteral()
	switch {
	case len(intrinsic) > 0 && len(name) == 0 && len(literal) == 0:
		v.validateIntrinsic(intrinsic)
	case len(intrinsic) == 0 && len(name) > 0 && len(literal) == 0:
		v.validateName(name)
	case len(intrinsic) == 0 && len(name) == 0 && len(literal) > 0:
		v.validateLiteral(literal)
	default:
		var message = v.formatError(
			"An element must contain exactly one intrinsic, name, or literal.",
		)
		panic(message)
	}
}

func (v *validator_) validateExpression(expression ExpressionLike) {
	var alternatives = expression.GetAlternatives()
	if alternatives == nil || alternatives.IsEmpty() {
		var message = v.formatError(
			"Each expression must have at least one alternative.",
		)
		panic(message)
	}
	var iterator = alternatives.GetIterator()
	for iterator.HasNext() {
		var alternative = iterator.GetNext()
		v.validateAlternative(alternative)
	}
}

func (v *validator_) validateFactor(factor FactorLike) {
	var predicate = factor.GetPredicate()
	if predicate == nil {
		var message = v.formatError(
			"A factor must contain a predicate.",
		)
		panic(message)
	}
	v.validatePredicate(predicate)
	var cardinality = factor.GetCardinality()
	if cardinality != nil {
		v.validateCardinality(cardinality)
	}
}

func (v *validator_) validateGlyph(glyph GlyphLike) {
	var first = glyph.GetFirst()
	v.validateCharacter(first)
	var last = glyph.GetLast()
	if len(last) > 0 {
		v.validateCharacter(last)
		if first > last {
			var message = v.formatError(
				"The first character in a glyph cannot come later than the last.",
			)
			panic(message)
		}
	}
}

func (v *validator_) validateGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	if statements == nil || statements.IsEmpty() {
		var message = "The grammar must contain at least one statement.\n"
		panic(message)
	}
	var iterator = statements.GetIterator()
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.validateStatement(statement)
	}
}

func (v *validator_) validateIntrinsic(intrinsic string) {
	var matches = ScannerClass().MatchIntrinsic(intrinsic)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid intrinsic.",
		)
		panic(message)
	}
}

func (v *validator_) validateLiteral(literal string) {
	var matches = ScannerClass().MatchLiteral(literal)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid literal.",
		)
		panic(message)
	}
	if v.inInversion && len([]rune(literal)) > 3 {
		var message = v.formatError(
			"A multi-character literal is not allowed in an inversion.",
		)
		panic(message)
	}
}

func (v *validator_) validateName(name string) {
	var matches = ScannerClass().MatchName(name)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid name.",
		)
		panic(message)
	}
	if uni.IsLower([]rune(name)[0]) {
		if v.isToken {
			var message = v.formatError(
				"A token definition cannot contain a rule name.",
			)
			panic(message)
		}
		if v.inInversion {
			var message = v.formatError(
				"An inverted assertion cannot contain a rule name.",
			)
			panic(message)
		}
	}
}

func (v *validator_) validateNote(note string) {
	var matches = ScannerClass().MatchNote(note)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid note.",
		)
		panic(message)
	}
}

func (v *validator_) validateNumber(number string) {
	var matches = ScannerClass().MatchNumber(number)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid number.",
		)
		panic(message)
	}
}

func (v *validator_) validatePrecedence(precedence PrecedenceLike) {
	var expression = precedence.GetExpression()
	if precedence == nil {
		var message = v.formatError(
			"A precedence must contain an expression.",
		)
		panic(message)
	}
	v.validateExpression(expression)
}

func (v *validator_) validatePredicate(predicate PredicateLike) {
	var isInverted = predicate.IsInverted()
	if isInverted && v.inInversion {
		var message = v.formatError(
			"Inverted assertions cannot be nested.",
		)
		panic(message)
	}
	if isInverted {
		v.inInversion = true
	}
	var assertion = predicate.GetAssertion()
	if assertion == nil {
		var message = v.formatError(
			"A predicate must have an assertion.",
		)
		panic(message)
	}
	v.validateAssertion(assertion)
}

func (v *validator_) validateStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.validateComment(comment)
	} else {
		var definition = statement.GetDefinition()
		if definition == nil {
			panic("A statement must contain either a comment or a definition.")
		}
		v.validateDefinition(definition)
	}
}

func (v *validator_) validateSymbol(symbol string) {
	var matches = ScannerClass().MatchSymbol(symbol)
	if len(matches) == 0 {
		var message = v.formatError(
			"Found an invalid symbol.",
		)
		panic(message)
	}
	v.isToken = uni.IsUpper([]rune(matches[1])[0])
}
