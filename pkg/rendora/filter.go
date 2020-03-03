/*
Copyright 2018 George Badawi.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rendora

import (
	"strings"
        "fmt"
	"github.com/gin-gonic/gin"
)

func isKeywordInSlice(slice []string, str string) bool {
	for _, s := range slice {
		if strings.Index(str, s) >= 0 {
			return true
		}
	}
	return false
}

func isInSlice(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func hasPrefixinSlice(slice []string, str string) bool {
	for _, s := range slice {
		if strings.HasPrefix(str, s) {
			return true
		}
	}
	return false
}

//isWhitelisted checks whether the current request is whitelisted (i.e. should be SSR'ed) or not
func (R *Rendora) isWhitelisted(c *gin.Context) bool {
	mua := c.Request.Header.Get("User-Agent")
	muaLower := strings.ToLower(mua)
	filters := &R.c.Filters
        
	fmt.Println("Passed in MUA: ", mua, muaLower)
	fmt.Println("UA Default: ", filters.UserAgent.Default)
	
	lenKeywords := len(filters.UserAgent.Exceptions.Keywords)
	lenExceptions := len(filters.UserAgent.Exceptions.Exact)
        fmt.Println("lenK: ", lenKeywords)
	fmt.Println("lenE: ", lenExceptions)

	switch filters.UserAgent.Default {
	case "whitelist":

		if lenKeywords > 0 && isKeywordInSlice(filters.UserAgent.Exceptions.Keywords, muaLower) {
		        fmt.Println("false: 1")
			return false
		}
		if lenExceptions > 0 && isInSlice(filters.UserAgent.Exceptions.Exact, mua) {
		        fmt.Println("false: 2")
			return false
		}
		break
	case "blacklist":
		if lenKeywords == 0 && lenExceptions == 0 {
		        fmt.Println("false: 3")
			return false
		}
		if lenKeywords > 0 && isKeywordInSlice(filters.UserAgent.Exceptions.Keywords, muaLower) == false {
		        fmt.Println("false: 4")
			return false
		}

		if lenExceptions > 0 && isInSlice(filters.UserAgent.Exceptions.Exact, mua) == false {
		        fmt.Println("false: 5")
			return false
		}

	}

	uri := c.Request.RequestURI

	switch filters.Paths.Default {
	case "blacklist":
		if len(filters.Paths.Exceptions.Exact) > 0 && isInSlice(filters.Paths.Exceptions.Exact, uri) {
			return true
		}

		if len(filters.Paths.Exceptions.Prefix) > 0 && hasPrefixinSlice(filters.Paths.Exceptions.Prefix, uri) {
			return true
		}
		return false
	case "whitelist":
		if len(filters.Paths.Exceptions.Exact) > 0 && isInSlice(filters.Paths.Exceptions.Exact, uri) {
		        fmt.Println("false: 6")
			return false
		}

		if len(filters.Paths.Exceptions.Prefix) > 0 && hasPrefixinSlice(filters.Paths.Exceptions.Prefix, uri) {
		        fmt.Println("false: 7")
			return false
		}
		return true
	default:
	        fmt.Println("false: default")
		return false
	}

}
