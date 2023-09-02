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
	fmt "fmt"
)

// FACTOR IMPLEMENTATION

// This function returns a string describing the type of the specified
// factor. This approach is used because the type switch CANNOT distinguish
// between abstract "Like" types if they support exactly the same method sets.
// The type switch CAN distinguish between the private structure types.
func GetType(factor Factor) string {
	switch value := factor.(type) {
	case Character:
		return "Character"
	case Intrinsic:
		return "Intrinsic"
	case Identifier:
		return "Identifier"
	case Literal:
		return "Literal"
	case *inversion:
		return "Inversion"
	case *grouping:
		return "Grouping"
	default:
		var message = fmt.Sprintf("An invalid factor type was found: %T", value)
		panic(message)
	}
}
