package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"instasafe/common"
	"instasafe/database"
	"instasafe/middlewares"
	"instasafe/repository"
	"instasafe/resource"
	"instasafe/service"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/go-chassis/openlog"
)

func getService(dbname string) *service.Service {
	Repo := repository.Repository{DBClient: database.GetClient(), DBName: dbname}
	return &service.Service{Rep: &Repo}
}

func duplicateInArray(arr []interface{}) []string {
	visited := make(map[string]bool)
	result := make(map[string]bool)
	for _, obj := range arr {
		errdetails := obj.(map[string]interface{})
		errcode := errdetails["errorcode"].(string)
		if visited[errcode] {
			result[errcode] = true
		} else {
			visited[errcode] = true
		}
	}
	res := []string{}
	for k, _ := range result {
		res = append(res, k)
	}
	return res
}

func CheckandLoadErrors() error {
	errbytes, err := ioutil.ReadFile("./../server/conf/errorcode.json")
	if err != nil {
		openlog.Error(err.Error())
		return err
	}
	errorslist := make([]interface{}, 0)

	err = json.Unmarshal(errbytes, &errorslist)
	if err != nil {
		openlog.Error(err.Error())
		return err
	}

	dups := duplicateInArray(errorslist)
	if len(dups) > 0 {
		fmt.Println("Duplicates exists in Errorcode: ", dups)
		return errors.New(" Duplicates exists in Errorcode")
	}

	for _, errs := range errorslist {
		ERROR := errs.(map[string]interface{})
		errcode := ERROR["errorcode"].(string)
		common.Errorcodes[errcode] = ERROR
	}

	return nil
}

func main() {
	resource := resource.Resource{ServiceProvider: getService}

	common.Dbname = "InstaSafeCollection"
	if err := database.Connect(); err != nil {
		openlog.Fatal("Error occured while connecting to database")
		return
	}

	err := CheckandLoadErrors()
	if err != nil {
		openlog.Fatal(err.Error())
		return
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(middlewares.CORSMiddleware())
	resource.URLRoutes(router)
	router.Run("0.0.0.0:5055")

}
