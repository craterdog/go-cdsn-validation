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
func Group(definition DefinitionLike, type_ GroupType, number Number) GroupLike {
	var v = &group{}
	v.SetDefinition(definition)
	v.SetType(type_)
	v.SetNumber(number)
	return v
}

// This type defines the structure and methods associated with a group.
type group struct {
	definition DefinitionLike
	type_      GroupType
	number     Number
}

// This method returns the definition for this group.
func (v *group) GetDefinition() DefinitionLike {
	return v.definition
}

// This method sets the definition for this group.
func (v *group) SetDefinition(definition DefinitionLike) {
	if definition == nil {
		panic("A group requires a definition.")
	}
	v.definition = definition
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
