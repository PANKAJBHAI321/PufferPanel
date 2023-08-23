/*
 Copyright 2019 Padduck, LLC
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

package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func NotImplemented(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotImplemented, pufferpanel.ErrNotImplemented)
}

func CreateOptions(options ...string) gin.HandlerFunc {
	replacement := make([]string, len(options)+1)

	copy(replacement, options)

	replacement[len(options)] = http.MethodOptions
	res := strings.Join(replacement, ",")

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", res)
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", res)
		c.Header("Content-Type", binding.MIMEJSON)
		c.AbortWithStatus(http.StatusOK)
	}
}

func HandleError(c *gin.Context, err error, statusCode int) bool {
	if err != nil {
		logging.Error.Printf("%s", err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatusJSON(statusCode, &pufferpanel.ErrorResponse{Error: pufferpanel.FromError(err)})
		}

		return true
	}

	return false
}
