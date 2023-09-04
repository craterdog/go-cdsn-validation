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
func Grouping(rule RuleLike, type_ GroupingType, count CountLike) GroupingLike {
	var v = &grouping{}
	v.SetRule(rule)
	v.SetType(type_)
	v.SetCount(count)
	return v
}

// This type defines the structure and methods associated with a grouping.
type grouping struct {
	rule  RuleLike
	type_ GroupingType
	count CountLike
}

// This method returns the rule for this grouping.
func (v *grouping) GetRule() RuleLike {
	return v.rule
}

// This method sets the rule for this grouping.
func (v *grouping) SetRule(rule RuleLike) {
	if rule == nil {
		panic("A grouping requires a rule.")
	}
	v.rule = rule
}

// This method returns the grouping type for this grouping.
func (v *grouping) GetType() GroupingType {
	return v.type_
}

// This method sets the grouping type for this grouping.
func (v *grouping) SetType(type_ GroupingType) {
	v.type_ = type_
}

// This method returns the count for this grouping.
func (v *grouping) GetCount() CountLike {
	return v.count
}

// This method sets the count for this grouping.
func (v *grouping) SetCount(count CountLike) {
	v.count = count
}
