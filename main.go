package main

import (
	"os"

	"github.com/gustnv/blockchain_project/cli"
	"github.com/gustnv/blockchain_project/wallet"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()

	wallet := wallet.MakeWallet()
	wallet.Address()
	
}
