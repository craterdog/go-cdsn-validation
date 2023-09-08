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
func Range(first Rune, last Rune) RangeLike {
	var v = &rng{}
	v.SetFirstRune(first)
	v.SetLastRune(last)
	return v
}

// This type defines the structure and methods associated with a range.
type rng struct {
	first Rune
	last  Rune
}

// This method returns the first rune for this range.
func (v *rng) GetFirstRune() Rune {
	return v.first
}

// This method sets the first rune for this range.
func (v *rng) SetFirstRune(first Rune) {
	if len(first) == 0 {
		panic("A range requires two runes.")
	}
	v.first = first
}

// This method returns the last rune for this range.
func (v *rng) GetLastRune() Rune {
	return v.last
}

// This method sets the last rune for this range.
func (v *rng) SetLastRune(last Rune) {
	if len(last) == 0 {
		panic("A range requires two runes.")
	}
	v.last = last
}
