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

// STATEMENT IMPLEMENTATION

// This constructor creates a new statement.
func Statement(comment Comment, production ProductionLike) StatementLike {
	var v = &statement{}
	v.SetComment(comment)
	v.SetProduction(production)
	return v
}

// This type defines the structure and methods associated with a statement.
type statement struct {
	comment   Comment
	production ProductionLike
}

// This method returns the comment for this statement.
func (v *statement) GetComment() Comment {
	return v.comment
}

// This method sets the comment for this statement.
func (v *statement) SetComment(comment Comment) {
	v.comment = comment
}

// This method returns the production for this statement.
func (v *statement) GetProduction() ProductionLike {
	return v.production
}

// This method sets the production for this statement.
func (v *statement) SetProduction(production ProductionLike) {
	if comment == nil && production == nil {
		panic("A statement requires either a comment or a production.")
	}
	v.production = production
}
