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

// RANGE INTERFACE

// This interface defines the methods supported by all range-like components.
type RangeLike interface {
	GetFirstCharacter() Character
	SetFirstCharacter(first Character)
	GetLastCharacter() Character
	SetLastCharacter(last Character)
}

// This constructor creates a new range.
func Range(first Character, last Character) RangeLike {
	var v = &range_{}
	v.SetFirstCharacter(first)
	v.SetLastCharacter(last)
	return v
}

// RANGE IMPLEMENTATION

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
		panic("A range requires at least one character.")
	}
	v.first = first
}

// This method returns the last character for this range.
func (v *range_) GetLastCharacter() Character {
	return v.last
}

// This method sets the last character for this range.
func (v *range_) SetLastCharacter(last Character) {
	v.last = last
}

// This private method appends a formatted range to the result.
func (v *formatter) formatRange(range_ RangeLike) {
	var first = range_.GetFirstCharacter()
	var last = range_.GetLastCharacter()
	v.formatCharacter(first)
	if len(last) > 0 {
		v.appendString("..")
		v.formatCharacter(last)
	}
}
