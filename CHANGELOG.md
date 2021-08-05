# Changelog
All notable changes to this project will be documented in this file.

## 0.11.0
 - bumped `github.com/chromedp/chromedp`.

## 0.10.0
 - improved `-q` option, it is now acts as `document.querySelectorAll()` so it will return all the matching nodes.

## 0.9.1
 - added retries with backoff to HTML calls;
 - fixed `-q` option, which allows querying with CSS selectors.

## 0.8.2
 - small fix: changed size of the read buffer from `102481024` to `1024*1024`.

## 0.8.0
 - better output to STDIN in cli, line by line instead of printing slice of lines.

## 0.7.0
 - increased buffer size to 1MB;
 - introduced option for setting timeout `-t`, defaults to 5 seconds;
 - introduced option for enabling verbose mode `-v`.

## 0.6.0
 - reverted back `#New` to have only single argument - `source` for backwards compatibility;
 - re-used "self-referential functions and the design of options" approach by Rob Pike by introducing `Option` and new constructor `#NewInOut` which uses it.

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