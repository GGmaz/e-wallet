package main

import (
	cfg "github.com/GGmaz/wallet-arringo/config"
	srv "github.com/GGmaz/wallet-arringo/internal/server"
)

func main() {
	config := cfg.NewConfig()
	server := srv.NewServer(config)
	server.Start()
}
