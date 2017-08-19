# Gastly - Simplistic Go Source File Rewriter

```
gastly <infile> <outfile> <package> [from=to] ...

Copies a Go source file with rewriting rules.

<infile> will be read.
<outfile> will be written.
<package> will be the new package name.

Each from=to specifies a rewrite rule, replacing any occurrences of "from" with
"to". These will be applied in order, so a leftmost replacement may affect what
the following replacements match.

If "to" begins with "droptype:" it will also drop any type specification that
matches "from" exactly. For example, "NumericType=droptype:int" will replace
"NumericType" with "int" everywhere and it will drop any "type NumericType ..."
specification it finds.
```

A real-world example can be found at
https://github.com/gholt/holdme/blob/master/internal/package.go

This was inspired by https://github.com/cheekybits/genny but I wanted something
simpler.

> Copyright See AUTHORS. All rights reserved.  
> Use of this source code is governed by a BSD-style  
> license that can be found in the LICENSE file.
