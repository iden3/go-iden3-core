package notificationsrv

import (
	"bytes"
	"fmt"
	"io/ioutil"

	// "crypto/ecdsa"
	"encoding/json"
	// "io/ioutil"
	"net/http"

	// "github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/services/signedpacketsrv"
)

const MSGPROOFCLAIMV01 = "iden3.proofclaim.v0_1"
const MSGTXT = "iden3.txt.v0_1"

type HttpError struct {
	Response   *http.Response
	StatusCode int
	Body       string
}

func NewHttpError(response *http.Response) (error, error) {
	if response.StatusCode == 200 {
		return nil, nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &HttpError{Response: response, StatusCode: response.StatusCode, Body: string(body)}, nil
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP error %d: %s", e.StatusCode, e.Body)
}

// Notification defines the fields to be sent to the server
type Notification struct {
	Type string      `json:"type" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

// Service defines fields to notification service
type Service struct {
	URL                string
	Token              string
	signedPacketSigner *signedpacketsrv.SignedPacketSigner
}

// New creates new notification service
func New(urlService string, signedPacketSigner *signedpacketsrv.SignedPacketSigner) *Service {
	return &Service{urlService, "", signedPacketSigner}
}

// Login get signed packet from service, signs the packeta and submit login to receive login token
// func (ns *Service) Login(keyStore *keystore.KeyStore, keySignPk *ecdsa.PublicKey,
// 	proofKSign core.ProofClaim, idenAssertForm *signedpacketsrv.IdenAssertForm) error {
// 	// Get signed packet from notification server
// 	req, err := requestLogin(ns.URL + "/login")
//
// 	// Sign packet
// 	signedPacket, err := signedpacketsrv.NewSignIdenAssertV01(req, idenAssertForm,
// 		keyStore, ns.IDAddr, keySignPk, proofKSign, 600)
// 	if err != nil {
// 		return err
// 	}
// 	// Submit sign packet and get login token
// 	ns.Token, err = submitLogin(ns.URL+"/login", signedPacket)
// 	if err != nil {
// 		return err
// 	}
//
// 	return err
// }

// ResLogin defines data response from http notification server call GET url/login
// type ResLogin struct {
// 	SignedPacket signedpacketsrv.RequestIdenAssert `json:"sigReq" binding:"required"`
// }

// requestLogin get signed packet from service
// func requestLogin(url string) (*signedpacketsrv.RequestIdenAssert, error) {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	var resLogin ResLogin
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = json.Unmarshal(data, &resLogin); err != nil {
// 		return nil, err
// 	}
//
// 	return &resLogin.SignedPacket, err
// }

// ResSubmit defines data response from http notification server call POST url/login
// type ResSubmit struct {
// 	Token string `json:"token" binding:"required"`
// }

// submitLogin get login token from service
// func submitLogin(url string, signedPacket *signedpacketsrv.SignedPacket) (string, error) {
// 	sp, err := signedPacket.MarshalJSON()
// 	if err != nil {
// 		return "", err
// 	}
// 	var resSubmit ResSubmit
// 	res, err := http.Post(url, "application/json", bytes.NewBuffer(sp))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer res.Body.Close()
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	if err = json.Unmarshal(data, &resSubmit); err != nil {
// 		return "", err
// 	}
// 	return resSubmit.Token, nil
// }

// SendNotification send notification to service
// func (ns *Service) SendNotification(not *Notification) error {
// 	// Check service has alredy a token
// 	if ns.Token == "" {
// 		return fmt.Errorf("No login token has been found")
// 	}
// 	// Post notification
// 	bytesNot, err := json.Marshal(not)
// 	if err != nil {
// 		return err
// 	}
// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", ns.URL+"/auth/notifications/"+ns.IDAddr.String(), bytes.NewBuffer(bytesNot))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Authorization", "Bearer "+ns.Token)
// 	_, err = client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func NewMsgTxt(txt string) Notification {
	return Notification{
		Type: MSGTXT,
		Data: txt,
	}
}

func NewMsgProofClaim(proofClaim *core.ProofClaim) Notification {
	return Notification{
		Type: MSGPROOFCLAIMV01,
		Data: proofClaim,
	}
}

// SendNotification send notification to service
func (ns *Service) SendNotification(notif Notification, id core.ID) error {
	signedPacket, err := ns.signedPacketSigner.
		NewSignMsgV01(0xdeadbeef, notif.Type, notif.Data)
	if err != nil {
		return err
	}
	signedPacketJSON, err := json.Marshal(signedPacket)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/notifications/%s", ns.URL, id.String()),
		"application/json", bytes.NewBuffer(signedPacketJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	httpErr, err := NewHttpError(resp)
	if err != nil {
		return err
	}
	return httpErr
}
