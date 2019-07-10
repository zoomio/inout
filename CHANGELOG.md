# Changelog
All notable changes to this project will be documented in this file.

## 0.5.0
 - CSS query (`-q` option): concat all found texts, surround found texts with corresponding HTML tags, added tests for HTTP and File readers.

## 0.4.0
 - added `-q` option to allow query for contents of certain DOM elements via CSS selector;
 - constructors of the `Reader` are now returning reference instead of value.

## 0.3.1
 - fixed bug in `Reader#Read` leading to stack overflow error.

## 0.3.0
 - added `Reader#Read` method and updated `Reader#Close` method, to be compatible with `io.ReadCloser` interface.

## 0.2.3
 - binary release.

## 0.2.2
 - updated deployment configuration.

## 0.2.0
 - breaking change: renamed API methods `#ReadString` to `#ReadLine`, `#ReadAllStrings` to `#ReadWords` and `#LinesFromReader` to `#ReadLines`.

## 0.1.0
 - first release.