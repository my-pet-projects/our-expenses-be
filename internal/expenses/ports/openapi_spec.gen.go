// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
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

	"H4sIAAAAAAAC/8xYTXPbNhD9Kxy0R9Zyk5tuaex63EM846bTQ8YHiFxJyJAAvVgq0Xj43zsL8AMSIZFu",
	"5JmcTBOL/Xhv8bDUi8hMWRkNmqxYvgibbaGU7vGjJNgY3PNzhaYCJAVuRWVG81/aVyCWwhIqvRFNKlTO",
	"r3OwGaqKFFuJf7R6riFReWLWCW0hyTq/qVgbLCWJpahrlYt07LCAHRRBKKUJNoC8pGUJ42ifZAmRQCPH",
	"lcSuZEVQuodfEdZiKX5ZDJAsWjwWPRhN70wiyr1omlQgPNcKIRfLL8LV4XLjILQVXRFP/Uaz+goZsafO",
	"6+33CrT14B5CnQUkzM0PAm+zimvDj2tLxQalzj8bksWUk7vBskmFrVdtSgpej3KPxxTaAcVBpjGobyRB",
	"5/4RKoN0Guzb1wI4nXcqckmuX/ue5xe/kXKdMupP+J5tpd7AoySYwV9o/L9ZO8LWJZyOQTlwf5xpDPpb",
	"RIMRtE0eOcDOOHFrgT4oTe/fDTgFMlCCtXJz0lG3PIL4uJF8wM48WkZQ6LialbTwsUYEne2j2vg6/rG1",
	"Hi2QxA3QmUhxFg/SG3lp401VHdOny5Y8/8g5FqbkoWvhwzpPdGkrgssXIYviYS2WX85n8Am+DcJ5Cd2e",
	"dXu24j59eY7vpafmaajzlAYyZH5tPhcRbb3YVRJhtEtvUvPvDiIeq4/eARLkU8n0Vxrxw3xQ3L57vTaT",
	"Xdo6ToOcYtXcawLc+VpA16WHg5u6NNoNGnuQGGwdjlbQqicvvftI93WkJvc3r+4/LqcsQVNUGi6oGxWq",
	"7Mjc1KsisNV1ufK3xXMtNSnazzSnnsNXkY2qmhbmAPeuhiC/A9FqZWzIJtYeA8ePYCsT5TqmMAO1Gr4V",
	"+0TmOeQ/JDOR5OJ35tkm6GmdgHGAye+IhT+lAUH8Q1DchiTwPUrP8gGMb7N1mcjS1Jom4WIvAdUnc+96",
	"8Ec1zKDaKD2twP2Gbg6ZO36OKuwjjmtjFCGrUdH+b3bVDlIgEfBDTdvhvz+7Hvzr388i9R+p7MmvDihv",
	"iSrRsGPVAnbEz8PNA1srKtj8ocakG2sTC7hzvnaA1pv/fnV9de1gq0DLSomleO9e+U87l+4i/NqqjKWI",
	"iiLw1JFIPmHdyUqUdofO7i1BeSVcEJS8hZVYfMjz2/4MMqBg6Q+T71veqRVWWVWFyty2xVfrP8s9L1Os",
	"heMLI3Y0PbdZkmFB4D9DsiJkmLAGR7nXHAfDu+vrN0izl7Uz6WJvwxZrWRd0sUz8Z0wkeK2Z04xYNlsb",
	"pw+l5PkvTr+zWeAwaG0g0jl3oMFNrH3X+B3jduks2/HL/cAhSyBA6wbZQ79rNGXCtwoTu1YFASYrljnF",
	"q881uC9q/wuLMx5RngagzbmmeUI+zIHM7AzIvEF8BFsXZBPVjVXx2MHy6QzOtU0/tjXN0xuek8Ox/uwR",
	"YYuf8qTE27078V6gfTvXWLRyb5eLxcvWWGK+mgWrdCp2EpVcFR7lbtEfr7ZSUZhMFrzEzp+a/wIAAP//",
	"kagtzAgVAAA=",
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
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
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
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
