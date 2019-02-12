#!/usr/bin/env python3
"""Test for iden3/go-iden3/cmd/backupserver
"""

import json
import requests
import provoj

URL = "http://127.0.0.1:5001/api/unstable"
username = 'user0'
password = 'pass0'
backup_data = 'this is the backup data'

T = provoj.NewTest("Backup")

R = requests.post(URL + "/register", json={'username': username, 'password': password})
T.rStatus("register", R)

R = requests.post(URL + "/backup/upload", json={'username': username, 'password': password, 'backup': backup_data})
T.rStatus("register", R)

R = requests.post(URL + "/backup/download", json={'username': username, 'password': password})
T.rStatus("register", R)
jsonR = R.json()
T.equal("checking backup username", jsonR["username"], username)
T.equal("checking backup data", jsonR["backup"], backup_data)

T.printScores()
