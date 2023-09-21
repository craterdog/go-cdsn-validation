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
	var v = &rng{}
	v.SetFirstCharacter(first)
	v.SetLastCharacter(last)
	return v
}

// This type defines the structure and methods associated with a range.
type rng struct {
	first Character
	last  Character
}

// This method returns the first character for this range.
func (v *rng) GetFirstCharacter() Character {
	return v.first
}

// This method sets the first character for this range.
func (v *rng) SetFirstCharacter(first Character) {
	if len(first) == 0 {
		panic("A range requires at least one character.")
	}
	v.first = first
}

// This method returns the last character for this range.
func (v *rng) GetLastCharacter() Character {
	return v.last
}

// This method sets the last character for this range.
func (v *rng) SetLastCharacter(last Character) {
	v.last = last
}

// This method attempts to parse a range. It returns the range and whether or
// not the range was successfully parsed.
func (v *parser) parseRange() (RangeLike, *Token, bool) {
	var ok bool
	var token *Token
	var first Character
	var last Character
	var range_ RangeLike
	first, token, ok = v.parseCharacter()
	if !ok {
		// This is not a range.
		return range_, token, false
	}
	_, _, ok = v.parseDelimiter("..")
	if ok {
		last, token, ok = v.parseCharacter()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("CHARACTER",
				"$range",
				"$CHARACTER")
			panic(message)
		}
	}
	range_ = Range(first, last)
	return range_, token, true
}
