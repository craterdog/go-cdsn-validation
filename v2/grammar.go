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
	col "github.com/craterdog/go-collection-framework/v2"
)

// GRAMMAR INTERFACE

// This interface defines the methods supported by all grammar-like
// components.
type GrammarLike interface {
	GetStatements() col.Sequential[StatementLike]
	SetStatements(statements col.Sequential[StatementLike])
}

// This constructor creates a new grammar.
func Grammar(statements col.Sequential[StatementLike]) GrammarLike {
	var v = &grammar{}
	v.SetStatements(statements)
	return v
}

// GRAMMAR IMPLEMENTATION

// This type defines the structure and methods associated with a grammar.
type grammar struct {
	statements col.Sequential[StatementLike]
}

// This method returns the statements for this grammar.
func (v *grammar) GetStatements() col.Sequential[StatementLike] {
	return v.statements
}

// This method sets the statements for this grammar.
func (v *grammar) SetStatements(statements col.Sequential[StatementLike]) {
	if statements == nil || statements.IsEmpty() {
		panic("A grammar requires at least one statement.")
	}
	v.statements = statements
}

// This private method appends a formatted grammar to the result.
func (v *formatter) formatGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.formatStatement(statement)
	}
}
