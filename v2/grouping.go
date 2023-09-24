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

type Grouping any

// This method attempts to parse a grouping. It returns the grouping and whether
// or not the grouping was successfully parsed.
func (v *parser) parseGrouping() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseExactlyN()
	if !ok {
		factor, token, ok = v.parseZeroOrOne()
	}
	if !ok {
		factor, token, ok = v.parseZeroOrMore()
	}
	if !ok {
		factor, token, ok = v.parseOneOrMore()
	}
	return factor, token, ok
}
