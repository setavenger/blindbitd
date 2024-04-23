package main

import "github.com/setavenger/blindbitd/src"

func main() {

	wallet := src.NewWallet()
	c := src.Client{BaseUrl: "http://localhost:8000"}

	daemon := src.Daemon{
		Wallet: wallet,
		Client: &c,
	}

	daemon.Run()

}
