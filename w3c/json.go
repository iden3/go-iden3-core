package w3c

import "encoding/json"

func (did *DID) UnmarshalJSON(bytes []byte) error {
	var didStr string
	err := json.Unmarshal(bytes, &didStr)
	if err != nil {
		return err
	}

	did3, err := ParseDID(didStr)
	if err != nil {
		return err
	}
	*did = *did3
	return nil
}

func (did DID) MarshalJSON() ([]byte, error) {
	return json.Marshal(did.String())
}
