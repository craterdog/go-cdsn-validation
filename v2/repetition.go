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

// REPETITION INTERFACE

// This interface defines the methods supported by all repetition-like
// components.
type RepetitionLike interface {
	GetCONSTRAINT() CONSTRAINT
	SetCONSTRAINT(constraint CONSTRAINT)
	GetFirstNUMBER() NUMBER
	SetFirstNUMBER(first NUMBER)
	GetLastNUMBER() NUMBER
	SetLastNUMBER(last NUMBER)
}

// This constructor creates a new repetition.
func Repetition(constraint CONSTRAINT, first, last NUMBER) RepetitionLike {
	if len(constraint) == 0 && len(first) == 0 && len(last) == 0 {
		panic("A repetition requires at least one of its attributes to be set.")
	}
	var v = &repetition{}
	v.SetCONSTRAINT(constraint)
	v.SetFirstNUMBER(first)
	v.SetLastNUMBER(last)
	return v
}

// REPETITION IMPLEMENTATION

// This type defines the structure and methods associated with a repetition.
type repetition struct {
	constraint CONSTRAINT
	first      NUMBER
	last       NUMBER
}

// This method returns the constraint for this repetition.
func (v *repetition) GetCONSTRAINT() CONSTRAINT {
	return v.constraint
}

// This method sets the constraint for this repetition.
func (v *repetition) SetCONSTRAINT(constraint CONSTRAINT) {
	if len(constraint) > 0 {
		v.constraint = constraint
		v.first = ""
		v.last = ""
	}
	v.constraint = constraint
}

// This method returns the first number in the range for this repetition.
func (v *repetition) GetFirstNUMBER() NUMBER {
	return v.first
}

// This method sets the first number in the range for this repetition.
func (v *repetition) SetFirstNUMBER(first NUMBER) {
	if len(first) > 0 {
		v.constraint = ""
		v.first = first
		v.last = ""
	}
}

// This method returns the last number in the range for this repetition.
func (v *repetition) GetLastNUMBER() NUMBER {
	return v.last
}

// This method sets the last number in the range for this repetition.
func (v *repetition) SetLastNUMBER(last NUMBER) {
	if len(last) > 0 {
		v.constraint = ""
		if len(v.first) == 0 {
			panic("A repetition requires that the first number be set if the second number is set.")
		}
		v.last = last
	}
}
