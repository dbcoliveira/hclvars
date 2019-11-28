package hclvars

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	hclParser "github.com/hashicorp/hcl2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type variable struct {
	Name        string
	Description string
	VarType     []string
	Def         []string
}

type HCLVars struct {
	Variables []variable
}

func (h HCLVars) DecodeType(cv cty.Value) []string {
	var res []string
	cvType := cv.Type().GoString()

	switch cvType {
	case "cty.Bool":
		res = append(res, strconv.FormatBool(cv.True()))
	case "cty.Number":
		i, _ := cv.AsBigFloat().Int64()
		res = append(res, strconv.FormatInt(i, 10))
	case "cty.String":
		res = append(res, cv.AsString())

	case "cty.DynamicPseudoType":
		// something do be done here.
	}
	if strings.HasPrefix(cvType, "cty.Tuple") {
		for _, t := range cv.AsValueSlice() {
			res = append(res, h.DecodeType(t)...)
		}
	}
	if strings.HasPrefix(cvType, "cty.Object") {
		m := make(map[string]interface{})

		for k, v := range cv.AsValueMap() {
			m[k] = h.DecodeType(v)[0]
		}

		jsonStr, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, string(jsonStr))
	}
	return res
}

func (h *HCLVars) ParseHCL(hclfile *hcl.File) {
	// expected tfvars structure.
	var config struct {
		Variables []struct {
			Name   string                 `hcl:",label"`
			Labels map[string]interface{} `hcl:",remain"`
		} `hcl:"variable,block"`
	}

	gohcl.DecodeBody(hclfile.Body, nil, &config)

	for _, v := range config.Variables {
		var description string
		var def, ty []string

		// description label
		if v.Labels["description"] != nil {
			var cv cty.Value
			gohcl.DecodeExpression(v.Labels["description"].(*hcl.Attribute).Expr, nil, &cv)
			description = cv.AsString()
		}
		// type label
		if v.Labels["type"] != nil {
			vars := v.Labels["type"].(*hcl.Attribute).Expr.Variables()
			if len(vars) > 0 {
				for _, vt := range vars {
					varType := vt.RootName()
					ty = append(ty, varType)
				}
			}
		}
		// default label
		if v.Labels["default"] != nil {
			var cv cty.Value
			gohcl.DecodeExpression(v.Labels["default"].(*hcl.Attribute).Expr, nil, &cv)
			// decode dynamic types
			def = append(def, h.DecodeType(cv)...)
		}

		h.Variables = append(h.Variables, variable{
			Name:        v.Name,
			Description: description,
			Def:         def,
			VarType:     ty})
	}
}

// ParseHCLFile receives a file name from local filesytem as input.
func (h HCLVars) ParseHCLFile(filename string) (*hcl.File, hcl.Diagnostics) {
	return hclParser.NewParser().ParseHCLFile(filename)
}

func (h HCLVars) ParseHCLBytes(b []byte) (*hcl.File, hcl.Diagnostics) {
	// no filename provided
	return hclParser.NewParser().ParseHCL(b, "filename")
}
