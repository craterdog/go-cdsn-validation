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

// CLASS NAMESPACE

// Private Class Namespace Type

type constraintClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var constraintClass = &constraintClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func ConstraintClass() ConstraintClassLike {
	return constraintClass
}

// Public Class Constructors

func (c *constraintClass_) FromNumber(number string) ConstraintLike {
	var constraint = &constraint_{
		// This class does not initialize any attributes.
	}
	constraint.SetFirst(number)
	constraint.SetLast(number)
	return constraint
}

func (c *constraintClass_) FromRange(first, last string) ConstraintLike {
	var constraint = &constraint_{
		// This class does not initialize any attributes.
	}
	constraint.SetFirst(first)
	constraint.SetLast(last)
	return constraint
}

// CLASS INSTANCES

// Private Class Type Definition

type constraint_ struct {
	first string
	last  string
}

// Public Interface

func (v *constraint_) GetFirst() string {
	return v.first
}

func (v *constraint_) GetLast() string {
	return v.last
}

func (v *constraint_) SetFirst(first string) {
	if len(first) < 1 {
		panic("A constraint requires a first number.")
	}
	v.first = first
}

func (v *constraint_) SetLast(last string) {
	v.last = last
}
