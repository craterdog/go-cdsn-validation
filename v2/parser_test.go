/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn_test

import (
	fmt "fmt"
	cds "github.com/craterdog/go-cdsn-validation/v2"
	ass "github.com/stretchr/testify/assert"
	osx "os"
	sts "strings"
	tes "testing"
)

const testDirectory = "./grammars/"

func TestParsingRoundtrips(t *tes.T) {

	var files, err = osx.ReadDir(testDirectory)
	if err != nil {
		panic("Could not find the " + testDirectory + " directory.")
	}

	for _, file := range files {
		var filename = testDirectory + file.Name()
		if sts.HasSuffix(filename, ".cdsn") {
			fmt.Println(filename)
			var expected, _ = osx.ReadFile(filename)
			var grammar = cds.ParseDocument(expected)
			var document = cds.FormatGrammar(grammar)
			ass.Equal(t, string(expected), string(document))
		}
	}
}
