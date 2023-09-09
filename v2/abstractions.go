/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
This cdsn package defines a parser and a canonical formatter for language
grammars containing Crater Dog Syntax Notation™ (CDSN).
*/
package cdsn

import (
	col "github.com/craterdog/go-collection-framework/v2"
)

// TYPE DEFINITIONS

// The following define the native Go token types.
type (
	Annotation   string
	Comment      string
	Digit        string
	Identifier   string
	Intrinsic    string
	Letter       string
	Note         string
	Number       string
	Rune         string
	String       string
	Symbol       string
)

// The following define the Go rule related types.
type (
	Factor       any
	GroupingType int64
)

// CONSTANTS

// The POSIX standard end-of-line character.
const EOL = "\n"

// The allowed grouping types.
const (
	Optional GroupingType = iota
	ExactNumber
	MinimumNumber
	MaximumNumber
)

// INDIVIDUAL INTERFACES

// This interface defines the methods supported by all alternative-like
// components.
type AlternativeLike interface {
	GetFactors() col.Sequential[Factor]
	SetFactors(factors col.Sequential[Factor])
	GetNote() Note
	SetNote(note Note)
}

// This interface defines the methods supported by all definition-like
// components.
type DefinitionLike interface {
	IsMultilined() bool
	SetMultilined(multilined bool)
	GetAlternatives() col.Sequential[AlternativeLike]
	SetAlternatives(alternatives col.Sequential[AlternativeLike])
}

// This interface defines the methods supported by all grammar-like
// components.
type GrammarLike interface {
	GetStatements() col.Sequential[StatementLike]
	SetStatements(statements col.Sequential[StatementLike])
}

// This interface defines the methods supported by all grouping-like
// components.
type GroupingLike interface {
	GetDefinition() DefinitionLike
	SetDefinition(definition DefinitionLike)
	GetType() GroupingType
	SetType(type_ GroupingType)
	GetNumber() Number
	SetNumber(count Number)
}

// This interface defines the methods supported by all inversion-like
// components.
type InversionLike interface {
	GetFactor() Factor
	SetFactor(factor Factor)
}

// This interface defines the methods supported by all production-like
// components.
type ProductionLike interface {
	GetSymbol() Symbol
	SetSymbol(symbol Symbol)
	GetDefinition() DefinitionLike
	SetDefinition(definition DefinitionLike)
}

// This interface defines the methods supported by all range-like components.
type RangeLike interface {
	GetFirstRune() Rune
	SetFirstRune(first Rune)
	GetLastRune() Rune
	SetLastRune(last Rune)
}

// This interface defines the methods supported by all statement-like
// components.
type StatementLike interface {
	GetAnnotation() Annotation
	SetAnnotation(annotation Annotation)
	GetProduction() ProductionLike
	SetProduction(production ProductionLike)
}
