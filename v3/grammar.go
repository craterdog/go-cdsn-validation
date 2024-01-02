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
	col "github.com/craterdog/go-collection-framework/v3"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type grammarClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var grammarClass = &grammarClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func GrammarClass() GrammarClassLike {
	return grammarClass
}

// Public Class Constructors

func (c *grammarClass_) FromStatements(statements col.Sequential[StatementLike]) GrammarLike {
	var grammar = &grammar_{
		// This class does not initialize any attributes.
	}
	grammar.SetStatements(statements)
	return grammar
}

// CLASS INSTANCES

// Private Class Type Definition

type grammar_ struct {
	statements col.Sequential[StatementLike]
}

// Public Interface

func (v *grammar_) GetStatements() col.Sequential[StatementLike] {
	return v.statements
}

func (v *grammar_) SetStatements(statements col.Sequential[StatementLike]) {
	if statements == nil || statements.IsEmpty() {
		panic("An grammar must have at least one statement.")
	}
	v.statements = statements
}
