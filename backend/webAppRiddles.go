package main

import (
	"net/http"
	"time"
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
	var product = riddles{}.create(4)
	var row = riddleRow{
		product:   product,
		createdAt: time.Now(),
	}
	me.db.insertRiddle(&row)
	var responseObject = riddleResponse{
		Id:      row.id,
		Product: row.product,
		Steps:   me.steps,
	}
	response.Write(encodeJson(responseObject))
}

func (me *webAppRiddles) answerRiddle(response http.ResponseWriter, request *http.Request) {
	var id = parseWebInt(request, "id")
	var product = parseWebInt(request, "product")
	var keys []int
	decodeWebJson(request.Body, &keys)
	var isCorrect = false
	me.db.processRiddle(id, product, func(item *riddleRow) {
		if nil == item {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		var product = 1
		for _, key := range keys {
			product = multiplyLimited(product, key, 1000_000)
		}
		isCorrect = product == item.product
	})
	response.Write(encodeJson(isCorrect))
}
