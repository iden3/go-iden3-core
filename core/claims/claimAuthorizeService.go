package claims

// import (
// 	"encoding/binary"
//
// 	"github.com/iden3/go-iden3-core/merkletree"
// )
//
// // ServiceType
// var (
// 	// ServiceTypeRelay is the type for authorize Relays
// 	ServiceTypeRelay = NewServiceType(0)
// 	// ServiceTypeNotificationsServer is the type for authorize Notification Server
// 	ServiceTypeNotificationsServer = NewServiceType(1)
// 	// ServiceTypeDiscoveryNode is the type for authorize DiscoveryNode
// 	ServiceTypeDiscoveryNode = NewServiceType(2)
// )
//
// // ServiceTypeLen is the length in bytes of the type of the Services
// const ServiceTypeLen = 64 / 8
//
// // ServiceType is the type used to store a claim type.
// type ServiceType [ServiceTypeLen]byte
//
// // NewServiceType to set type of authorized services
// func NewServiceType(num uint64) *ServiceType {
// 	st := ServiceType{}
// 	binary.BigEndian.PutUint64(st[:], num)
// 	return &st
// }
//
// // ClaimAuthorizeService is a claim to authorize a Service for the identity that performs the claim
// type ClaimAuthorizeService struct {
// 	// Version is the claim version.
// 	Version uint32
// 	// ServiceType is the type of the authorized service
// 	ServiceType *ServiceType
// 	// ServiceAddr is the hash of the addr
// 	ServiceAddr [248 / 8]byte
// 	// ServicePubK is the hash of the pubK
// 	ServicePubK [248 / 8]byte
// 	// ServiceUrl is the hash of the domain
// 	ServiceUrl [248 / 8]byte
// }
//
// // NewClaimAuthorizeService returns a ClaimAuthorizeService with the provided data.
// func NewClaimAuthorizeService(serviceType *ServiceType, serviceAddr, servicePubK, serviceUrl string) *ClaimAuthorizeService {
// 	return &ClaimAuthorizeService{
// 		Version:     0,
// 		ServiceType: serviceType,
// 		ServiceAddr: HashString(serviceAddr),
// 		ServicePubK: HashString(servicePubK),
// 		ServiceUrl:  HashString(serviceUrl),
// 	}
// }
//
// // NewClaimAuthorizeServiceFromEntry deserializes a ClaimAuthorizeService from an Entry.
// func NewClaimAuthorizeServiceFromEntry(e *merkletree.Entry) *ClaimAuthorizeService {
// 	c := &ClaimAuthorizeService{}
// 	_, c.Version = GetClaimTypeVersion(e)
// 	var serviceType [64 / 8]byte
// 	copyFromElemBytes(serviceType[:], ClaimTypeVersionLen, &e.Data[3])
// 	c.ServiceType = NewServiceType(binary.BigEndian.Uint64(serviceType[:]))
// 	copyFromElemBytes(c.ServiceAddr[:], 0, &e.Data[2])
// 	copyFromElemBytes(c.ServicePubK[:], 0, &e.Data[1])
// 	copyFromElemBytes(c.ServiceUrl[:], 0, &e.Data[0])
// 	return c
// }
//
// // Entry serializes the claim into an Entry.
// func (c *ClaimAuthorizeService) Entry() *merkletree.Entry {
// 	e := &merkletree.Entry{}
// 	SetClaimTypeVersion(e, c.Type(), c.Version)
// 	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, c.ServiceType[:])
// 	copyToElemBytes(&e.Data[2], 0, c.ServiceAddr[:])
// 	copyToElemBytes(&e.Data[1], 0, c.ServicePubK[:])
// 	copyToElemBytes(&e.Data[0], 0, c.ServiceUrl[:])
// 	return e
// }
//
// // Type returns the ClaimType of the claim.
// func (c *ClaimAuthorizeService) Type() ClaimType {
// 	return *ClaimTypeAuthorizeService
// }
