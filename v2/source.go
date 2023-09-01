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

// SOURCE IMPLEMENTATION

// This constructor creates a new source.
func Source(statements col.Sequential[StatementLike]) SourceLike {
	var v = &source{}
	v.SetStatements(statements)
	return v
}

// This type defines the structure and methods associated with a source.
type source struct {
	statements col.Sequential[StatementLike]
}

// This method returns the statements for this source.
func (v *source) GetStatements() col.Sequential[StatementLike] {
	return v.statements
}

// This method sets the statements for this source.
func (v *source) SetStatements(statements col.Sequential[StatementLike]) {
	if statements == nil || statements.IsEmpty() {
		panic("A source requires at least one statement.")
	}
	v.statements = statements
}
