// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/perfectbui/chat/clients"
	"github.com/perfectbui/chat/config"
	"github.com/perfectbui/chat/job"
	"google.golang.org/grpc"

	"github.com/perfectbui/chat/middlewares"
)

var addr = flag.String("addr", ":7070", "http service address")
var authClientAddr = flag.String("authClientAddr", ":5000", "http service address")
var authClientHost = flag.String("authClientHost", "localhost", "http service address")

func loadAuthClient() {
	conn, err := grpc.Dial(*authClientHost+*authClientAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	clients.LoadAuthClient(conn)
}

func main() {
	flag.Parse()
	config.CreateRedisClient()
	job.LoadProducer()
	// config.InitProducer()
	// config.InitConsumer()
	loadAuthClient()
	wsServer := newWsServer()
	// http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access-token")
		if err != nil {
			return
		}
		refreshToken, err := r.Cookie("refresh-token")
		if err != nil {
			return
		}
		_, _, err = middlewares.CheckAuth(accessToken.Value, refreshToken.Value)
		if err != nil {
			return
		}

		userID, ok := r.URL.Query()["userID"]
		if !ok || len(userID[0]) < 1 {
			log.Println("Url Param 'userID' is missing")
			return
		}
		n, err := strconv.ParseInt(userID[0], 10, 64)
		if err != nil {
			return
		}
		fmt.Printf("nhan ket noi tu %v \n", userID[0])
		serveWs(wsServer, w, r, n)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
