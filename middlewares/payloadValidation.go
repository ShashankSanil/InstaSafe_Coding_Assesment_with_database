package middlewares

import (
	"encoding/json"
	"fmt"
	"instasafe/common"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chassis/openlog"
	"github.com/xeipuuv/gojsonschema"
)

func PayloadValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := make(map[string]interface{})
		if err := c.BindJSON(&payload); err != nil {
			return
		}
		payload_bytes, _ := json.Marshal(payload)
		handlerName := c.HandlerName()
		arr := strings.Split(handlerName, ".")
		fmt.Println("@@@@@@@@@@@@@@@@", arr[2])
		schemaPath := "./../payloadSchemas/" + arr[2] + ".json"
		schema_bytes, _ := ioutil.ReadFile(schemaPath)
		schemaLoader := gojsonschema.NewBytesLoader(schema_bytes)
		documentLoader := gojsonschema.NewBytesLoader(payload_bytes)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			openlog.Error("error occured here" + err.Error())
			c.JSON(http.StatusBadRequest, common.ErrorHandler("105", result, 0, "en"))
			return
		}
		if result.Valid() {
			openlog.Info("Payload Validation completed")
			c.Set("Payload", payload)
			c.Next()
			return
		} else {
			validationErrors := make([]string, 0)
			for _, desc := range result.Errors() {
				validationErrors = append(validationErrors, desc.String())
			}
			fmt.Println(validationErrors, "######################################")
			c.JSON(http.StatusBadRequest, common.ErrorHandler("100", validationErrors, 0, "en"))
			openlog.Error("schema validation errors")
			return
		}

	}
}
