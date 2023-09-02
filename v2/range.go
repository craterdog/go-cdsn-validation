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
func Range(first Character, last Character) RangeLike {
	var v = &range_{}
	v.SetFirstCharacter(first)
	v.SetLastCharacter(last)
	return v
}

// This type defines the structure and methods associated with a range.
type range_ struct {
	first Character
	last  Character
}

// This method returns the first character for this range.
func (v *range_) GetFirstCharacter() Character {
	return v.first
}

// This method sets the first character for this range.
func (v *range_) SetFirstCharacter(first Character) {
	if len(first) == 0 {
		panic("A range requires two characters.")
	}
	v.first = first
}

// This method returns the last character for this range.
func (v *range_) GetLastCharacter() Character {
	return v.last
}

// This method sets the last character for this range.
func (v *range_) SetLastCharacter(last Character) {
	if len(last) == 0 {
		panic("A range requires two characters.")
	}
	v.last = last
}
