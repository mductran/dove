package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/pion/webrtc/v3"
)

func trickle(conn net.Conn) {
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

	buf := make([]byte, 1024)
	for {
		// read and process http responses
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		var candidate webrtc.ICECandidateInit
		var offer webrtc.SessionDescription

		switch {
		case json.Unmarshal(buf[:n], &offer) == nil && offer.SDP != "":
			if err = pc.SetRemoteDescription(offer); err != nil {
				panic(err)
			}

			answer, answerErr := pc.CreateAnswer(nil)
			if answerErr != nil {
				panic(answerErr)
			}

			if err = pc.SetLocalDescription(offer); err != nil {
				panic(err)
			}

			outbound, marshalErr := json.Marshal(answer)
			if marshalErr != nil {
				panic(marshalErr)
			}

			if _, err = conn.Write(outbound); err != nil {
				panic(err)
			}

		case json.Unmarshal(buf[:n], &candidate) == nil && candidate.Candidate != "":
			if err = pc.AddICECandidate(candidate); err != nil {
				panic(err)
			}
		default:
			panic("Unknown message")
		}
	}
}
