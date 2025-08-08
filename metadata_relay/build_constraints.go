//go:build !test
// +build !test

package main

// This file ensures that test dependencies are not included in production builds
// Test dependencies should only be imported in files with //go:build test
