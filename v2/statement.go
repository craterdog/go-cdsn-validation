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

// STATEMENT INTERFACE

// This interface defines the methods supported by all statement-like
// components.
type StatementLike interface {
	GetCOMMENT() COMMENT
	SetCOMMENT(comment COMMENT)
	GetDefinition() DefinitionLike
	SetDefinition(definition DefinitionLike)
}

// This constructor creates a new statement.
func Statement(comment COMMENT, definition DefinitionLike) StatementLike {
	var v = &statement{}
	v.SetCOMMENT(comment)
	v.SetDefinition(definition)
	return v
}

// STATEMENT IMPLEMENTATION

// This type defines the structure and methods associated with a statement.
type statement struct {
	comment    COMMENT
	definition DefinitionLike
}

// This method returns the comment for this statement.
func (v *statement) GetCOMMENT() COMMENT {
	return v.comment
}

// This method sets the comment for this statement.
func (v *statement) SetCOMMENT(comment COMMENT) {
	v.comment = comment
}

// This method returns the definition for this statement.
func (v *statement) GetDefinition() DefinitionLike {
	return v.definition
}

// This method sets the definition for this statement.
func (v *statement) SetDefinition(definition DefinitionLike) {
	if len(v.comment) == 0 && definition == nil {
		panic("A statement requires either a comment or a definition.")
	}
	v.definition = definition
}
