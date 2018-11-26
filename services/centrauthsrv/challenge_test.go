package centrauthsrv

import (
	"strconv"
	"testing"
	"time"
)

func TestVerifyChallengeTimestamp(t *testing.T) {
	challenge := "uuid-" + strconv.Itoa(int(time.Now().Unix())) + "-randstr"
	err := VerifyChallengeTimestamp(challenge)
	if err != nil {
		t.Errorf(err.Error())
	}
}
