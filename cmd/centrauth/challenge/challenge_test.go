package challenge

import (
	"strconv"
	"testing"
	"time"
)

func TestVerifyTimestamp(t *testing.T) {
	challenge := "uuid-" + strconv.Itoa(int(time.Now().Unix())) + "-randstr"
	err := VerifyTimestamp(challenge)
	if err != nil {
		t.Errorf(err.Error())
	}
}
