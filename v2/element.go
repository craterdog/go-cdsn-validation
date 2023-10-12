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
	SetINTRINSIC(intrinsic_ INTRINSIC)
	GetNAME() NAME
	SetNAME(name_ NAME)
	GetSTRING() STRING
	SetSTRING(string_ STRING)
}

// This constructor creates a new element.
func Element(intrinsic_ INTRINSIC, name_ NAME, string_ STRING) ElementLike {
	var v = &element_{}
	v.SetINTRINSIC(intrinsic_)
	v.SetNAME(name_)
	v.SetSTRING(string_)
	return v
}

// ELEMENT IMPLEMENTATION

// This type defines the structure and methods associated with a element.
type element_ struct {
	intrinsic_ INTRINSIC
	name_      NAME
	string_    STRING
}

// This method returns the intrinsic for this element.
func (v *element_) GetINTRINSIC() INTRINSIC {
	return v.intrinsic_
}

// This method sets the intrinsic for this element.
func (v *element_) SetINTRINSIC(intrinsic_ INTRINSIC) {
	if len(intrinsic_) > 0 {
		v.intrinsic_ = intrinsic_
		v.name_ = ""
		v.string_ = ""
	}
}

// This method returns the name for this element.
func (v *element_) GetNAME() NAME {
	return v.name_
}

// This method sets the name for this element.
func (v *element_) SetNAME(name_ NAME) {
	if len(name_) > 0 {
		v.intrinsic_ = ""
		v.name_ = name_
		v.string_ = ""
	}
}

// This method returns the string for this element.
func (v *element_) GetSTRING() STRING {
	return v.string_
}

// This method sets the string for this element.
func (v *element_) SetSTRING(string_ STRING) {
	if len(string_) > 0 {
		v.intrinsic_ = ""
		v.name_ = ""
		v.string_ = string_
	}
}
