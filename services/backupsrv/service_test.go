package backupsrv

import (
	"encoding/json"
	"testing"

	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
)

func TestPoWBackupData(t *testing.T) {
	backupDataJSON := `
		{
      "idaddrhex" : "0xa6d0cfcb340d63e092d94f1418ed747e06600ce8",
      "data" : "ob5+ScOOVN63jQoIuqQabXWYNQV54Wilxt0F9v6lQX9HD6qIJCoNoeb57yZB68akelkq",
      "datasignature" : "0x9a18e8de5f783c248d274233a2c7df143c89df3463b63ca75cb5eceb912cac197e6a963c85a660d05d97d29359f7ac533e5a12b44a7e1fc5ca16d418e540a6cf1c",
      "type" : "testtype",
      "ksign" : "0xafe7617048aa672cf291b67681aaec12739a1bc6",
      "proofofksignhex" : {
          "claimproof" : {
              "leaf" : "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000afe7617048aa672cf291b67681aaec12739a1bc62077bb3f0400dd62421c97220536fd6ed2be29228e8db1315e8c6d7525f4bdf4dad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb34062000000005b816b9e000000005b816b9e",
              "proof" : "0x000000000000000000000000000000000000000000000000000000000000000142fc7e9dec0b9f522d25471a8402841fb64a87b144b20f2173ee35ce60f06421",
              "root" : "0x2e497cfc68e1bf0173991556a79a19d1ab32ce470d09ef4ef8f26d108d8461aa"
          },
          "setrootclaimproof" : {
              "leaf" : "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000001a6d0cfcb340d63e092d94f1418ed747e06600ce82e497cfc68e1bf0173991556a79a19d1ab32ce470d09ef4ef8f26d108d8461aa",
              "proof" : "0x00000000000000000000000000000000000000000000000000000000000000015afcaa2ff30f41b318b74a92701932083a612562257356df111e18c447ebc36d",
              "root" : "0x05e2729903acfde930ad687011291f321d71674470eae0136a30ea076391494b"
          },
          "claimnonrevocationproof" : {
              "leaf" : "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000001afe7617048aa672cf291b67681aaec12739a1bc62077bb3f0400dd62421c97220536fd6ed2be29228e8db1315e8c6d7525f4bdf4dad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb34062000000005b816b9e000000005b816b9e",
              "proof" : "0x0000000000000000000000000000000000000000000000000000000000000003170f3c6f9a16d26c58a803eb1c82ac194ea9a21f1ba813f436e5b99e1752583942fc7e9dec0b9f522d25471a8402841fb64a87b144b20f2173ee35ce60f06421",
              "root" : "0x2e497cfc68e1bf0173991556a79a19d1ab32ce470d09ef4ef8f26d108d8461aa"
          },
          "setrootclaimnonrevocationproof" : {
              "leaf" : "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000002a6d0cfcb340d63e092d94f1418ed747e06600ce82e497cfc68e1bf0173991556a79a19d1ab32ce470d09ef4ef8f26d108d8461aa",
              "proof" : "0x0000000000000000000000000000000000000000000000000000000000000009f51742846a2aba8a633aa5e82b1a4224cf1989bb69714856b6dee41c71e4d81c5afcaa2ff30f41b318b74a92701932083a612562257356df111e18c447ebc36d",
              "root" : "0x05e2729903acfde930ad687011291f321d71674470eae0136a30ea076391494b"
          },
          "date" : 1542617349,
          "signature" : "0x69cdee010706666d9aeab310378310e3d62276c9bcbfb5bcc1c8a80a39b6377500b392a4d2500b3fe12a8369906e6e3346fefb31ec6f59d641714c8f9047684100"
      },
      "relayaddr" : "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c",
      "timestamp" : 1542617353
    }
	`

	var backupData BackupData
	json.Unmarshal([]byte(backupDataJSON), &backupData)
	assert.Equal(t, backupData.ProofOfKSignHex.SetRootClaimProof.Root, "0x05e2729903acfde930ad687011291f321d71674470eae0136a30ea076391494b")

	data, err := utils.PoW(backupData, 2)
	assert.Nil(t, err)
	b, err := json.Marshal(data)
	assert.Nil(t, err)
	hash := merkletree.HashBytes(b)
	assert.True(t, utils.CheckPoW(hash, 2))
}
