package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/pion/webrtc/v3"
)

func Trickle(conn net.Conn) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
				},
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	defer func() {
		cErr := pc.Close()
		if cErr != nil {
			panic(cErr)
		}
	}()

	// trickle candidates in OnIceConnection
	// send candidate to client as soon as it's ready
	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		msg, err := json.Marshal(i.ToJSON().Candidate)
		if err != nil {
			panic(err)
		}

		// send candidate info over HTTP connection
		if _, err := conn.Write(msg); err != nil {
			panic(err)
		}
	})

	// set handler for ICE connection state
	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		fmt.Println("ICE connection state changed", pcs.String())
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			for range time.Tick(time.Second * 3) {
				if err := dc.SendText(time.Now().String()); err != nil {
					panic(err)
				}
			}
		})
	})
}
