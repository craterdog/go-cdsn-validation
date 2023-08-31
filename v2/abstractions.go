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

type (
	Comment    string
	Identifier string
	Intrinsic  string
	Factor     any
	Literal    string
	Note       string
	Symbol     string
)

// INDIVIDUAL INTERFACES

// This interface defines the methods supported by all alternative-like components.
type AlternativeLike interface {
	GetNote() Note
	GetOption() OptionLike
}

// This interface defines the methods supported by all grammar-like components.
type GrammarLike interface {
	GetStatements() col.Sequential[StatementLike]
}

// This interface defines the methods supported by all group-like components.
type GroupLike interface {
	GetLeftBracket() string
	GetRightBracket() string
	GetRule() RuleLike
}

// This interface defines the methods supported by all inversion-like components.
type InversionLike interface {
	GetFactor() Factor
}

// This interface defines the methods supported by all option-like components.
type OptionLike interface {
	GetFactors() col.Sequential[Factor]
}

// This interface defines the methods supported by all production-like components.
type ProductionLike interface {
	GetNote() Note
	GetRule() RuleLike
	GetSymbol() Symbol
}

// This interface defines the methods supported by all range-like components.
type RangeLike interface {
	GetFirstLiteral() Literal
	GetLastLiteral() Literal
}

// This interface defines the methods supported by all rule-like components.
type RuleLike interface {
	GetOption() OptionLike
	GetAlternatives() col.Sequential[AlternativeLike]
}

// This interface defines the methods supported by all statement-like components.
type StatementLike interface {
	GetComment() Comment
	GetProduction() ProductionLike
}
