/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn

import (
	reg "regexp"
)

type CONSTRAINT string

const TokenCONSTRAINT TokenType = "CONSTRAINT"
const (
	number     = digit + `+`
	constraint = `[~?*+]|` + number
)

// This scanner is used for matching constraint tokens.
var constraintScanner = reg.MustCompile(`^(?:` + constraint + `)`)

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
	GetCONSTRAINT() CONSTRAINT
	SetCONSTRAINT(constraint CONSTRAINT)
	GetFactor() FactorLike
	SetFactor(factor FactorLike)
}

// This constructor creates a new constraint.
func Constraint(constraint CONSTRAINT, factor FactorLike) ConstraintLike {
	var v = &constraint{}
	v.SetCONSTRAINT(constraint)
	v.SetFactor(factor)
	return v
}

// CONSTRAINT IMPLEMENTATION

// This type defines the structure and methods associated with a constraint.
type constraint struct {
	constraint CONSTRAINT
	factor     FactorLike
}

// This method returns the number for this constraint.
func (v *constraint) GetCONSTRAINT() CONSTRAINT {
	return v.constraint
}

// This method sets the number for this constraint.
func (v *constraint) SetCONSTRAINT(constraint CONSTRAINT) {
	if len(constraint) == 0 {
		panic("A constraint requires a constraint.")
	}
	v.constraint = constraint
}

// This method returns the factor for this constraint.
func (v *constraint) GetFactor() FactorLike {
	return v.factor
}

// This method sets the factor for this constraint.
func (v *constraint) SetFactor(factor FactorLike) {
	if factor == nil {
		panic("A constraint requires a factor.")
	}
	v.factor = factor
}
