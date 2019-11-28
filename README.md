# hclvars
Parses HCL variables files for further processing.

## Details
HCLVars traverses terraform variable files and store them in a simple array structure (hclvars.HCLVars). This is useful to parse dynamic random terraform variable.tf files part of a module.  

The variable file should be in the format:
```
variable "v" {
  description = <description>
  type        = <type>
  default     = <default value>
}
```
The result hclvars.HCLvars handles with many different types of default and type attributes including maps and arrays.

## Variable files source examples

https://raw.githubusercontent.com/GoogleCloudPlatform/terraform-google-nat-gateway/master/variables.tf
https://raw.githubusercontent.com/hashicorp/terraform-azurerm-consul/master/vars.tf

Examples takaen from https://registry.terraform.io/browse/modules

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
