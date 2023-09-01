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

import (
	col "github.com/craterdog/go-collection-framework/v2"
)

// RULE IMPLEMENTATION

// This constructor creates a new rule.
func Rule(option OptionLike, alternatives col.Sequential[AlternativeLike]) RuleLike {
	var v = &rule{}
	v.SetOption(option)
	v.SetAlternatives(alternatives)
	return v
}

// This type defines the structure and methods associated with a rule.
type rule struct {
	option       OptionLike
	alternatives col.Sequential[AlternativeLike]
}

// This method returns the option for this rule.
func (v *rule) GetOption() OptionLike {
	return v.option
}

// This method sets the option for this rule.
func (v *rule) SetOption(option OptionLike) {
	v.option = option
}

// This method returns the alternatives for this rule.
func (v *rule) GetAlternatives() col.Sequential[AlternativeLike] {
	return v.alternatives
}

// This method sets the alternatives for this rule.
func (v *rule) SetAlternatives(alternatives col.Sequential[AlternativeLike]) {
	v.alternatives = alternatives
}
