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

// ELEMENT INTERFACE

// This interface defines the methods supported by all element-like
// components.
type ElementLike interface {
	GetINTRINSIC() INTRINSIC
	SetINTRINSIC(intrinsic INTRINSIC)
	GetNAME() NAME
	SetNAME(name NAME)
	GetLITERAL() LITERAL
	SetLITERAL(literal LITERAL)
}

// This constructor creates a new element.
func Element(intrinsic INTRINSIC, name NAME, literal LITERAL) ElementLike {
	var v = &element{}
	v.SetINTRINSIC(intrinsic)
	v.SetNAME(name)
	v.SetLITERAL(literal)
	return v
}

// ELEMENT IMPLEMENTATION

// This type defines the structure and methods associated with a element.
type element struct {
	intrinsic INTRINSIC
	name      NAME
	literal   LITERAL
}

// This method returns the intrinsic for this element.
func (v *element) GetINTRINSIC() INTRINSIC {
	return v.intrinsic
}

// This method sets the intrinsic for this element.
func (v *element) SetINTRINSIC(intrinsic INTRINSIC) {
	if len(intrinsic) > 0 {
		v.intrinsic = intrinsic
		v.name = ""
		v.literal = ""
	}
}

// This method returns the name for this element.
func (v *element) GetNAME() NAME {
	return v.name
}

// This method sets the name for this element.
func (v *element) SetNAME(name NAME) {
	if len(name) > 0 {
		v.intrinsic = ""
		v.name = name
		v.literal = ""
	}
}

// This method returns the literal for this element.
func (v *element) GetLITERAL() LITERAL {
	return v.literal
}

// This method sets the literal for this element.
func (v *element) SetLITERAL(literal LITERAL) {
	if len(literal) > 0 {
		v.intrinsic = ""
		v.name = ""
		v.literal = literal
	}
}
