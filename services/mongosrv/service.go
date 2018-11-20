package mongosrv

import mgo "gopkg.in/mgo.v2"

type Service interface {
	GetCollections() map[string]*mgo.Collection
}

type ServiceImpl struct {
	// session mgo.Session
	collections map[string]*mgo.Collection
}

func New(url string, databaseName string, collectionsArray []string) (*ServiceImpl, error) {
	session, err := mgo.Dial("mongodb://" + url)
	if err != nil {
		return &ServiceImpl{}, err
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// get all the collections specified in the collectionsArray
	collections := make(map[string]*mgo.Collection)
	for _, collection := range collectionsArray {
		collections[collection] = session.DB(databaseName).C(collection)
	}
	return &ServiceImpl{collections}, nil
}

// GetCollections returns the current mongodb collections of the mongosrv
func (mongosrv *ServiceImpl) GetCollections() map[string]*mgo.Collection {
	return mongosrv.collections
}
