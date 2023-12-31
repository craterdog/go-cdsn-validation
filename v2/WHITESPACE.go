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

type WHITESPACE string

const TokenWHITESPACE TokenType = "WHITESPACE"
const space = ` `
const whitespace = `[` + space + eol + `]+`

// This scanner is used for matching whitespace tokens.
var whitespaceScanner = reg.MustCompile(`^(?:` + whitespace + `)`)
