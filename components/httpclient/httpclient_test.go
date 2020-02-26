package httpclient

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Res struct {
	A string `json:"a" validate:"required"`
	B string `json:"b" validate:"required"`
}

func TestHttpClient(t *testing.T) {

	myHandler := http.NewServeMux()
	myHandler.HandleFunc("/good", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "a": "test1", "b": "test2" }`)
	})
	myHandler.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "a": "test1" }`)
	})
	s := &http.Server{
		Addr:           "127.0.0.1:14567",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	sErrCh := make(chan error)
	go func() {
		sErrCh <- s.ListenAndServe()
	}()
	var err error
	for i := 0; i < 10; i++ {
		conn, errDial := net.Dial("tcp", "127.0.0.1:14567")
		err = errDial
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}
	require.Nil(t, err)

	var res0 Res
	// Unmarshal the response into a struct with validation
	httpClient := NewHttpClient("http://127.0.0.1:14567")
	err = httpClient.DoRequest(httpClient.NewRequest().Path("good").Get(""), &res0)
	require.Nil(t, err)

	// Should error because there's a missing required field
	var res1 Res
	err = httpClient.DoRequest(httpClient.NewRequest().Path("bad").Get(""), &res1)
	require.NotNil(t, err)

	// Not interested in the response
	err = httpClient.DoRequest(httpClient.NewRequest().Path("good").Get(""), nil)
	require.Nil(t, err)

	// Unmarshal the response into a map
	res2 := make(map[string]string)
	err = httpClient.DoRequest(httpClient.NewRequest().Path("good").Get(""), &res2)
	require.Nil(t, err)
	require.Equal(t, map[string]string{"a": "test1", "b": "test2"}, res2)

	s.Close()
	err = <-sErrCh
	require.Equal(t, http.ErrServerClosed, err)
}
