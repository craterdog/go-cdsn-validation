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

// CONSTRAINT INTERFACE

// This interface defines the methods supported by all constraint-like
// components.
type ConstraintLike interface {
	GetLIMIT() LIMIT
	SetLIMIT(limit LIMIT)
	GetFactor() Factor
	SetFactor(factor Factor)
}

// This constructor creates a new constraint.
func Constraint(limit LIMIT, factor Factor) ConstraintLike {
	var v = &constraint{}
	v.SetLIMIT(limit)
	v.SetFactor(factor)
	return v
}

// CONSTRAINT IMPLEMENTATION

// This type defines the structure and methods associated with a constraint.
type constraint struct {
	limit  LIMIT
	factor Factor
}

// This method returns the number for this constraint.
func (v *constraint) GetLIMIT() LIMIT {
	return v.limit
}

// This method sets the number for this constraint.
func (v *constraint) SetLIMIT(limit LIMIT) {
	if len(limit) == 0 {
		panic("A constraint requires a limit.")
	}
	v.limit = limit
}

// This method returns the factor for this constraint.
func (v *constraint) GetFactor() Factor {
	return v.factor
}

// This method sets the factor for this constraint.
func (v *constraint) SetFactor(factor Factor) {
	if factor == nil {
		panic("A constraint requires a factor.")
	}
	v.factor = factor
}
