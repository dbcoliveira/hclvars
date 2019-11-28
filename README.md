# hclvars
Parses HCL variables files for further processing.


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
