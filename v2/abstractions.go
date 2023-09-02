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
	Comment      string
	Identifier   string
	Intrinsic    string
	Factor       any
	Literal      string
	Note         string
	Range        string
	Symbol       string
	GroupingType int64
)

const (
	Precedence GroupingType = iota
	ZeroOrOne
	ZeroOrMore
	OneOrMore
)

// INDIVIDUAL INTERFACES

// This interface defines the methods supported by all alternative-like components.
type AlternativeLike interface {
	GetOption() OptionLike
	SetOption(option OptionLike)
	GetNote() Note
	SetNote(note Note)
}

// This interface defines the methods supported by all grammar-like components.
type GrammarLike interface {
	GetStatements() col.Sequential[StatementLike]
	SetStatements(statements col.Sequential[StatementLike])
}

// This interface defines the methods supported by all grouping-like components.
type GroupingLike interface {
	GetRule() RuleLike
	SetRule(rule RuleLike)
	GetType() GroupingType
	SetType(type_ GroupingType)
}

// This interface defines the methods supported by all inversion-like components.
type InversionLike interface {
	GetFactor() Factor
	SetFactor(factor Factor)
}

// This interface defines the methods supported by all option-like components.
type OptionLike interface {
	GetFactors() col.Sequential[Factor]
	SetFactors(factors col.Sequential[Factor])
}

// This interface defines the methods supported by all production-like components.
type ProductionLike interface {
	GetSymbol() Symbol
	SetSymbol(symbol Symbol)
	GetRule() RuleLike
	SetRule(rule RuleLike)
	GetNote() Note
	SetNote(note Note)
}

// This interface defines the methods supported by all rule-like components.
type RuleLike interface {
	GetOption() OptionLike
	SetOption(option OptionLike)
	GetAlternatives() col.Sequential[AlternativeLike]
	SetAlternatives(alternatives col.Sequential[AlternativeLike])
}

// This interface defines the methods supported by all source-like components.
type SourceLike interface {
	GetStatements() col.Sequential[StatementLike]
	SetStatements(statements col.Sequential[StatementLike])
}

// This interface defines the methods supported by all statement-like components.
type StatementLike interface {
	GetComment() Comment
	SetComment(comment Comment)
	GetProduction() ProductionLike
	SetProduction(production ProductionLike)
}
