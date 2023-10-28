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

// CARDINALITY INTERFACE

// This interface defines the methods supported by all cardinality-like
// components.
type CardinalityLike interface {
	GetFirstNUMBER() NUMBER
	SetFirstNUMBER(first NUMBER)
	GetLastNUMBER() NUMBER
	SetLastNUMBER(last NUMBER)
}

// This constructor creates a new cardinality.
func Cardinality(first, last NUMBER) CardinalityLike {
	var v = &cardinality{}
	v.SetFirstNUMBER(first)
	v.SetLastNUMBER(last)
	return v
}

// CARDINALITY IMPLEMENTATION

// This type defines the structure and methods associated with a cardinality.
type cardinality struct {
	first NUMBER
	last  NUMBER
}

// This method returns the first number in the range for this cardinality.
func (v *cardinality) GetFirstNUMBER() NUMBER {
	return v.first
}

// This method sets the first number in the range for this cardinality.
func (v *cardinality) SetFirstNUMBER(first NUMBER) {
	if len(first) == 0 {
		panic("A cardinality requires that at least the first number be set.")
	}
	v.first = first
}

// This method returns the last number in the range for this cardinality.
func (v *cardinality) GetLastNUMBER() NUMBER {
	return v.last
}

// This method sets the last number in the range for this cardinality.
func (v *cardinality) SetLastNUMBER(last NUMBER) {
	v.last = last
}
