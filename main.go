package main

import (
	"regexp"
	"errors"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/oauth2"
	"sync"
)

var server = "lava.cosmicpe.me:19132"
var did = "69c17a07d86a4bc3bb5bf0c92399dba3"

func main() {

	token, err := auth.RequestLiveToken()
	if err != nil {
		panic(err)
	}

	src := auth.RefreshTokenSource(token)

	p, err := minecraft.NewForeignStatusProvider(server)

	if err != nil {
		panic(err)
	}

	listener, err := minecraft.ListenConfig{StatusProvider: p}.Listen("raknet", "127.0.0.1:19132")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleBot(c.(*minecraft.Conn), listener, src)
	}

}

func handleBot(conn *minecraft.Conn, listener *minecraft.Listener, src oauth2.TokenSource){

	serverConn, err := minecraft.Dialer{
		TokenSource: src,
		ClientData: login.ClientData{
			DeviceID: did,
		},
	}.Dial("raknet", server)

	if err != nil {
		panic(err)
	}

	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}
		g.Done()
	}()

	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		g.Done()
	}()
	g.Wait()

	go func(){
		defer listener.Disconnect(conn, "disconnected")
		defer serverConn.Close()
		for {
			pk, err := conn.ReadPacket()

			if err != nil {
				panic(err)
			}

			if err := serverConn.WritePacket(pk); err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(conn, disconnect.Error())
				}
				return
			}

			//client side here

		}
	}()

	go func(){
		defer serverConn.Close()
		defer listener.Disconnect(conn, "disconnected")

		for {

			pk, err := serverConn.ReadPacket()
			if err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(conn, disconnect.Error())
				}
				return
			}

			if err := conn.WritePacket(pk); err != nil {
				return
			}

			if text, ok := pk.(*packet.Text); ok {
				fmt.Println(stripColors(text.Message))
			}



		}
	}()

}

var re = regexp.MustCompile(`(?i)§[0-9A-GK-OR]`)

func stripColors(message string) string {
	var str = []byte(message)
	var sub = []byte(" ")
	str = re.ReplaceAll(str, sub)
	return string(str)
}