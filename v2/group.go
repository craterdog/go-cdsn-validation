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

// GROUPING IMPLEMENTATION

// This constructor creates a new group.
func Group(expression ExpressionLike, type_ GroupType, number Number) GroupLike {
	var v = &group{}
	v.SetExpression(expression)
	v.SetType(type_)
	v.SetNumber(number)
	return v
}

// This type defines the structure and methods associated with a group.
type group struct {
	expression ExpressionLike
	type_      GroupType
	number     Number
}

// This method returns the expression for this group.
func (v *group) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this group.
func (v *group) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("A group requires an expression.")
	}
	v.expression = expression
}

// This method returns the group type for this group.
func (v *group) GetType() GroupType {
	return v.type_
}

// This method sets the group type for this group.
func (v *group) SetType(type_ GroupType) {
	v.type_ = type_
}

// This method returns the number for this group.
func (v *group) GetNumber() Number {
	return v.number
}

// This method sets the number for this group.
func (v *group) SetNumber(number Number) {
	v.number = number
}
