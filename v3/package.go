/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
This package file defines the INTERFACE to this package.  Any additions to the
types defined in this file require a MINOR version change.  Any deletions from,
or changes to, the types defined in this file require a MAJOR version change.

The package provides a parser and formatter for documents written using Crater
Dog Syntax Notation™ (CDSN).  The parser performs validation on the resulting
parse tree.  The formatter takes a validated parse tree and generates the
corresponding CDSN document using the canonical format.

For detailed documentation on this package refer to the wiki:

	https://github.com/craterdog/go-cdsn-validation/wiki

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:

	https://github.com/craterdog/go-coding-conventions/wiki

Additional implementations of the classes provided by this package can be
developed and used seamlessly since the interface definitions only depend on
other interfaces and primitive types; and the class implementations only depend
on interfaces, not on each other.
*/
package cdsn

import (
	col "github.com/craterdog/go-collection-framework/v3"
)

// PACKAGE ABSTRACTIONS

// Abstract Types

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all alternative-class-like types.
type AlternativeClassLike interface {
	FromFactors(factors col.Sequential[FactorLike]) AlternativeLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all alternative-like types.
type AlternativeLike interface {
	GetFactors() col.Sequential[FactorLike]
	GetNote() string
	SetFactors(factors col.Sequential[FactorLike])
	SetNote(note string)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all assertion-class-like types.
type AssertionClassLike interface {
	FromElement(element ElementLike) AssertionLike
	FromGlyph(glyph GlyphLike) AssertionLike
	FromPrecedence(precedence PrecedenceLike) AssertionLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all assertion-like types.
type AssertionLike interface {
	GetElement() ElementLike
	GetGlyph() GlyphLike
	GetPrecedence() PrecedenceLike
	SetElement(element ElementLike)
	SetGlyph(glyph GlyphLike)
	SetPrecedence(precedence PrecedenceLike)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all cardinality-class-like types.
type CardinalityClassLike interface {
	FromConstraint(constraint ConstraintLike) CardinalityLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all cardinality-like types.
type CardinalityLike interface {
	GetConstraint() ConstraintLike
	SetConstraint(constraint ConstraintLike)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all constraint-class-like types.
type ConstraintClassLike interface {
	FromNumber(number string) ConstraintLike
	FromRange(first, last string) ConstraintLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all constraint-like types.
type ConstraintLike interface {
	GetFirst() string
	GetLast() string
	SetFirst(first string)
	SetLast(last string)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all definition-class-like types.
type DefinitionClassLike interface {
	FromSymbolAndExpression(symbol string, expression ExpressionLike) DefinitionLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all definition-like types.
type DefinitionLike interface {
	GetExpression() ExpressionLike
	GetSymbol() string
	SetExpression(expression ExpressionLike)
	SetSymbol(symbol string)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all document-class-like types.
type DocumentClassLike interface {
	FromGrammar(grammar GrammarLike) DocumentLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all document-like types.
type DocumentLike interface {
	GetGrammar() GrammarLike
	SetGrammar(grammar GrammarLike)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all element-class-like types.
type ElementClassLike interface {
	FromIntrinsic(intrinsic string) ElementLike
	FromName(name string) ElementLike
	FromLiteral(literal string) ElementLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all element-like types.
type ElementLike interface {
	GetIntrinsic() string
	GetName() string
	GetLiteral() string
	SetIntrinsic(intrinsic string)
	SetName(name string)
	SetLiteral(literal string)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all expression-class-like types.
type ExpressionClassLike interface {
	FromAlternatives(alternatives col.Sequential[AlternativeLike]) ExpressionLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all expression-like types.
type ExpressionLike interface {
	GetAlternatives() col.Sequential[AlternativeLike]
	IsMultilined() bool
	SetAlternatives(alternatives col.Sequential[AlternativeLike])
	SetMultilined(isMultilined bool)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all factor-class-like types.
type FactorClassLike interface {
	FromPredicate(predicate PredicateLike) FactorLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all factor-like types.
type FactorLike interface {
	GetCardinality() CardinalityLike
	GetPredicate() PredicateLike
	SetCardinality(cardinality CardinalityLike)
	SetPredicate(predicate PredicateLike)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all formatter-class-like types.
type FormatterClassLike interface {
	Default() FormatterLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all formatter-like types.
type FormatterLike interface {
	FormatDocument(document DocumentLike) string
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all glyph-class-like types.
type GlyphClassLike interface {
	FromCharacter(character string) GlyphLike
	FromRange(first, last string) GlyphLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all glyph-like types.
type GlyphLike interface {
	GetFirst() string
	GetLast() string
	SetFirst(first string)
	SetLast(last string)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all grammar-class-like types.
type GrammarClassLike interface {
	FromStatements(statements col.Sequential[StatementLike]) GrammarLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all grammar-like types.
type GrammarLike interface {
	GetStatements() col.Sequential[StatementLike]
	SetStatements(statements col.Sequential[StatementLike])
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all parser-class-like types.
type ParserClassLike interface {
	Default() ParserLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all parser-like types.
type ParserLike interface {
	ParseDocument(document string) DocumentLike
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all precedence-class-like types.
type PrecedenceClassLike interface {
	FromExpression(expression ExpressionLike) PrecedenceLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all precedence-like types.
type PrecedenceLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all predicate-class-like types.
type PredicateClassLike interface {
	FromAssertion(assertion AssertionLike, inverted bool) PredicateLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all predicate-like types.
type PredicateLike interface {
	GetAssertion() AssertionLike
	IsInverted() bool
	SetAssertion(assertion AssertionLike)
	SetInverted(inverted bool)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all statement-class-like types.
type StatementClassLike interface {
	FromComment(comment string) StatementLike
	FromDefinition(definition DefinitionLike) StatementLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all statement-like types.
type StatementLike interface {
	GetComment() string
	GetDefinition() DefinitionLike
	SetComment(comment string)
	SetDefinition(definition DefinitionLike)
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all validator-class-like types.
type ValidatorClassLike interface {
	Default() ValidatorLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all validator-like types.
type ValidatorLike interface {
	ValidateDocument(document DocumentLike)
}
