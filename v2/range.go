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

// RANGE IMPLEMENTATION

// This constructor creates a new range.
func Range(first Literal, last Literal) RangeLike {
	var v = &range_{}
	v.SetFirstLiteral(first)
	v.SetLastLiteral(last)
	return v
}

// This type defines the structure and methods associated with a range.
type range_ struct {
	first Literal
	last  Literal
}

// This method returns the first literal for this range.
func (v *range_) GetFirstLiteral() Literal {
	return v.first
}

// This method sets the first literal for this range.
func (v *range_) SetFirstLiteral(first Literal) {
	if first == nil {
		panic("A range requires at least one literal.")
	}
	v.first = first
}

// This method returns the last literal for this range.
func (v *range_) GetLastLiteral() Literal {
	return v.last
}

// This method sets the last literal for this range.
func (v *range_) SetLastLiteral(last Literal) {
	v.last = last
}
