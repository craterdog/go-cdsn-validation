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

// ALTERNATIVE IMPLEMENTATION

// This constructor creates a new alternative.
func Alternative(note Note, option OptionLike) AlternativeLike {
	var v = &alternative{}
	v.SetNote(note)
	v.SetOption(option)
	return v
}

// This type defines the structure and methods associated with an alternative.
type alternative struct {
	note   Note
	option OptionLike
}

// This method returns the note for this alternative.
func (v *alternative) GetNote() Note {
	return v.note
}

// This method sets the note for this alternative.
func (v *alternative) SetNote(note Note) {
	v.note = note
}

// This method returns the option for this alternative.
func (v *alternative) GetOption() OptionLike {
	return v.option
}

// This method sets the option for this alternative.
func (v *alternative) SetOption(option OptionLike) {
	if option == nil {
		panic("An alternative requires an option.")
	}
	v.option = option
}
