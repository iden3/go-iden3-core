package issuer

import (
	"fmt"

	"github.com/iden3/go-iden3-core/components/idensigner"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/idenstatewriter"
)

var (
	dbPrefixClaimsTree     = []byte("treeclaims:")
	dbPrefixRevocationTree = []byte("treerevocation:")
	dbPrefixRootsTree      = []byte("treeroots:")
)

type Config struct {
	MaxLevelsClaimsTree     int
	MaxLevelsRevocationTree int
	MaxLevelsRootsTree      int
}

var ConfigDefault = Config{MaxLevelsClaimsTree: 140, MaxLevelsRevocationTree: 140, MaxLevelsRootsTree: 140}

type Issuer struct {
	id              *core.ID
	claimsMt        *merkletree.MerkleTree
	revMt           *merkletree.MerkleTree
	rootsMt         *merkletree.MerkleTree
	idenStateWriter idenstatewriter.IdenStateWriter
	signer          idensigner.IdenSigner
	storage         db.Storage
	cfg             Config
}

// TODO
func New(storage db.Storage, cfg Config, signer idensigner.IdenSigner, idenStateWriter idenstatewriter.IdenStateWriter) (*Issuer, error) {
	cltStorage := storage.WithPrefix(dbPrefixClaimsTree)
	retStorage := storage.WithPrefix(dbPrefixRevocationTree)
	rotStorage := storage.WithPrefix(dbPrefixRootsTree)

	clt, err := merkletree.NewMerkleTree(cltStorage, cfg.MaxLevelsClaimsTree)
	if err != nil {
		return nil, err
	}
	ret, err := merkletree.NewMerkleTree(retStorage, cfg.MaxLevelsRevocationTree)
	if err != nil {
		return nil, err
	}
	rot, err := merkletree.NewMerkleTree(rotStorage, cfg.MaxLevelsRootsTree)
	if err != nil {
		return nil, err
	}
	return &Issuer{
		id:              nil,
		claimsMt:        clt,
		revMt:           ret,
		rootsMt:         rot,
		idenStateWriter: idenStateWriter,
		signer:          signer,
		storage:         storage,
		cfg:             cfg,
	}, nil
}

func Load(storage db.Storage, signer idensigner.IdenSigner, idenStateWriter idenstatewriter.IdenStateWriter) (*Issuer, error) {
	return nil, fmt.Errorf("TODO")
}

func (i *Issuer) ID() *core.ID {
	return i.id
}

func (i *Issuer) GenCredentialExistence(claim merkletree.Entrier) (*proof.CredentialExistence, error) {
	return nil, fmt.Errorf("TODO")
}

func (i *Issuer) IssueClaim(claim merkletree.Entrier) error {
	return fmt.Errorf("TODO")
}

func (i *Issuer) PublishState() error {
	return fmt.Errorf("TODO")
}

func (i *Issuer) RevokeClaim(claim merkletree.Entrier) error {
	return fmt.Errorf("TODO")
}

func (i *Issuer) UpdateClaim(hIndex *merkletree.Hash, value []merkletree.ElemBytes) error {
	return fmt.Errorf("TODO")
}

func (i *Issuer) Sign(string) (string, error) {
	return "", fmt.Errorf("TODO")
}

func (i *Issuer) SignBinary(string) (string, error) {
	return "", fmt.Errorf("TODO")
}
