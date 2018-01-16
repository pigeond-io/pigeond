// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package edge

import (
	. "github.com/pigeond-io/pigeond/common"
	"github.com/pigeond-io/pigeond/edge/client"
	"github.com/pigeond-io/pigeond/edge/client/message"
	"github.com/pigeond-io/pigeond/edge/hub"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func Init() {
	Info.Print("Edge server initialization started....")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		initClientListener(8002)
	}()

	go func() {
		defer wg.Done()
		initHubListener(8001)
	}()

	Info.Print("Edge server initialization done ")

	wg.Wait()
}

func initClientListener(port int) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Websocket connection")
		client.Handler(w, r, message.DefaultMessageReader{})
	})

	addr := ":" + strconv.Itoa(port)
	Info.Print("Client listening on websocket : ", addr)
	response := http.ListenAndServe(addr, nil)
	Error.Println(response)
}

func initHubListener(port int) {
	hub.Listen(port, 2048)
}
