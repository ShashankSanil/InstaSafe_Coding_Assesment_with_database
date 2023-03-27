package service

import (
	"fmt"
	"instasafe/common"
	"instasafe/repository"
	"strconv"
	"time"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	Rep *repository.Repository
}

func (s *Service) CreateEndUser(payload map[string]interface{}, language string) common.HTTPResponse {
	data, errcode, err := s.Rep.FindUserByEmail(payload["Email"].(string))
	if err != nil {
		return common.ErrorHandler(errcode, nil, 0, language)
	}
	fmt.Println(data)
	if len(data) > 0 {
		return common.ErrorHandler("107", nil, 0, language)
	}
	payload["city"] = ""
	payload["resetLocation"] = false
	res, errcode, err := s.Rep.CreateEndUser(payload)
	if err != nil {
		common.ErrorHandler(errcode, nil, 0, language)
	}
	return common.ErrorHandler(errcode, res, 1, language)
}

func (s *Service) CreateTransaction(payload map[string]interface{}, language string) common.HTTPResponse {
	var fixedDiff time.Duration = 60000000000
	currentTime := time.Now().UTC()
	fmt.Println("@@@@@@@@@@@@@@@  ", currentTime)
	timeStamp, err1 := time.Parse(time.RFC3339, payload["timestamp"].(string))
	if err1 != nil {
		return common.ErrorHandler("112", nil, 0, language)
	}
	diff := currentTime.Sub(timeStamp)
	if diff < 0 {
		return common.ErrorHandler("111", nil, 0, language)
	}
	if diff > fixedDiff {
		return common.ErrorHandler("110", nil, 0, language)
	}
	res, errcode, err := s.Rep.CreateTransaction(payload)
	if err != nil {
		return common.ErrorHandler(errcode, nil, 0, language)
	}
	return common.ErrorHandler(errcode, res, 1, language)
}

func (s *Service) GetStatistics(userid, city, language string) common.HTTPResponse {
	output := make(map[string]interface{})
	userDetailes, errcode, err := s.Rep.GetUserByID(userid)
	if err != nil {
		return common.ErrorHandler(errcode, userDetailes, 0, language)
	}
	var sum, min, max float64
	var count int
	if resetLocation, cok := userDetailes["resetLocation"].(bool); cok {
		var filter primitive.M
		if resetLocation {
			filter = bson.M{}
		} else {
			if userDetailes["city"].(string) != city {
				return common.ErrorHandler("123", nil, 0, language)
			}
			filter = bson.M{"city": city}
		}
		transactions, errcode, err1 := s.Rep.GetAllTransactions(filter)
		if err1 != nil {
			return common.ErrorHandler(errcode, userDetailes, 0, language)
		}
		var fixedDiff time.Duration = 60000000000
		currentTime := time.Now().UTC()
		for _, transaction := range transactions {
			timeStamp, err3 := time.Parse(time.RFC3339, transaction["timestamp"].(string))
			if err3 != nil {
				openlog.Error(err3.Error())
				return common.ErrorHandler("104", nil, 0, language)
			}
			diff := currentTime.Sub(timeStamp)
			if diff <= fixedDiff {
				amt, err2 := strconv.ParseFloat(transaction["amount"].(string), 32)
				if err2 != nil {
					openlog.Error(err2.Error())
					return common.ErrorHandler("122", nil, 0, language)
				}
				sum += amt
				count++
				if count == 1 {
					min = amt
				}
				if amt < min {
					min = amt
				}
				if amt > max {
					max = amt
				}
			}
		}
		output["Sum"] = sum
		if count != 0 {
			output["avg"] = sum / float64(count)
		} else {
			output["avg"] = 0
		}
		output["max"] = max
		output["min"] = min
		output["count"] = count
	}
	return common.ErrorHandler("124", output, 1, language)
}

func (s *Service) DeleteAllTransactions(language string) common.HTTPResponse {
	res, errcode, err := s.Rep.DeleteAllTransactions()
	if err != nil {
		return common.ErrorHandler(errcode, res, 0, language)
	}
	return common.ErrorHandler(errcode, res, 0, language)
}

func (s *Service) AddLoaction(uid string, payload map[string]interface{}, language string) common.HTTPResponse {
	payload["resetLocation"] = false
	res, errcode, err := s.Rep.UpdateLocation(uid,payload)
	if err != nil {
		return common.ErrorHandler(errcode, res, 0, language)
	}
	return common.ErrorHandler("116", res, 0, language)
}

func (s *Service) ResetLoaction(uid string, payload map[string]interface{}, language string) common.HTTPResponse {
	payload["resetLocation"] = true
	res, errcode, err := s.Rep.UpdateLocation(uid,payload)
	if err != nil {
		return common.ErrorHandler(errcode, res, 0, language)
	}
	return common.ErrorHandler("117", res, 0, language)
}
