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

// DOCUMENT INTERFACE

// This interface defines the methods supported by all document-like
// components.
type DocumentLike interface {
	GetStatements() col.Sequential[StatementLike]
	SetStatements(statements col.Sequential[StatementLike])
}

// This constructor creates a new document.
func Document(statements col.Sequential[StatementLike]) DocumentLike {
	var v = &document{}
	v.SetStatements(statements)
	return v
}

// DOCUMENT IMPLEMENTATION

// This type defines the structure and methods associated with a document.
type document struct {
	statements col.Sequential[StatementLike]
}

// This method returns the statements for this document.
func (v *document) GetStatements() col.Sequential[StatementLike] {
	return v.statements
}

// This method sets the statements for this document.
func (v *document) SetStatements(statements col.Sequential[StatementLike]) {
	if statements == nil || statements.IsEmpty() {
		panic("A document requires at least one statement.")
	}
	v.statements = statements
}
