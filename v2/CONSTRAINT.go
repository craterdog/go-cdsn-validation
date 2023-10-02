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
