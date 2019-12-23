// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"bytes"
	"encoding/gob"

	"github.com/pkg/errors"

	pb "github.com/lfaoro/spark/proto/api/vault"
)

func (s Server) encryptCard(data *pb.PutCardRequest) ([]byte, error) {
	ec, err := encodeCard(data)
	if err != nil {
		return nil, err
	}

	ciphertext, err := s.kms.Encrypt(ec)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (s Server) decryptCard(data []byte) (*pb.GetCardResponse, error) {
	plaintext, err := s.kms.Decrypt(data)
	if err != nil {
		return nil, err
	}

	vc, err := decodeCard(plaintext)
	if err != nil {
		return nil, err
	}

	return vc, err
}

type cardData struct {
	Holder   string
	Number   string
	ExpYear  uint32
	ExpMonth uint32
	CVC      uint32
}

func encodeCard(card *pb.PutCardRequest) ([]byte, error) {
	cd := cardData{
		Number:   card.Number,
		ExpYear:  card.ExpYear,
		ExpMonth: card.ExpMonth,
		CVC:      card.Cvc,
	}
	data := bytes.Buffer{}
	err := gob.NewEncoder(&data).Encode(cd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode cardData")
	}

	return data.Bytes(), nil
}

func decodeCard(card []byte) (*pb.GetCardResponse, error) {
	cd := &cardData{}
	data := bytes.NewReader(card)
	err := gob.NewDecoder(data).Decode(cd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode cardData")
	}

	vc := &pb.GetCardResponse{
		Holder:   cd.Holder,
		Number:   cd.Number,
		ExpMonth: cd.ExpMonth,
		ExpYear:  cd.ExpYear,
		Cvc:      cd.CVC,
	}

	return vc, nil
}
