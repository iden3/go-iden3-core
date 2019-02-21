#!/usr/bin/env python3
"""go-iden3/cmd/claimserver endpoints test
"""

import requests
import provoj

URL = "http://127.0.0.1:6000/api/unstable"

t = provoj.NewTest("claimserver")

r = requests.get(URL + "/root")
t.rStatus("get root", r)

aux = "0x" + str.encode("asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf").hex()
r = requests.post(URL + "/claims", json={"indexData": aux, "data": aux})
t.rStatus("post claim", r)
postClaimRes = r.json()

hi = "0x15a329f60308d935a9665bb922b36b3bbdd031260cab1a3cef027b0055dea55f"
r = requests.get(URL + "/claims/" + hi + "/proof")
t.rStatus("get claim proof by hi", r)
getClaimRes = r.json()

t.equal("postClaim response == getClaim response (proofClaim)", postClaimRes, getClaimRes)

t.printScores()
