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

// INVERSION INTERFACE

// This interface defines the methods supported by all inversion-like
// components.
type InversionLike interface {
	GetFactor() FactorLike
	SetFactor(factor FactorLike)
}

// This constructor creates a new inversion.
func Inversion(factor FactorLike) InversionLike {
	var v = &inversion{}
	v.SetFactor(factor)
	return v
}

// INVERSION IMPLEMENTATION

// This type defines the structure and methods associated with a inversion.
type inversion struct {
	factor FactorLike
}

// This method returns the factor for this inversion.
func (v *inversion) GetFactor() FactorLike {
	return v.factor
}

// This method sets the factor for this inversion.
func (v *inversion) SetFactor(factor FactorLike) {
	if factor == nil {
		panic("A inversion requires a factor.")
	}
	v.factor = factor
}
