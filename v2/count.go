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

import (
	col "github.com/craterdog/go-collection-framework/v2"
)

// COUNT IMPLEMENTATION

// This constructor creates a new count.
func Count(digits col.Sequential[Digit]) CountLike {
	var v = &count{}
	v.SetDigits(digits)
	return v
}

// This type defines the structure and methods associated with a count.
type count struct {
	digits col.Sequential[Digit]
}

// This method determines whether or not this count is the default (one).
func (v *count) IsDefault() bool {
	return v.digits.IsEmpty()
}

// This method returns the digits for this count.
func (v *count) GetDigits() col.Sequential[Digit] {
	return v.digits
}

// This method sets the digits for this count.
func (v *count) SetDigits(digits col.Sequential[Digit]) {
	if digits == nil {
		panic("A count cannot be nil.")
	}
	v.digits = digits
}
