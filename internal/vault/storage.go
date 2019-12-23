// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/lfaoro/spark/ent"
	"github.com/lfaoro/spark/ent/card"
	"github.com/lfaoro/spark/internal"
	pb "github.com/lfaoro/spark/proto/api/vault"
)

func (s Server) storeCard(ctx context.Context, w *pb.PutCardResponse, cipheredCard []byte) error {
	expOn, err := ptypes.Timestamp(w.ExpiresOn)
	if err != nil {
		return errors.Wrap(err, "failed to transform timestamp")
	}
	autoDelOn, err := ptypes.Timestamp(w.AutoDeleteOn)
	if err != nil {
		return errors.Wrap(err, "failed to transform timestamp")
	}

	// todo: remove fake id generator\
	// var buf [32]byte
	// _, err = rand.Read(buf[:])
	// if err != nil {
	// 	panic(err) // out of randomness, should never happen
	// }
	// id := rand.Intn(1000)

	id, ok := internal.GetValueFrom(ctx, "project_id").(int)
	if !ok {
		return fmt.Errorf("failed to get project_id from context")
	}

	dbcard, err := s.db.Card.Create().
		SetOwnerID(id).
		SetCipheredCard(cipheredCard).
		SetToken(w.Token).
		SetHash(w.Hash).
		SetFirstSix(w.FirstSix).
		SetLastFour(w.LastFour).
		SetExpiresOn(expOn).
		SetAutoDeleteOn(autoDelOn).
		SetRequestIP(w.RequestIp).
		SetUserAgent(w.UserAgent).
		Save(ctx)
	if err != nil {
		return err
	}

	// bail if no metadata
	if w.Metadata == nil || w.Metadata.ServiceMessage != "" {
		return nil
	} else {
		m := w.Metadata
		_, err = s.db.Metadata.Create().
			SetOwnerID(dbcard.ID).
			SetScheme(m.Scheme).
			SetBrand(m.Brand).
			SetType(m.Type).
			SetCurrency(m.Currency).
			SetIsPrepaid(m.IsPrepaid).
			SetIssuerName(m.Issuer.Name).
			SetIssuerURL(m.Issuer.Name).
			SetIssuerPhone(m.Issuer.Phone).
			SetIssuerCity(m.Issuer.City).
			SetIssuerCountry(m.Issuer.Country).
			SetIssuerCountrycode(m.Issuer.CountryCode).
			SetIssuerLatitude(float64(m.Issuer.Latitude)).
			SetIssuerLongitude(float64(m.Issuer.Longitude)).
			SetIssuerMap(m.Issuer.Map).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Server) findCard(ctx context.Context, token string) (*ent.Card, error) {
	value, err := s.db.Card.Query().
		Where(card.Token(token)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s Server) findCardMetadata(ctx context.Context, token string) (*ent.Metadata, error) {
	value, err := s.db.Card.Query().
		Where(card.Token(token)).
		QueryMetadata().
		First(ctx)
	if err != nil {
		return nil, err
	}
	return value, nil
}
