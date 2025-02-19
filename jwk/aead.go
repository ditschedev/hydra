/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk

import (
	"context"
	"encoding/base64"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/driver/config"

	"github.com/gtank/cryptopasta"
	"github.com/pkg/errors"
)

type AEAD struct {
	c *config.DefaultProvider
}

func NewAEAD(c *config.DefaultProvider) *AEAD {
	return &AEAD{c: c}
}

func aeadKey(key []byte) *[32]byte {
	var result [32]byte
	copy(result[:], key[:32])
	return &result
}

func (c *AEAD) Encrypt(ctx context.Context, plaintext []byte) (string, error) {
	keys := append([][]byte{c.c.GetGlobalSecret(ctx)}, c.c.GetRotatedGlobalSecrets(ctx)...)
	if len(keys) == 0 {
		return "", errors.Errorf("at least one encryption key must be defined but none were")
	}

	if len(keys[0]) != 32 {
		return "", errors.Errorf("key must be exactly 32 long bytes, got %d bytes", len(keys[0]))
	}

	ciphertext, err := cryptopasta.Encrypt(plaintext, aeadKey(keys[0]))
	if err != nil {
		return "", errorsx.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *AEAD) Decrypt(ctx context.Context, ciphertext string) (p []byte, err error) {
	keys := append([][]byte{c.c.GetGlobalSecret(ctx)}, c.c.GetRotatedGlobalSecrets(ctx)...)
	if len(keys) == 0 {
		return nil, errors.Errorf("at least one decryption key must be defined but none were")
	}

	for _, key := range keys {
		if p, err = c.decrypt(ciphertext, key); err == nil {
			return p, nil
		}
	}

	return nil, err
}

func (c *AEAD) decrypt(ciphertext string, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.Errorf("key must be exactly 32 long bytes, got %d bytes", len(key))
	}

	raw, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	plaintext, err := cryptopasta.Decrypt(raw, aeadKey(key))
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return plaintext, nil
}
