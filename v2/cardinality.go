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
	GetCONSTRAINT() CONSTRAINT
	SetCONSTRAINT(constraint CONSTRAINT)
	GetFirstNUMBER() NUMBER
	SetFirstNUMBER(first NUMBER)
	GetLastNUMBER() NUMBER
	SetLastNUMBER(last NUMBER)
}

// This constructor creates a new cardinality.
func Cardinality(constraint CONSTRAINT, first, last NUMBER) CardinalityLike {
	if len(constraint) == 0 && len(first) == 0 && len(last) == 0 {
		panic("A cardinality requires at least one of its attributes to be set.")
	}
	var v = &cardinality{}
	v.SetCONSTRAINT(constraint)
	v.SetFirstNUMBER(first)
	v.SetLastNUMBER(last)
	return v
}

// CARDINALITY IMPLEMENTATION

// This type defines the structure and methods associated with a cardinality.
type cardinality struct {
	constraint CONSTRAINT
	first      NUMBER
	last       NUMBER
}

// This method returns the constraint for this cardinality.
func (v *cardinality) GetCONSTRAINT() CONSTRAINT {
	return v.constraint
}

// This method sets the constraint for this cardinality.
func (v *cardinality) SetCONSTRAINT(constraint CONSTRAINT) {
	if len(constraint) > 0 {
		v.constraint = constraint
		v.first = ""
		v.last = ""
	}
	v.constraint = constraint
}

// This method returns the first number in the range for this cardinality.
func (v *cardinality) GetFirstNUMBER() NUMBER {
	return v.first
}

// This method sets the first number in the range for this cardinality.
func (v *cardinality) SetFirstNUMBER(first NUMBER) {
	if len(first) > 0 {
		v.constraint = ""
		v.first = first
		v.last = ""
	}
}

// This method returns the last number in the range for this cardinality.
func (v *cardinality) GetLastNUMBER() NUMBER {
	return v.last
}

// This method sets the last number in the range for this cardinality.
func (v *cardinality) SetLastNUMBER(last NUMBER) {
	if len(last) > 0 {
		v.constraint = ""
		if len(v.first) == 0 {
			panic("A cardinality requires that the first number be set if the second number is set.")
		}
		v.last = last
	}
}
