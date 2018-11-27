package centrauthsrv

import (
	"errors"
	"strconv"
	"strings"

	"github.com/iden3/go-iden3/utils"
)

// VerifyChallengeTimestamp checks that the given timestamp is correct
func VerifyChallengeTimestamp(challenge string) error {
	// verify challenge timestamp < 30 seconds ago
	if len(strings.Split(challenge, "-")) < 2 {
		return errors.New("VerifyChallengeTimestamp: challenge timestamp error")
	}
	unixTimeChallenge, err := strconv.Atoi(strings.Split(challenge, "-")[1])
	if err != nil {
		return errors.New("VerifyChallengeTimestamp: challenge timestamp error")
	}
	// t := time.Unix(int64(unixTimeChallenge), 10)
	// elapsed := time.Since(t)
	// if elapsed.Seconds() > 30000 { // 30 seconds to resolve challenge // DEV in development we use more time
	// 	return errors.New("VerifyTimstamp: too much time elapsed since the challenge was sent")
	// }
	verified := utils.VerifyTimestamp(uint64(unixTimeChallenge), 30000) // 30 seconds to resolve challenge // DEV in development we use more time
	if !verified {
		return errors.New("VerifyTimstamp: too much time elapsed since the challenge was sent")
	}
	return nil
}
