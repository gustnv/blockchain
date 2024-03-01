package wallet

import (
    "bytes"
    "crypto/elliptic"
    "crypto/ecdsa"
    "encoding/gob"
    "io/ioutil"
    "log"
	"math/big"
    "os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
    Wallets map[string]*Wallet
}

func CreateWallets() (*Wallets, error) {
    wallets := Wallets{}
    wallets.Wallets = make(map[string]*Wallet)

    err := wallets.LoadFile()

    return &wallets, err
}

func (ws *Wallets) AddWallet() string {
    wallet := MakeWallet()
    address := string(wallet.Address())

    ws.Wallets[address] = wallet

    return address
}

func (ws *Wallets) GetAllAddresses() []string {
    var addresses []string

    for address := range ws.Wallets {
        addresses = append(addresses, address)
    }

    return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
    return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile() error {
    if _, err := os.Stat(walletFile); os.IsNotExist(err) {
        return err
    }

    var wallets Wallets

    fileContent, err := ioutil.ReadFile(walletFile)
    if err != nil {
        return err
    }

    gob.Register(elliptic.P256())
    decoder := gob.NewDecoder(bytes.NewReader(fileContent))
    err = decoder.Decode(&wallets)
    if err != nil {
        return err
    }

    for _, wallet := range wallets.Wallets {
		curve := elliptic.P256()
		x := new(big.Int).SetBytes(wallet.PublicKey[:len(wallet.PublicKey)/2]) 
		y := new(big.Int).SetBytes(wallet.PublicKey[len(wallet.PublicKey)/2:]) 
   
		privKey := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: curve,
				X:     x,
				Y:     y,
			},
			D: new(big.Int).SetBytes(wallet.PrivateKey),  // Extract the big.Int
		} 
   
		// Convert the private key to bytes before storing
		wallet.PrivateKey = privKey.D.Bytes() 
   }

    ws.Wallets = wallets.Wallets

    return nil
}

func (ws *Wallets) SaveFile() {
    var content bytes.Buffer

    gob.Register(elliptic.P256())

    encoder := gob.NewEncoder(&content)
    err := encoder.Encode(ws)
    if err != nil {
        log.Panic(err)
    }

    err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
    if err != nil {
        log.Panic(err)
    }
}
