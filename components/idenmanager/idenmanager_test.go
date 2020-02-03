package idenmanager

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/components/idensigner"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/db"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	idenstatewritemock "github.com/iden3/go-iden3-core/services/idenstatewriter/mock"
	"github.com/iden3/go-iden3-crypto/babyjub"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var debug = true

var service *IdenManager
var idenStateWriter *idenstatewritemock.IdenStateWriteMock
var mt *merkletree.MerkleTree

const relayIdHex = "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf"
const relaySkHex = "4be5471a938bdf3606888472878baace4a6a64e14a153adf9a1333969e4e573c"

var relayID core.ID
var relaySk babyjub.PrivateKey
var relayPkComp *babyjub.PublicKeyComp
var relayPk *babyjub.PublicKey

var rmDirs []string

func newTestingMerkle(numLevels int) (*merkletree.MerkleTree, error) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}
	sto, err := db.NewLevelDbStorage(dir, false)
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}

	mt, err := merkletree.NewMerkleTree(sto, numLevels)
	return mt, err
}

/*
func loadConfig() {
	c.Server.Port = "5000"
	c.Server.PrivK = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
	c.Geth.URL = ""
	c.ContractsAddress.Identities = "0x101d2fa51f8259df207115af9eaa73f3f4e52e60"
	c.Domain = "iden3.io"
	c.Namespace = "iden3.io"
}
*/
func initializeEnvironment(t *testing.T) {

	// MerkleTree leveldb
	var err error
	mt, err = newTestingMerkle(140)
	if err != nil {
		t.Error(err)
	}

	pass := []byte("my passphrase")
	storage := babykeystore.MemStorage([]byte{})
	keyStore, err := babykeystore.NewKeyStore(&storage, babykeystore.LightKeyStoreParams)
	if err != nil {
		panic(err)
	}

	if _, err := hex.Decode(relaySk[:], []byte(relaySkHex)); err != nil {
		panic(err)
	}
	if relayPkComp, err = keyStore.ImportKey(relaySk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(relayPkComp, pass); err != nil {
		panic(err)
	}
	if relayPk, err = relayPkComp.Decompress(); err != nil {
		panic(err)
	}

	signer := idensigner.New(keyStore, *relayPk)

	relayID, err = core.IDFromString(relayIdHex)
	if err != nil {
		panic(err)
	}
	idenStateWriter = idenstatewritemock.New()
	service = New(&relayID, mt, idenStateWriter, *signer)
}

func TestGetNextVersion(t *testing.T) {
	initializeEnvironment(t)

	indexData := []byte("c1")
	data := []byte{}
	var indexSlot [claims.IndexSlotBytes]byte
	var dataSlot [claims.DataSlotBytes]byte
	// copy(indexSlot[:], indexData[:400/8])
	// copy(dataSlot[:], data[:496/8])
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])
	claim := claims.NewClaimBasic(indexSlot, dataSlot, 0)

	version, err := GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(0), version)
	claim.Version = version
	if err := mt.AddClaim(claim); err != nil {
		panic(err)
	}
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(1), version)

	claim.Version = version
	if err := mt.AddClaim(claim); err != nil {
		panic(err)
	}
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(2), version)

	claim.Version = version
	if err := mt.AddClaim(claim); err != nil {
		panic(err)
	}
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(3), version)

	claim.Version = version
	if err := mt.AddClaim(claim); err != nil {
		panic(err)
	}
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(4), version)
}

func TestGetNonRevocationProof(t *testing.T) {
	initializeEnvironment(t)
	indexData := []byte("c1")
	data := []byte{}
	var indexSlot [claims.IndexSlotBytes]byte
	var dataSlot [claims.DataSlotBytes]byte
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])
	claim := claims.NewClaimBasic(indexSlot, dataSlot, 0)

	if err := mt.AddClaim(claim); err != nil {
		panic(err)
	}
	version, err := GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(1), version)

	claimProof, err := getNonRevocationProof(mt, claim.Entry().HIndex())
	require.Nil(t, err)

	require.Equal(t,
		"0x030000000000000000000000000000000000000000000000000000000000000055340ba059d27b18de92b8ac4fb42a158bfb2e389ae3294ffeb132af91a7da1c81b7e168ae02356bb67bef5d540125f040ee4d7bff6eb64a35cfcdf2d5761a02",
		common3.HexEncode(claimProof.Proof))
	require.Equal(t,
		"0x187f605d25901b5a2c563f491220ec6a54ef47367a55cc5c22e71bdda8a26c0d",
		claimProof.Root.Hex())

	proof, err := merkletree.NewProofFromBytes(claimProof.Proof)
	require.Nil(t, err)

	// VerifyProof with HIndex and HValue from claimProof.Leaf
	var leafBytes [merkletree.ElemBytesLen * merkletree.DataLen]byte
	copy(leafBytes[:], claimProof.Leaf)
	dataNonE := merkletree.NewDataFromBytes(leafBytes)
	claimProofEntry := merkletree.Entry{
		Data: *dataNonE,
	}
	verified := merkletree.VerifyProof(&claimProof.Root, proof, claimProofEntry.HIndex(), claimProofEntry.HValue())
	require.True(t, verified)
}

func TestGetClaimProof(t *testing.T) {
	initializeEnvironment(t)

	// id, err := core.IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
	// require.Nil(t, err)

	// Basic Claim
	indexData := []byte("index01")
	data := []byte("data01")
	var indexSlot [claims.IndexSlotBytes]byte
	var dataSlot [claims.DataSlotBytes]byte
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])
	claimBasic := claims.NewClaimBasic(indexSlot, dataSlot, 0)

	var kSignSk babyjub.PrivateKey
	if _, err := hex.Decode(kSignSk[:], []byte("9b3260823e7b07dd26ef357ccfed23c10bcef1c85940baa3d02bbf29461bbbbe")); err != nil {
		panic(err)
	}
	kSignPk := kSignSk.Public()
	claimAuthKSignBabyJub := claims.NewClaimAuthorizeKSignBabyJub(kSignPk, 0)

	// TMP commented due ClaimAuthorizeKSignSecp256k1 is not updated yet to new spec
	// // KSign Claim
	// sk, err := crypto.HexToECDSA("0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4")
	// if err != nil {
	//         panic(err)
	// }
	// pk := sk.Public().(*ecdsa.PublicKey)
	// claimAuthKSign := core.NewClaimAuthorizeKSignSecp256k1(pk)

	// open the MerkleTree of the user
	userMT, err := newTestingMerkle(140)
	require.Nil(t, err)

	// add claimBasic in User ID Merkle Tree
	err = userMT.AddClaim(claimBasic)
	require.Nil(t, err)

	// // add claimAuthKSign in User ID Merkle Tree
	// err = userMT.AddClaim(claimAuthKSign)
	// require.Nil(t, err)

	// add claimAuthKSignBabyJub in User ID Merkle Tree
	err = userMT.AddClaim(claimAuthKSignBabyJub)
	require.Nil(t, err)

	/*
					// TMP commented due SetRootClaim is not updated yet to new spec
						// setRootClaim of the user in the Relay Merkle Tree
						setRootClaim, err := core.NewClaimSetRootKey(&id, userMT.RootKey())
						require.Nil(t, err)
						// setRootClaim.BaseIndex.Version++ // TODO autoincrement
						// add User's ID Merkle Root into the Relay's Merkle Tree
						err = mt.AddClaim(setRootClaim)
						require.Nil(t, err)

				idenStateWriter.On("GetRoot", &relayID).Return(
					&core.RootData{BlockN: 123, BlockTimestamp: 456, Root: mt.RootKey()}, nil).Once()
					proofClaim, err := service.GetClaimProofByHiBlockchain(setRootClaim.Entry().HIndex())
				require.Nil(t, err)
				p, err := json.Marshal(proofClaim)
				require.Nil(t, err)
				if debug {
					fmt.Printf("\n\tSetRoot claim proof\n\n")
					fmt.Println(string(p))
				}

			ok, err := proofClaim.Verify(proofClaim.Proof.Root)
			if !ok || err != nil {
				panic(err)
			}

		// ClaimAssignName
		// id, err := core.IDFromString("1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY")
		// require.Nil(t, err)
		claimAssignName := core.NewClaimAssignName("testName@iden3.eth", id)
		// add assignNameClaim in User ID Merkle Tree
		err = mt.AddClaim(claimAssignName)
		require.Nil(t, err)
		fmt.Printf("> A %+v\n", mt.RootKey().String())
		idenStateWriter.On("GetRoot", &relayID).Return(
			&core.RootData{BlockN: 123, BlockTimestamp: 456, Root: mt.RootKey()}, nil).Once()
		fmt.Printf("> R %+v\n", mt.RootKey().String())
		proofClaimAssignName, err := service.GetClaimProofByHiBlockchain(claimAssignName.Entry().HIndex())
		require.Nil(t, err)
		p, err = json.Marshal(proofClaimAssignName)
		if debug {
			fmt.Printf("\n\tclaim assign name claim proof\n\n")
			fmt.Println(string(p))
		}
	*/

	//proofClaim, err := service.GetClaimProofUserByHi(ethAddr, *claim.Entry().HIndex())
	//require.Nil(t, err)
	//if err != nil {
	//	panic(err)
	//}
	//claimProof := proofClaim.ClaimProof
	//claimNonRevocationProof := proofClaim.ClaimNonRevocationProof
	//setRootClaimProof := proofClaim.SetRootClaimProof
	//setRootClaimNonRevocationProof := proofClaim.SetRootClaimNonRevocationProof

	//require.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db000000048000000006461746161736466", common3.HexEncode(claimProof.Leaf))
	//require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.HexEncode(claimProof.Proof))
	//require.Equal(t, "0x1415376b054a9ab3c7f9bd0ec956b0f403ae98d7e37dcbafdf26b465b23dd970", claimProof.Root.Hex())
	//require.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000970e8128ab834e8eac17ab8e3812f010678cf7911415376b054a9ab3c7f9bd0ec956b0f403ae98d7e37dcbafdf26b465b23dd970", common3.HexEncode(setRootClaimProof.Leaf))
	//require.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000474c3e76aebd3df03ff91325d245e72ea9ad9599777f5d2c5e560b3f049d68309", common3.HexEncode(setRootClaimProof.Proof))
	//require.Equal(t, "0xf73c98cbaa1d43ada4ed5520300c348985dd47cc283e3cf7186434a07a46886a", setRootClaimProof.Root.Hex())
	//require.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000025563046fb69f065953f0fdb0b3033f721457184adfae2824c02932090bf8f281", common3.HexEncode(claimNonRevocationProof.Proof))
	//require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000014367d7c39348c9b7f2d488a7cd2edfbae56d608ec92a1b0a747adda3c4aaf763d74c3e76aebd3df03ff91325d245e72ea9ad9599777f5d2c5e560b3f049d68309", common3.HexEncode(setRootClaimNonRevocationProof.Proof))

	//require.Equal(t, claimProof.Root.Bytes(), claimNonRevocationProof.Root.Bytes())
	//require.Equal(t, setRootClaimProof.Root.Bytes(), setRootClaimNonRevocationProof.Root.Bytes())

	//var leafBytes [128]byte
	//copy(leafBytes[:], claimProof.Leaf)
	//entry := merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err := merkletree.NewProofFromBytes(claimProof.Proof)
	//require.Nil(t, err)
	//verified := merkletree.VerifyProof(&claimProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], setRootClaimProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(setRootClaimProof.Proof)
	//require.Nil(t, err)
	//verified = merkletree.VerifyProof(&setRootClaimProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], claimNonRevocationProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(claimNonRevocationProof.Proof)
	//require.Nil(t, err)
	//verified = merkletree.VerifyProof(&claimNonRevocationProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], setRootClaimNonRevocationProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(setRootClaimNonRevocationProof.Proof)
	//require.Nil(t, err)
	//verified = merkletree.VerifyProof(&setRootClaimNonRevocationProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)
}

func initializeIdService(t *testing.T) *IdenManager {
	relayId, err := core.IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	if err != nil {
		t.Error(err)
	}

	// MerkleTree leveldb
	mt, err := newTestingMerkle(140)
	if err != nil {
		t.Error(err)
	}
	// sto := db.NewMemoryStorage()
	idenStateWriteMock := idenstatewritemock.New()
	idenStateWriteMock.On("SetRoot", mock.Anything).Return()

	pass := []byte("my passphrase")
	storage := babykeystore.MemStorage([]byte{})
	keyStore, err := babykeystore.NewKeyStore(&storage, babykeystore.LightKeyStoreParams)
	if err != nil {
		panic(err)
	}

	if _, err := hex.Decode(relaySk[:], []byte(relaySkHex)); err != nil {
		panic(err)
	}
	if relayPkComp, err = keyStore.ImportKey(relaySk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(relayPkComp, pass); err != nil {
		panic(err)
	}
	if relayPk, err = relayPkComp.Decompress(); err != nil {
		panic(err)
	}

	signer := idensigner.New(keyStore, *relayPk)

	return New(&relayId, mt, idenStateWriteMock, *signer)
}

func TestCreateIdGenesisRandomLoop(t *testing.T) {
	idsrv := initializeIdService(t)

	// turn this to 'true' to compute this test. Currently disabled as needs more than 100s to compute
	if false {
		for i := 0; i < 1024; i++ {
			kOpSk := babyjub.NewRandPrivKey()
			kop := kOpSk.Public()
			if debug {
				fmt.Println("kop", kop)
			}
			kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
			kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
			kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

			id, proofKOp, err := idsrv.CreateIdGenesis(kop, kDis, kReen, kUpdateRoot)
			require.Nil(t, err)

			id2, _, err := genesis.CalculateIdGenesisFrom4Keys(kop, kDis, kReen, kUpdateRoot)
			require.Nil(t, err)
			require.Equal(t, id, id2)

			proofKOpVerified, err := proofKOp.Verify(proofKOp.Proof.Root)
			require.Nil(t, err)
			require.True(t, proofKOpVerified)
		}
	}
}

/*
// TMP commented due the ClaimAuthorizeKSignSecp256k is not updated to new spec and causes crash

func TestCreateIdGenesisHardcoded(t *testing.T) {
	idsrv := initializeIdService(t)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	// krecStr := "0x03f9737be33b5829e3da80160464b2891277dae7d7c23609f9bb34bd4ede397bbf"
	// krevStr := "0x02d2da59d3022b4c1589e4910baa6cbaddd01f95ed198fdc3068d9dc1fb784a9a4"

	var kopComp babyjub.PublicKeyComp
	err := kopComp.UnmarshalText([]byte(kopStr))
	require.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	require.Nil(t, err)
	kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, proofKOp, err := idsrv.CreateIdGenesis(kopPub, kDis, kReen, kUpdateRoot)
	require.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	require.Equal(t, "1KVNQLxwiiXyFzVwJqMPUhRMk6TeXEbzXhz2R6aw2", id.String())

	id2, _, err := core.CalculateIdGenesisFrom4Keys(kopPub, kDis, kReen, kUpdateRoot)
	require.Nil(t, err)
	require.Equal(t, id, id2)

	proofKOpVerified, err := proofKOp.Verify(proofKOp.Proof.Root)
	require.Nil(t, err)
	require.True(t, proofKOpVerified)
}
*/

func TestMain(m *testing.M) {
	result := m.Run()
	for _, dir := range rmDirs {
		os.RemoveAll(dir)
	}
	os.Exit(result)
}
