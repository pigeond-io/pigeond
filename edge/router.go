// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package edge

import (
	"github.com/pigeond-io/pigeond/common/log"
	"github.com/pigeond-io/pigeond/edge/client"
	"github.com/pigeond-io/pigeond/edge/client/message"
	"github.com/pigeond-io/pigeond/edge/hub"
	"net/http"
	"strconv"
	"sync"
)

func Init(wsPort int, udpPort int, bufferSize int) {
	log.Info("Edge server initialization started....")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		initClientListener(wsPort)
	}()

	go func() {
		defer wg.Done()
		initHubListener(udpPort, bufferSize)
	}()

	log.Info("Edge server initialization done ")

	wg.Wait()
}

func initClientListener(port int) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Websocket connection")
		client.Handler(w, r, message.DefaultMessageReader{})
	})

	addr := ":" + strconv.Itoa(port)
	log.Info("Client listening on websocket : ", addr)
	response := http.ListenAndServe(addr, nil)
	log.Fatal(response)
}

func initHubListener(port int, bufferSize int) {
	hub.Listen(port, bufferSize)
}
