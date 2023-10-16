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
	GetDefinition() DefinitionLike
	SetDefinition(definition DefinitionLike)
	GetCOMMENT() COMMENT
	SetCOMMENT(comment COMMENT)
}

// This constructor creates a new statement.
func Statement(definition DefinitionLike, comment COMMENT) StatementLike {
	if definition == nil && len(comment) == 0 {
		panic("A statement requires at least one of its attributes to be set.")
	}
	var v = &statement{}
	v.SetDefinition(definition)
	v.SetCOMMENT(comment)
	return v
}

// STATEMENT IMPLEMENTATION

// This type defines the structure and methods associated with an statement.
type statement struct {
	definition DefinitionLike
	comment    COMMENT
}

// This method returns the definition for this statement.
func (v *statement) GetDefinition() DefinitionLike {
	return v.definition
}

// This method sets the definition for this statement.
func (v *statement) SetDefinition(definition DefinitionLike) {
	if definition != nil {
		v.definition = definition
		v.comment = ""
	}
}

// This method returns the comment for this statement.
func (v *statement) GetCOMMENT() COMMENT {
	return v.comment
}

// This method sets the comment for this statement.
func (v *statement) SetCOMMENT(comment COMMENT) {
	if len(comment) > 0 {
		v.definition = nil
		v.comment = comment
	}
}
