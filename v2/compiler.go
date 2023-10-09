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
	byt "bytes"
	col "github.com/craterdog/go-collection-framework/v2"
	osx "os"
	sts "strings"
	uni "unicode"
)

// COMPILER INTERFACE

// This function compiles the specified document into its corresponding parser.
func CompileDocument(directory, packageName string, document DocumentLike) {
	var v = &compiler{directory: directory, packageName: packageName}
	v.compileDocument(document)
}

// COMPILER IMPLEMENTATION

// This private function determines whether or not the specified name is a token
// name.
func isTokenName(name NAME) bool {
	return uni.IsUpper(rune(name[1]))
}

// This private function replaces all occurences of the target string with the
// specified name.
func replaceName(template string, target string, name string) string {
	var nameLower, nameUpper string
	var targetLower = "#" + target + "#"
	var targetUpper = "#" + sts.ToUpper(target[0:1]) + target[1:] + "#"
	if isTokenName(NAME(name)) {
		nameLower = sts.ToLower(name)
		nameUpper = name
	} else {
		nameLower = name
		nameUpper = sts.ToUpper(name[0:1]) + name[1:]
	}
	template = sts.ReplaceAll(template, targetLower, nameLower)
	template = sts.ReplaceAll(template, targetUpper, nameUpper)
	return template
}

// This type defines the structure and methods for a compiler agent.
type compiler struct {
	directory     string
	packageName   string
	scannerBuffer byt.Buffer
	parserBuffer  byt.Buffer
}

// This private method creates a new configuration (package.go) file if one
// does not already exist.
func (v *compiler) initializeConfiguration() {
	var err error
	var template []byte
	var configuration = v.directory + "/package.go"
	_, err = osx.Open(configuration)
	if err != nil {
		template, err = osx.ReadFile("./templates/package.tp")
		if err != nil {
			panic(err)
		}
		template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
		err = osx.WriteFile(configuration, template, 0666)
		if err != nil {
			panic(err)
		}
	}
}

// This private method creates the byte buffer for the generated scanner code.
func (v *compiler) initializeScanner() {
	var template, err = osx.ReadFile("./templates/scanner.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.scannerBuffer.Write(template[0 : len(template)-1])
}

// This private method creates the byte buffer for the generated parser code.
func (v *compiler) initializeParser() {
	var template, err = osx.ReadFile("./templates/parser.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.parserBuffer.Write(template[0 : len(template)-1])
}

// This private method appends the scan token template for the specified name to
// the scanner byte buffer.
func (v *compiler) appendScanToken(name NAME) {
	const template = `

// This method adds a new #token# token with the current scanner information
// to the token channel. It returns true if a new #token# token was found.
func (v *scanner) scan#Token#() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(#token#Scanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token#Token#)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}`
	var snippet = replaceName(template, "token", string(name))
	v.scannerBuffer.WriteString(snippet)
}

// This private method appends the parse token template for the specified name
// to the parser byte buffer.
func (v *compiler) appendParseToken(name NAME) {
	const template = `

// This method attempts to parse a new #token# token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parse#Token#() (#Token#, *Token, bool) {
	var #token#_ #Token#
	var token = v.nextToken()
	if token.Type != Token#Token# {
		v.backupOne(token)
		return #token#_, token, false
	}
	#token#_ = #Token#(token.Value)
	return #token#_, token, true
}`
	var snippet = replaceName(template, "token", string(name))
	v.parserBuffer.WriteString(snippet)
}

// This private method appends the parse rule start template for the specified
// name to the parser byte buffer.
func (v *compiler) appendParseRuleStart(name NAME) {
	const template = `

// This method attempts to parse a new #rule#. It returns the #rule#
// and whether or not the #rule# was successfully parsed.
func (v *parser) parse#Rule#() (#Rule#Like, *Token, bool) {
	var ok bool
	var token *Token
	var #rule#_ #Rule#Like`
	var snippet = replaceName(template, "rule", string(name))
	v.parserBuffer.WriteString(snippet)
}

// This private method appends the parse rule end template for the specified
// name to the parser byte buffer.
func (v *compiler) appendParseRuleEnd(name NAME) {
	const template = `
	#rule#_ = #Rule#(#arguments#)
	return #rule#_, token, true
}`
	var snippet = replaceName(template, "rule", string(name))
	v.parserBuffer.WriteString(snippet)
}

// This private method writes the byte buffer for the generated scanner code into
// a file.
func (v *compiler) finalizeScanner() {
	var filename = v.directory + "scanner.go"
	var err = osx.WriteFile(filename, v.scannerBuffer.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

// This private method writes the byte buffer for the generated parser code into
// a file.
func (v *compiler) finalizeParser() {
	var filename = v.directory + "parser.go"
	var err = osx.WriteFile(filename, v.parserBuffer.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

// This private method compiles the specified definition.
func (v *compiler) compileDefinition(definition DefinitionLike) {
	var symbol = definition.GetSYMBOL()
	var name = symbol.GetNAME()
	switch {
	case string(name) == "INTRINSIC":
		// Intrinsics are automatically part of every parser.
	case isTokenName(name):
		v.appendScanToken(name)
		v.appendParseToken(name)
	default:
		v.appendParseRuleStart(name)
	}
	if !isTokenName(name) {
		v.appendParseRuleEnd(name)
	}
}

// This private method compiles the specified document.
func (v *compiler) compileDocument(document DocumentLike) {
	v.initializeConfiguration()
	v.initializeScanner()
	v.initializeParser()
	var statements = document.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.compileStatement(statement)
	}
	v.finalizeScanner()
	v.finalizeParser()
}

// This private method compiles the specified statement.
func (v *compiler) compileStatement(statement Statement) {
	switch actual := statement.(type) {
	case *definition:
		v.compileDefinition(actual)
	case COMMENT:
	}
}
