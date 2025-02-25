// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethereum

import (
	"log"

	"github.com/dominant-strategies/go-quai/common"
)

// ChecksumAddress ensures an Ethereum hex address
// is in Checksum Format. If the address cannot be converted,
// it returns !ok.
func ChecksumAddress(address string, location common.Location) (common.Address, bool) {
	addr, err := common.NewMixedcaseAddressFromString(address, location)
	if err != nil {
		return common.Address{}, false
	}

	return addr.Address(), true
}

// MustChecksum ensures an address can be converted
// into a valid checksum. If it does not, the program
// will exit.
func MustChecksum(address string, location common.Location) common.Address {
	addr, ok := ChecksumAddress(address, location)
	if !ok {
		log.Fatalf("invalid address %s", address)
	}

	return addr
}
