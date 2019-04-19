package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var (
	tv     Tv
	screen Screen
)

type Response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func ErrorResponse(rsp http.ResponseWriter) {
	payload := Response{
		Status: "error",
		Result: nil,
	}
	rsp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rsp).Encode(payload)
}

func SuccessResponse(rsp http.ResponseWriter, data interface{}) {
	payload := Response{
		Status: "ok",
		Result: data,
	}
	rsp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rsp).Encode(payload)
}

//打开电视
type openTvPayload struct {
}

func openTv(rsp http.ResponseWriter, req *http.Request) {
	var payload openTvPayload
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		ErrorResponse(rsp)
		return
	}
	//1. 屏幕打开状态
	if screen.isOpen() {
		log.Println("screen is open")
		if openTvAction() {
			SuccessResponse(rsp, nil)
			return
		}
		ErrorResponse(rsp)
		return
	}
	log.Println("screen is close")
	//2. 屏幕关闭状态
	err = screen.open()
	if err != nil {
		log.Println("screen open fail")
		log.Println(err)
		ErrorResponse(rsp)
		return
	}
	log.Println("screen open success")
	if openTvAction() {
		SuccessResponse(rsp, nil)
		return
	}
	ErrorResponse(rsp)
	return
}

func openTvAction() bool {
	//open tv action
	if tv.isOpen() {
		log.Println("tv is open")
		return true
	}
	err := tv.open()
	if err != nil {
		log.Println("tv open fail")
		log.Println(err)
		return false
	}
	log.Println("tv open success")
	return true
}

//关闭电视
type closeTvPayload struct {
}

func closeTv(rsp http.ResponseWriter, req *http.Request) {
	var payload closeTvPayload
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		ErrorResponse(rsp)
		return
	}
	//1. tv关闭状态
	if !tv.isOpen() {
		log.Println("tv is close")
		if closeScreenAction() {
			SuccessResponse(rsp, nil)
			return
		}
		ErrorResponse(rsp)
		return
	}
	log.Println("tv is open")
	//2. tv打开状态
	err = tv.close()
	if err != nil {
		log.Println("tv close fail")
		log.Println(err)
		ErrorResponse(rsp)
		return
	}
	log.Println("tv close success")
	for tv.isOpen() {
		time.Sleep(time.Second * 1)
	}
	if closeScreenAction() {
		SuccessResponse(rsp, nil)
		return
	}
	ErrorResponse(rsp)
	return
}

func closeScreenAction() bool {
	//close screen action
	if !screen.isOpen() {
		log.Println("screen is close")
		return true
	}
	err := screen.close()
	if err != nil {
		log.Println("screen close fail")
		log.Println(err)
		return false
	}
	log.Println("screen close success")
	return true
}

//暂停
type pauseTvPayload struct {
}

func pauseTv(rsp http.ResponseWriter, req *http.Request) {
	var payload pauseTvPayload
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		ErrorResponse(rsp)
		return
	}
	err = tv.pause()
	if err != nil {
		log.Println(err)
		ErrorResponse(rsp)
		return
	}
	SuccessResponse(rsp, nil)
	return
}
