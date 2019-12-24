// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcpkms

import (
	"context"
	"fmt"
	"log"

	gcpkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/grpc/connectivity"

	"github.com/pkg/errors"
	"google.golang.org/api/iterator"

	"github.com/lfaoro/spark/pkg/kms"
)

type GCPKMS struct {
	projectID  string
	locationID string
	ctx        context.Context
	client     *gcpkms.KeyManagementClient
	parent     string
	debug      bool
}

type Option func(*GCPKMS)

func Debug(v bool) Option {
	return func(c *GCPKMS) {
		c.debug = v
	}
}

func New(projectID, locationID string, opt ...Option) kms.Service {
	ctx := context.Background()
	client, err := gcpkms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	hsm := &GCPKMS{
		projectID:  projectID,
		locationID: locationID,

		ctx:    ctx,
		client: client,
		parent: fmt.Sprintf("projects/%s/locations/%s",
			projectID, locationID),
	}

	for _, o := range opt {
		o(hsm)
	}

	// TODO: consider returning a list of keys  instead
	//  of assuming the first key is what we want.
	hsm.setupKeys()

	return hsm
}

func (c GCPKMS) Encrypt(plaintext []byte) ([]byte, error) {
	res, err := c.client.Encrypt(c.ctx, &kmspb.EncryptRequest{
		Name:      c.parent,
		Plaintext: plaintext,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt plaintext")
	}

	return res.Ciphertext, nil
}

func (c GCPKMS) Decrypt(ciphertext []byte) ([]byte, error) {
	res, err := c.client.Decrypt(c.ctx, &kmspb.DecryptRequest{
		Name:       c.parent,
		Ciphertext: ciphertext,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt ciphertext")
	}
	return res.Plaintext, nil
}

func (c GCPKMS) State() connectivity.State {
	return c.client.Connection().GetState()
}

func (c *GCPKMS) setupKeys() {
	// TODO: investigate key versions
	lkr := c.client.ListKeyRings(c.ctx,
		&kmspb.ListKeyRingsRequest{
			Parent: c.parent,
		})
	for {
		res, err := lkr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed to list key rings: %v", err)
		}
		if c.debug {
			fmt.Printf("KeyRing: %q\n", res.Name)
		}
		c.parent = res.Name
		break
	}

	lck := c.client.ListCryptoKeys(c.ctx, &kmspb.ListCryptoKeysRequest{
		Parent: c.parent,
	})
	for {
		res, err := lck.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal("failed to list keys", err)
		}
		if c.debug {
			fmt.Println("keys:", res.Name)
		}
		c.parent = res.Name
		break
	}

}
