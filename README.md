# hclvars
Parses HCL variables files for further processing.

## Details
HCLVars traverses a terraform variable files and store them in a simple array structure (hclvars.HCLVars).
The variable file should be in the format:
```
variable "v" {
  description = "description"
  type        = <type>
  default     = <default value>
}
```
The result hclvars.HCLvars handles with many different types of default and type attributes including maps and arrays.

## How to use

```
package main

import (
	"fmt"
	"github.com/dbcoliveira/hclvars"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		os.Exit(1)
	}
	hcl, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err.Error)
	}

	var vars hclvars.HCLVars
	pbytes, _ := vars.ParseHCLBytes(hcl)

	vars.ParseHCL(pbytes)

	for _, v := range vars.Variables {
		fmt.Println(v)
	}

}
```
