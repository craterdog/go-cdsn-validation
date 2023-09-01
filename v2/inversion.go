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

// INVERSION IMPLEMENTATION

// This constructor creates a new inversion.
func Inversion(factor Factor) InversionLike {
	var v = &inversion{}
	v.SetFactor(factor)
	return v
}

// This type defines the structure and methods associated with an inversion.
type inversion struct {
	factor Factor
}

// This method returns the factor for this inversion.
func (v *inversion) GetFactor() Factor {
	return v.factor
}

// This method sets the factor for this inversion.
func (v *inversion) SetFactor(factor Factor) {
	if factor == nil {
		panic("An inversion requires a factor.")
	}
	v.factor = factor
}
