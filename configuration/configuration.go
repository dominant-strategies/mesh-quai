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

package configuration

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dominant-strategies/mesh-quai/ethereum"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/dominant-strategies/go-quai/params"
)

// Mode is the setting that determines if
// the implementation is "online" or "offline".
type Mode string

const (
	// Online is when the implementation is permitted
	// to make outbound connections.
	Online Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline Mode = "OFFLINE"

	// Mainnet is the Quai Mainnet (aka Colosseum).
	Mainnet string = "MAINNET"

	// Orchard is the Quai Orchard testnet.
	Orchard string = "ORCHARD"

	// Local is the Quai Local testnet.
	Local string = "LOCAL"

	// DataDirectory is the default location for all
	// persistent data.
	DataDirectory = "/data"

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"

	// GoQuaiEnv is an optional environment variable
	// used to connect rosetta-ethereum to an already
	// running geth node.
	GoQuaiEnv = "GOQUAI"

	// DefaultGoQuaiURL is the default URL for
	// a running go-quai node. This is used
	// when GoQuaiEnv is not populated.
	DefaultGoQuaiURL = "http://localhost:8545"

	// SkipGoQuaiAdminEnv is an optional environment variable
	// to skip geth `admin` calls which are typically not supported
	// by hosted node services. When not set, defaults to false.
	SkipGoQuaiAdminEnv = "SKIP_GO_QUAI_ADMIN"

	// MiddlewareVersion is the version of rosetta-ethereum.
	MiddlewareVersion = "0.0.4"
)

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	GenesisBlockIdentifier *types.BlockIdentifier
	GoQuaiURL              string
	RemoteGoQuai           bool
	Port                   int
	GoQuaiArguments        string
	SkipGoQuaiAdmin        bool

	// Block Reward Data
	Params *params.ChainConfig
}

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration() (*Configuration, error) {
	config := &Configuration{}

	modeValue := Mode(os.Getenv(ModeEnv))
	switch modeValue {
	case Online:
		config.Mode = Online
	case Offline:
		config.Mode = Offline
	case "":
		return nil, errors.New("MODE must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	networkValue := os.Getenv(NetworkEnv)
	switch networkValue {
	case Mainnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.MainnetNetwork,
		}
		config.GenesisBlockIdentifier = ethereum.MainnetGenesisBlockIdentifier
		config.Params = params.ProgpowColosseumChainConfig
		config.GoQuaiArguments = ethereum.MainnetGoQuaiArguments
	case Orchard:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.OrchardNetwork,
		}
		config.GenesisBlockIdentifier = ethereum.OrchardGenesisBlockIdentifier
		config.Params = params.ProgpowOrchardChainConfig
		config.GoQuaiArguments = ethereum.OrchardGoQuaiArguments
	case Local:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.DevNetwork,
		}
		config.GenesisBlockIdentifier = nil
		config.Params = params.ProgpowLocalChainConfig
		config.GoQuaiArguments = ethereum.LocalGoQuaiArguments
	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	config.GoQuaiURL = DefaultGoQuaiURL
	envGethURL := os.Getenv(GoQuaiEnv)
	if len(envGethURL) > 0 {
		config.RemoteGoQuai = true
		config.GoQuaiURL = envGethURL
	}

	config.SkipGoQuaiAdmin = false
	envSkipGethAdmin := os.Getenv(SkipGoQuaiAdminEnv)
	if len(envSkipGethAdmin) > 0 {
		val, err := strconv.ParseBool(envSkipGethAdmin)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to parse SKIP_GO_QUAI_ADMIN %s", err, envSkipGethAdmin)
		}
		config.SkipGoQuaiAdmin = val
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("%w: unable to parse port %s", err, portValue)
	}
	config.Port = port

	return config, nil
}
