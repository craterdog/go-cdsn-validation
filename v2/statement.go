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
func Statement(annotation Annotation, production ProductionLike) StatementLike {
	var v = &statement{}
	v.SetAnnotation(annotation)
	v.SetProduction(production)
	return v
}

// This type defines the structure and methods associated with a statement.
type statement struct {
	annotation Annotation
	production ProductionLike
}

// This method returns the annotation for this statement.
func (v *statement) GetAnnotation() Annotation {
	return v.annotation
}

// This method sets the annotation for this statement.
func (v *statement) SetAnnotation(annotation Annotation) {
	v.annotation = annotation
}

// This method returns the production for this statement.
func (v *statement) GetProduction() ProductionLike {
	return v.production
}

// This method sets the production for this statement.
func (v *statement) SetProduction(production ProductionLike) {
	if len(v.annotation) == 0 && production == nil {
		panic("A statement requires either an annotation or a production.")
	}
	v.production = production
}
