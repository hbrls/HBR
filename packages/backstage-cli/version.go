package main

// cliVersion 程序版本号
// 在构建时通过 ldflags 注入，例如：-ldflags "-X main.cliVersion=0.0.1"
var cliVersion = "dev"
