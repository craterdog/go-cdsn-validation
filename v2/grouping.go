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

// This constructor creates a new grouping.
func Grouping(definition DefinitionLike, type_ GroupingType, number Number) GroupingLike {
	var v = &grouping{}
	v.SetDefinition(definition)
	v.SetType(type_)
	v.SetNumber(number)
	return v
}

// This type defines the structure and methods associated with a grouping.
type grouping struct {
	definition DefinitionLike
	type_      GroupingType
	number     Number
}

// This method returns the definition for this grouping.
func (v *grouping) GetDefinition() DefinitionLike {
	return v.definition
}

// This method sets the definition for this grouping.
func (v *grouping) SetDefinition(definition DefinitionLike) {
	if definition == nil {
		panic("A grouping requires a definition.")
	}
	v.definition = definition
}

// This method returns the grouping type for this grouping.
func (v *grouping) GetType() GroupingType {
	return v.type_
}

// This method sets the grouping type for this grouping.
func (v *grouping) SetType(type_ GroupingType) {
	v.type_ = type_
}

// This method returns the number for this grouping.
func (v *grouping) GetNumber() Number {
	return v.number
}

// This method sets the number for this grouping.
func (v *grouping) SetNumber(number Number) {
	v.number = number
}
