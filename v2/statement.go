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
	GetComment() Comment
	SetComment(comment Comment)
	GetDefinition() DefinitionLike
	SetDefinition(definition DefinitionLike)
}

// This constructor creates a new statement.
func Statement(comment Comment, definition DefinitionLike) StatementLike {
	var v = &statement{}
	v.SetComment(comment)
	v.SetDefinition(definition)
	return v
}

// STATEMENT IMPLEMENTATION

// This type defines the structure and methods associated with a statement.
type statement struct {
	comment    Comment
	definition DefinitionLike
}

// This method returns the comment for this statement.
func (v *statement) GetComment() Comment {
	return v.comment
}

// This method sets the comment for this statement.
func (v *statement) SetComment(comment Comment) {
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

// This method attempts to parse a statement. It returns the statement and
// whether or not the statement was successfully parsed.
func (v *parser) parseStatement() (StatementLike, *Token, bool) {
	var ok bool
	var token *Token
	var comment Comment
	var definition DefinitionLike
	var statement StatementLike
	comment, _, ok = v.parseComment()
	if !ok {
		definition, token, ok = v.parseDefinition()
		if !ok {
			// This is not a statement.
			return statement, token, false
		}
	}
	statement = Statement(comment, definition)
	return statement, token, true
}

// This private method appends a formatted statement to the result.
func (v *formatter) formatStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.appendNewline()
		v.formatComment(comment)
	} else {
		var definition = statement.GetDefinition()
		v.formatDefinition(definition)
	}
	v.appendNewline()
	v.appendNewline()
}
