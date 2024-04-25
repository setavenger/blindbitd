package src

func (u *OwnedUTXO) LabelPubKey() []byte {
	if u.Label != nil {
		return u.Label.PubKey[:]
	} else {
		return nil
	}
}
