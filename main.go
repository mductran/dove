package main

import (
	"fmt"
	"os"

	"time"

	"dove/signal"

	"github.com/pion/webrtc/v3"
)

func main() {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println("error creating peer connection")
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Println("error closing peer connection")
		}
	}()

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		fmt.Println("peer connection state changed: ", state.String())

		if state == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			os.Exit(0)
		}
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Println("new data channe: ", d.Label(), d.ID())

		// register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message := signal.RandSeq(15) // TODO: replace with sending actual data
				fmt.Printf("Sending '%s' \n", message)

				// send message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})

		// message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Receiving message %s '%s'", d.Label(), string(msg.Data))
		})
	})

	offer := webrtc.SessionDescription{} // create session description https://webrtcforthecurious.com/docs/02-signaling/#full-example
	signal.Decode(signal.MustReadStdin(), &offer)

	// set remote session description
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection) // channel blocked until ICE Candidate gathering is complete

	// set local session description
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))

	// Block forever
	select {}

}
