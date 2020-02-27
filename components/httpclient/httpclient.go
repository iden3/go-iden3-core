package httpclient

import (
	"fmt"
	"reflect"

	"github.com/dghubble/sling"
	"gopkg.in/go-playground/validator.v9"
)

type ServerError struct {
	Err string `json:"error"`
}

func (e ServerError) Error() string {
	return fmt.Sprintf("server: %v", e.Err)
}

type HttpClient struct {
	UrlBase  string
	client   *sling.Sling
	validate *validator.Validate
}

func NewHttpClient(urlBase string) *HttpClient {
	if urlBase[len(urlBase)-1] != '/' {
		urlBase += "/"
	}
	client := sling.New().Base(urlBase)
	return &HttpClient{UrlBase: urlBase,
		client: client, validate: validator.New()}
}

func (p *HttpClient) NewRequest() *sling.Sling {
	return p.client.New()
}

func (p *HttpClient) DoRequest(s *sling.Sling, res interface{}) error {
	var serverError ServerError
	resp, err := s.Receive(res, &serverError)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return serverError
	}
	if res != nil {
		rv := reflect.ValueOf(res)
		if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
			return p.validate.Struct(res)
		}
	}
	return nil
}
