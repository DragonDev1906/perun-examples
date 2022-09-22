// Copyright 2022 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"github.com/perun-network/perun-eth-backend/wire"

	"perun.network/go-perun/channel"
)

// CreateContractBackend creates a new contract backend.
func CreateContractBackend(
	nodeURL string,
	chainID *big.Int,
	w *swallet.Wallet,
) (ethchannel.ContractBackend, error) {
	signer := types.LatestSignerForChainID(chainID)
	transactor := swallet.NewTransactor(w, signer)

	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return ethchannel.ContractBackend{}, err
	}

	return ethchannel.NewContractBackend(ethClient, ethchannel.MakeChainID(chainID), transactor, txFinalityDepth), nil
}

// WalletAddress returns the wallet address of the client.
func (c *PaymentClient) WalletAddress() common.Address {
	return common.Address(*c.account.(*ethwallet.Address))
}

// WireAddress returns the wire address of the client.
func (c *PaymentClient) WireAddress() *wire.Address {
	return &wire.Address{Address: ethwallet.AsWalletAddr(c.WalletAddress())}
}

// EthToWei converts a given amount in ETH to Wei.
func EthToWei(ethAmount *big.Float) (weiAmount *big.Int) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).Mul(ethAmount, weiPerEthFloat)
	weiAmount, _ = weiAmountFloat.Int(nil)
	return weiAmount
}

// WeiToEth converts a given amount in Wei to ETH.
func WeiToEth(weiAmount *big.Int) (ethAmount *big.Float) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).SetInt(weiAmount)
	return new(big.Float).Quo(weiAmountFloat, weiPerEthFloat)
}

// getChainAssets returns the assets on the different chains.
func getChainsAssets(chains []ChainConfig) []channel.Asset {
	assets := make([]channel.Asset, len(chains))
	for i, chain := range chains {
		assets[i] = ethchannel.NewAsset(chain.ChainID.Int, chain.AssetHolder)
	}
	return assets
}
