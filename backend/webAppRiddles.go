package main

import (
	"net/http"

	"github.com/goccy/go-json"
)

type webAppRiddles struct {
	db    *database
	steps int
}

func (me *webAppRiddles) init(db *database) []namedWebFunction {
	me.db = db
	me.steps = 4
	return []namedWebFunction{
		{"/api/riddles/primeNumbers", me.getPrimeNumbers},
		{"/api/riddles/new", me.createRiddle},
		{"/api/riddles/answer", me.answerRiddle},
	}
}

func (me *webAppRiddles) getPrimeNumbers(response http.ResponseWriter, request *http.Request) {
	response.Write(encodeJson(globalPrimeNumbers))
}

func (me *webAppRiddles) createRiddle(response http.ResponseWriter, request *http.Request) {
	var item = riddles{}.create(4)
	me.db.insertRiddle(&item)
	item.keys = nil
	response.Write(encodeJson(item))
}

func (me *webAppRiddles) answerRiddle(response http.ResponseWriter, request *http.Request) {
	var id = parseWebInt(request, "id")
	var product = parseWebInt(request, "product")
	var keys []int
	assertWebError(json.NewDecoder(request.Body).Decode(&keys), webError{
		Message: "Invalid JSON body", Status: http.StatusBadRequest,
	})
	var isCorrect = false
	me.db.processRiddle(id, product, func(item *riddleItem) {
		if nil == item {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		var product = 1
		for _, key := range keys {
			product = multiplyLimited(product, key, 1000_000)
		}
		isCorrect = product == item.Product
	})
	response.Write(encodeJson(isCorrect))
}
