package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	w               http.ResponseWriter
	statusCode      int
	responseMessage interface{}
}

func (r *Response) jsonResponse() {
	response, _ := json.Marshal(r.responseMessage)

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)
	r.w.Write(response)
}

func (a *App) getStatusOfOperation(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var taskId int
	if len(params) == 0 {
		response := Response{
			w:               w,
			statusCode:      http.StatusBadRequest,
			responseMessage: map[string]interface{}{"success": false, "result": "Request must contain parameter 'task_id'"},
		}
		response.jsonResponse()
		return
	}
	taskId, err := strconv.Atoi(params["task_id"][0])
	if err != nil {
		response := Response{
			w:               w,
			statusCode:      http.StatusBadRequest,
			responseMessage: map[string]interface{}{"success": false, "result": "Conversion error occurred"},
		}
		response.jsonResponse()
		return
	}
	result, ok := taskIdsToStatus[taskId]
	if !ok {
		response := Response{
			w:               w,
			statusCode:      http.StatusBadRequest,
			responseMessage: map[string]interface{}{"success": false, "result": "Task id does not exist"},
		}
		response.jsonResponse()
		return
	}
	if result != nil {
		response := Response{
			w:               w,
			statusCode:      http.StatusOK,
			responseMessage: map[string]interface{}{"success": true, "result": result},
		}
		response.jsonResponse()
		return
	} else {
		response := Response{
			w:               w,
			statusCode:      http.StatusOK,
			responseMessage: map[string]interface{}{"success": false, "result": "Operation has not yet been completed"},
		}
		response.jsonResponse()
	}
}

func (a *App) retrieveGoods(w http.ResponseWriter, r *http.Request) {
	goodsList := make([]Goods, 0)
	params := r.URL.Query()
	sellerId := params["seller_id"]
	offerId := params["offer_id"]
	query := params["query"]
	if len(sellerId) == 0 {
		sellerId = append(sellerId, "")
	}
	if len(offerId) == 0 {
		offerId = append(offerId, "")
	}
	if len(query) == 0 {
		query = append(query, "")
	}
	err, gs := a.getGoods(sellerId[0], offerId[0], query[0])
	if err != nil {
		response := Response{
			w:               w,
			statusCode:      http.StatusBadRequest,
			responseMessage: map[string]interface{}{"status": false, "result": "Failed to retrieve goods from a database"},
		}
		response.jsonResponse()
		return
	}
	for i := range gs {
		var goods Goods
		b, err := json.Marshal(gs[i])
		if err != nil {
			fmt.Println(err)
		}
		err_ := json.Unmarshal(b, &goods)
		if err_ != nil {
			fmt.Println(err_)
		}

		goodsList = append(goodsList, goods)
	}

	response := Response{
		w:               w,
		statusCode:      http.StatusOK,
		responseMessage: map[string]interface{}{"status": true, "result": goodsList},
	}
	response.jsonResponse()
	return
}

func (a *App) loadGoodsAsync(sellerId int, file multipart.File, result chan Stats, taskId int, done chan bool) {
	deleted := 0
	created := 0
	updated := 0
	errs := 0

	goods := parseExcelFile(file)

	for i := 0; i < len(goods); i++ {
		if a.checkSellerOfferExistence(sellerId, goods[i].OfferId) {
			del, err := a.updateGoods(*goods[i])
			if err != nil {
				errs++
			}
			if del {
				deleted++
			} else {
				updated++
			}
		} else {
			err := a.insertGoods(*goods[i], sellerId)
			if err != nil {
				errs++
			} else {
				created++
			}
		}
	}
	var stats Stats
	statMap := map[string]int{"Created": created, "Updated": updated, "Deleted": deleted, "Errors": errs}
	bs, _ := json.Marshal(statMap)

	json.Unmarshal(bs, &stats)

	taskIdsToStatus[taskId] = stats
}

func (a *App) loadGoods(w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 10000000
	taskId := rand.Intn(max - min + 1)

	sellerId, err := strconv.Atoi(r.FormValue("seller_id"))
	if err != nil {
		response := Response{
			w:               w,
			statusCode:      http.StatusBadRequest,
			responseMessage: map[string]interface{}{"status": false, "result": err},
		}
		response.jsonResponse()
		return
	}

	_, header, err := r.FormFile("goods_file")
	if err != nil {
		response := Response{
			w:               w,
			statusCode:      http.StatusBadRequest,
			responseMessage: map[string]interface{}{"status": false, "result": err},
		}
		response.jsonResponse()
		return
	}

	file, err := header.Open()
	if err != nil {
		response := Response{
			w:               w,
			statusCode:      http.StatusInternalServerError,
			responseMessage: map[string]interface{}{"status": false, "result": "Failed to open a file"},
		}
		response.jsonResponse()
		return
	}

	defer file.Close()
	result := make(chan Stats)
	done := make(chan bool)
	go a.loadGoodsAsync(sellerId, file, result, taskId, done)
	taskIdsToStatus[taskId] = nil
	w.Header().Set("Connection", "Keep-Alive")
	response := Response{
		w:               w,
		statusCode:      http.StatusOK,
		responseMessage: map[string]interface{}{"success": true, "result": taskId},
	}
	response.jsonResponse()
}
