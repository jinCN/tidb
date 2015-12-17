// Copyright 2014 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"strings"

	"github.com/juju/errors"
	"github.com/pingcap/tidb/kv"
)

// ScanMetaWithPrefix scans metadata with the prefix.
func ScanMetaWithPrefix(retriever kv.Retriever, prefix string, filter func([]byte, []byte) bool) error {
	iter, err := retriever.Seek([]byte(prefix))
	if err != nil {
		return errors.Trace(err)
	}
	defer iter.Close()

	for {
		if err != nil {
			return errors.Trace(err)
		}

		if iter.Valid() && strings.HasPrefix(iter.Key(), prefix) {
			if !filter([]byte(iter.Key()), iter.Value()) {
				break
			}
			err = iter.Next()
			if err != nil {
				return errors.Trace(err)
			}
		} else {
			break
		}
	}

	return nil
}

// DelKeyWithPrefix deletes keys with prefix.
func DelKeyWithPrefix(rm kv.RetrieverMutator, prefix string) error {
	var keys []string
	iter, err := rm.Seek([]byte(prefix))
	if err != nil {
		return errors.Trace(err)
	}

	defer iter.Close()
	for {
		if err != nil {
			return errors.Trace(err)
		}

		if iter.Valid() && strings.HasPrefix(iter.Key(), prefix) {
			keys = append(keys, iter.Key())
			err = iter.Next()
			if err != nil {
				return errors.Trace(err)
			}
		} else {
			break
		}
	}

	for _, key := range keys {
		err := rm.Delete([]byte(key))
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// RowKeyPrefixFilter returns a function which checks whether currentKey has decoded rowKeyPrefix as prefix.
func RowKeyPrefixFilter(rowKeyPrefix []byte) kv.FnKeyCmp {
	return func(currentKey kv.Key) bool {
		// Next until key without prefix of this record.
		return !bytes.HasPrefix(currentKey, rowKeyPrefix)
	}
}
