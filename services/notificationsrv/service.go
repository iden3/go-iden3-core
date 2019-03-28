package notificationsrv

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
	"github.com/iden3/go-iden3/core"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// Notification defines the fields to be sent to the server
type Notification struct {
	Type string				`json:"type" binding:"required"`
  Data interface{}  `json:"data" binding:"required"`
}


// ResLogin defines data response from http notification server call GET url/login
type ResLogin struct {
	SignedPacket signedpacketsrv.RequestIdenAssert `json:"sigReq" binding:"required"`
}

// ResSubmit defines data response from http notification server call POST url/login
type ResSubmit struct {
	Token string `json:"token" binding:"required"`
}

// Service defines fields to notification service 
type Service struct {
	URL       string
	IDAddr common.Address
	Token     string
}

// New creates new notification service
func New(urlService string, idAddr common.Address) *Service {
	return &Service{urlService, idAddr, ""}
}

// Login get signed packet from service, signs the packeta and submit login to receive login token
func (ns *Service) Login(keyStore *keystore.KeyStore, keySignPk *ecdsa.PublicKey,
	proofKSign core.ProofClaim, idenAssertForm *signedpacketsrv.IdenAssertForm ) error {
	// Get signed packet from notification server
	req, err := requestLogin(ns.URL + "/login")
	
	// Sign packet
	signedPacket, err := signedpacketsrv.NewSignIdenAssertV01(req, idenAssertForm,
		keyStore, ns.IDAddr, keySignPk, proofKSign, 600)
	if err != nil {
		return err
	}
	// Submit sign packet and get login token
	ns.Token, err = submitLogin(ns.URL + "/login", signedPacket)
	if err != nil {
		return err
	}
	
	return err
}

// requestLogin get signed packet from service
func requestLogin(url string) (*signedpacketsrv.RequestIdenAssert, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resLogin ResLogin
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &resLogin); err != nil {
		return nil, err
	}
	
	return &resLogin.SignedPacket, err
}

// submitLogin get login token from service
func submitLogin(url string, signedPacket *signedpacketsrv.SignedPacket) (string, error) {
	sp, err := signedPacket.MarshalJSON()
	if err != nil {
		return "", err
	}
	var resSubmit ResSubmit
	res, err := http.Post(url, "application/json", bytes.NewBuffer(sp))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(data, &resSubmit); err != nil {
		return "", err
	}
	return resSubmit.Token, nil
}

// SendNotification send notification to service
func(ns *Service) SendNotification(not *Notification) error {
	// Check service has alredy a token
	if ns.Token == "" {
		return fmt.Errorf("No login token has been found")
	}
	// Post notification
	bytesNot, err := json.Marshal(not)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", ns.URL + "/auth/notifications/" + ns.IDAddr.Hex(), bytes.NewBuffer(bytesNot))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer " + ns.Token)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
