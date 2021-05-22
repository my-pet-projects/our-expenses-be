// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package ports

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xUzW4TMRB+lWjguOoGuPmGSkFwCRKVOFQ9TLyTxNWubcbjqFG0747G3s0WBYkeuMSO",
	"5+f7mUnOsMVE31EOYKDF6KABG4YYPHlJYM6Q7IEGLNdbFNoHPuk9cojE4qhEXKefu8ADChjI2XXQgJwi",
	"gYEk7PwexgZ6OlKvmVPEeaE9sYY8DvQistREZPLy9XUANblSEhrK5S3TDgy8aRdd7SSqvSgaL82QGU+1",
	"l5pyxWhsgOlXdkwdmAcoPAr5F1Sn4lnwwuvxAhO2T2RFce6YA19bSvPzlcgkKDn9m9qU10ytrrG1F9nM",
	"Tk4/1JAKvCVk4o+5qq/fPs/Gf/t5D03dCe1Uo8skDiIRRm3s/C5ofUfJsoviggcD95tPG8120mv6JvPq",
	"7jmST5RWifhYeh2JU01/d7O+WavkEMnrchr4UJ6qwYVua+sMJ9v2JHqok6igujjwheR2yVKXUgwKqqnv",
	"12s9bPBCvhRjjL2zpbx9Sspk/hm8fp/Ugz+1S+gClNcd5l7+G2bdn78AZk/PkaxQt6I5p85cjU5gHs6Q",
	"uZ/Glkzbng8hiS7zOP0VHJEdbvtq1BysY51EQB8s9hrS9o/j7wAAAP//gvhgaVEEAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.Swagger, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.SwaggerLoader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadSwaggerFromData(specData)
	if err != nil {
		return
	}
	return
}
