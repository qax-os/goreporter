##Get info on length of functions in a Go package.

Given package is searched in  directories provided by envs in following order: GOPATH, GOROOT. AST is generated only for Go source files present in package path, ie, `flen crypto` shall only parse `crypto.go`. For parsing `crypto/sha1`, full package path needs to be provided, ie `flen crypto/sha1`. For externally implemented functions, line number will be 0, as their is no enough information available to get their line numbers.

###Install
`go get github.com/lafolle/flen/cmd/flen`

###Usage
```
Usage: flen [options] <pkg>
  -bs int
        bucket size (natural number) (default 5)
  -l int
        min length (inclusive)
  -t    include tests files
  -u int
        max length (exclusive) (default 1000000)	
```
###Examples
Simple usage  
```
$ flen strings
Full path of pkg:  /usr/local/go/src/strings
Externally implemented funcs
+-------+-----------+-------------------------------------------+---------+------+
| INDEX |   NAME    |                 FILEPATH                  | LINE NO | SIZE |
+-------+-----------+-------------------------------------------+---------+------+
|     0 | IndexByte | /usr/local/go/src/strings/strings_decl.go |       0 |    0 |
+-------+-----------+-------------------------------------------+---------+------+

[1-6)   -       ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
[6-11)  -       ∎∎∎∎∎∎∎∎∎∎
[11-16) -       ∎∎∎∎∎∎∎∎∎∎∎∎
[16-21) -       ∎∎∎∎∎
[21-26) -       ∎∎∎
[26-31) -       ∎∎∎∎
[31-36) -       ∎∎∎∎
[36-41) -       ∎∎∎
[41-46) -       ∎∎
[46-51) -       ∎
[51-56) -
[56-61) -       ∎
[61-66) -
$
```  

###Range of lengths
To get all function/methods whose lengths fall in given range:   
```
flen -l 16 -u 41 strings
Full path of pkg:  /usr/local/go/src/strings
Functions with length in range [16, 41)
+-------+---------------------+--------------------------------------+---------+------+
| INDEX |        NAME         |               FILEPATH               | LINE NO | SIZE |
+-------+---------------------+--------------------------------------+---------+------+
|     0 | Replace             | /usr/local/go/src/strings/replace.go |     365 |   17 |
|     1 | Seek                | /usr/local/go/src/strings/reader.go  |     109 |   17 |
|     2 | Join                | /usr/local/go/src/strings/strings.go |     388 |   18 |
|     3 | WriteString         | /usr/local/go/src/strings/replace.go |     432 |   20 |
|     4 | isSeparator         | /usr/local/go/src/strings/strings.go |     502 |   20 |
|     5 | WriteString         | /usr/local/go/src/strings/replace.go |     385 |   22 |
|     6 | genSplit            | /usr/local/go/src/strings/strings.go |     281 |   23 |
|     7 | explode             | /usr/local/go/src/strings/strings.go |      17 |   25 |
|     8 | Replace             | /usr/local/go/src/strings/replace.go |     461 |   26 |
|     9 | WriteString         | /usr/local/go/src/strings/replace.go |     490 |   27 |
|    10 | makeGenericReplacer | /usr/local/go/src/strings/replace.go |     240 |   29 |
|    11 | FieldsFunc          | /usr/local/go/src/strings/strings.go |     353 |   30 |
|    12 | Replace             | /usr/local/go/src/strings/strings.go |     681 |   31 |
|    13 | Index               | /usr/local/go/src/strings/strings.go |     147 |   33 |
|    14 | LastIndex           | /usr/local/go/src/strings/strings.go |     184 |   33 |
|    15 | lookup              | /usr/local/go/src/strings/replace.go |     193 |   33 |
|    16 | Map                 | /usr/local/go/src/strings/strings.go |     422 |   37 |
|    17 | WriteString         | /usr/local/go/src/strings/replace.go |     312 |   38 |
|    18 | makeStringFinder    | /usr/local/go/src/strings/search.go  |      48 |   40 |
+-------+---------------------+--------------------------------------+---------+------+

[1-6)   -       ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
[6-11)  -       ∎∎∎∎∎∎∎∎∎∎
[11-16) -       ∎∎∎∎∎∎∎∎∎∎∎∎
[16-21) -       ∎∎∎∎∎
[21-26) -       ∎∎∎
[26-31) -       ∎∎∎∎
[31-36) -       ∎∎∎∎
[36-41) -       ∎∎∎
[41-46) -       ∎∎
[46-51) -       ∎
[51-56) -
[56-61) -       ∎
[61-66) -
```

###Including test files
By default, test files are ignored. Enable parsing test files by `-t` flag
```
$ flen -t strings
Full path of pkg:  /usr/local/go/src/strings

Externally implemented funcs
+-------+-----------+-------------------------------------------+---------+------+
| INDEX |   NAME    |                 FILEPATH                  | LINE NO | SIZE |
+-------+-----------+-------------------------------------------+---------+------+
|     0 | IndexByte | /usr/local/go/src/strings/strings_decl.go |       0 |    0 |
+-------+-----------+-------------------------------------------+---------+------+

[1-6)		-	∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
[6-11)		-	∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
[11-16)		-	∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
[16-21)		-	∎∎∎∎∎∎∎∎∎∎∎∎
[21-26)		-	∎∎∎∎∎
[26-31)		-	∎∎∎∎∎∎
[31-36)		-	∎∎∎∎
[36-41)		-	∎∎∎∎
[41-46)		-	∎∎∎∎∎∎∎
[46-51)		-	∎
[51-56)		-	∎
[56-61)		-	∎
[61-66)		-	
[66-71)		-	∎
[71-76)		-	
[76-81)		-	
[81-86)		-	
[86-91)		-	
[91-96)		-	
[96-101)	-	
[101-106)	-	
[106-111)	-	
[111-116)	-	
[116-121)	-	
[121-126)	-	
[126-131)	-	
[131-136)	-	
[136-141)	-	
[141-146)	-	
[146-151)	-	
[151-156)	-	
[156-161)	-	
[161-166)	-	
[166-171)	-	
[171-176)	-	
[176-181)	-	
[181-186)	-	
[186-191)	-	
[191-196)	-	
[196-201)	-	
[201-206)	-	
[206-211)	-	
[211-216)	-	
[216-221)	-	
[221-226)	-	
[226-231)	-	
[231-236)	-	
[236-241)	-	
[241-246)	-	
[246-251)	-	
[251-256)	-	
[256-261)	-	
[261-266)	-	∎
[266-271)	-	
```
