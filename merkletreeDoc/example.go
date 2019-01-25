package main

func main() {
}

/*
func main() {
	storage, err := db.NewLevelDbStorage("./path", false)
	if err != nil {
		panic(err)
	}
	mt, err := merkletree.New(storage, 140) // new merkletree of 140 levels of depth
	if err != nil {
		panic(err)
	}
	defer mt.Storage().Close()

	indexData := []byte("this is a first claim, this is the data that we put in the claim, that will affect it's position in the merkletree")
	data := []byte("data that we put in the claim, that will not affect it's position in the merkletree")
	claim0 := core.NewGenericClaim("namespace", "default", indexData, data)

	indexData = []byte("this is a second claim")
	data = []byte("data that we put in the claim, that will not affect it's position in the merkletree")
	claim1 := core.NewGenericClaim("namespace", "default", indexData, data)

	fmt.Println("adding claim0")
	err = mt.Add(claim0)
	if err != nil {
		panic(err)
	}
	fmt.Println("merkle root: " + mt.Root().Hex())
	fmt.Println("adding claim1")
	err = mt.Add(claim1)
	if err != nil {
		panic(err)
	}

	mp, err := mt.GenerateProof(claim0.Hi())
	if err != nil {
		panic(err)
	}
	fmt.Println("merkle root: " + mt.Root().Hex())

	mpHex := common3.HexEncode(mp)
	fmt.Println("merkle proof: " + mpHex)
	checked := merkletree.CheckProof(mt.Root(), mp, claim0.Hi(), claim0.Ht(), mt.NumLevels())
	fmt.Println("merkle proof checked:", checked)

	claimInPosBytes, err := mt.GetValueInPos(claim0.Hi())
	if err != nil {
		panic(err)
	}
	// print true if the claimInPosBytes is the same than claim0.Bytes()
	fmt.Println("claim in position equals to the original:", bytes.Equal(claim0.Bytes(), claimInPosBytes))

	indexData = []byte("this claim will not be stored")
	data = []byte("")
	claim2 := core.NewGenericClaim("namespace", "default", indexData, data)

	mp, err = mt.GenerateProof(claim2.Hi())
	if err != nil {
		panic(err)
	}

	checked = merkletree.CheckProof(mt.Root(), mp, claim2.Hi(), merkletree.EmptyNodeValue, mt.NumLevels())

	fmt.Println("merkle proof of non existence checked:", checked)
}
*/
