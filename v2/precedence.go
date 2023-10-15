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

// PRECEDENCE INTERFACE

// This interface defines the methods supported by all precedence-like
// components.
type PrecedenceLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
	GetRepetition() RepetitionLike
	SetRepetition(repetition RepetitionLike)
}

// This constructor creates a new precedence.
func Precedence(expression ExpressionLike, repetition RepetitionLike) PrecedenceLike {
	var v = &precedence{}
	v.SetExpression(expression)
	v.SetRepetition(repetition)
	return v
}

// PRECEDENCE IMPLEMENTATION

// This type defines the structure and methods associated with a precedence.
type precedence struct {
	expression ExpressionLike
	repetition RepetitionLike
}

// This method returns the expression for this precedence.
func (v *precedence) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this precedence.
func (v *precedence) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("A precedence requires an expression.")
	}
	v.expression = expression
}

// This method returns the repetition for this precedence.
func (v *precedence) GetRepetition() RepetitionLike {
	return v.repetition
}

// This method sets the repetition for this precedence.
func (v *precedence) SetRepetition(repetition RepetitionLike) {
	v.repetition = repetition
}
