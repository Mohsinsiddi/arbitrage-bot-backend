package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ArbitrageContract wraps interaction with the arbitrage contract
type ArbitrageContract struct {
	address    common.Address
	abi        abi.ABI
	client     *ethclient.Client
	transactor *bind.TransactOpts
}

// NewArbitrageContract creates a new instance of the arbitrage contract wrapper
func NewArbitrageContract(address string, client *ethclient.Client, transactor *bind.TransactOpts) (*ArbitrageContract, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid contract address: %s", address)
	}

	contractABI, err := abi.JSON(strings.NewReader(arbitrageABI))
	if err != nil {
		return nil, fmt.Errorf("parsing ABI: %w", err)
	}

	return &ArbitrageContract{
		address:    common.HexToAddress(address),
		abi:        contractABI,
		client:     client,
		transactor: transactor,
	}, nil
}

// ExecuteArbitrage calls the contract's executeArbitrage method
func (c *ArbitrageContract) ExecuteArbitrage(
	ctx context.Context,
	borrowToken common.Address,
	borrowAmount *big.Int,
	intermediateToken common.Address,
	sourceDex uint8,
	slippage *big.Int,
) (*types.Transaction, error) {
	// Create the transaction
	tx, err := c.transactor.BindTransactor(
		&c.abi,
		c.address,
		"executeArbitrage",
		borrowToken,
		borrowAmount,
		intermediateToken,
		sourceDex,
		slippage,
	)
	if err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	return tx, nil
}

// CheckProfitability calls the contract's checkArbitrageProfitability method
func (c *ArbitrageContract) CheckProfitability(
	ctx context.Context,
	borrowToken common.Address,
	borrowAmount *big.Int,
	intermediateToken common.Address,
	sourceDex uint8,
) (bool, *big.Int, error) {
	var isProfitable bool
	var expectedProfit *big.Int

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	result, err := c.abi.Pack("checkArbitrageProfitability", borrowToken, borrowAmount, intermediateToken, sourceDex)
	if err != nil {
		return false, nil, fmt.Errorf("packing method args: %w", err)
	}

	res, err := c.client.CallContract(ctx, ethereum.CallMsg{
		To:   &c.address,
		Data: result,
	}, nil)
	if err != nil {
		return false, nil, fmt.Errorf("calling contract: %w", err)
	}

	err = c.abi.UnpackIntoInterface(&[]interface{}{&isProfitable, &expectedProfit}, "checkArbitrageProfitability", res)
	if err != nil {
		return false, nil, fmt.Errorf("unpacking result: %w", err)
	}

	return isProfitable, expectedProfit, nil
}

// WithdrawTokens calls the contract's withdrawTokens method
func (c *ArbitrageContract) WithdrawTokens(
	ctx context.Context,
	token common.Address,
) (*types.Transaction, error) {
	// Create the transaction
	tx, err := c.transactor.BindTransactor(
		&c.abi,
		c.address,
		"withdrawTokens",
		token,
	)
	if err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	return tx, nil
}

// ABI for the arbitrage contract (partial)
const arbitrageABI = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "borrowToken",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "borrowAmount",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "intermediateToken",
				"type": "address"
			},
			{
				"internalType": "uint8",
				"name": "sourceDex",
				"type": "uint8"
			},
			{
				"internalType": "uint256",
				"name": "slippage",
				"type": "uint256"
			}
		],
		"name": "executeArbitrage",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "borrowToken",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "borrowAmount",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "intermediateToken",
				"type": "address"
			},
			{
				"internalType": "uint8",
				"name": "sourceDex",
				"type": "uint8"
			}
		],
		"name": "checkArbitrageProfitability",
		"outputs": [
			{
				"internalType": "bool",
				"name": "isProfitable",
				"type": "bool"
			},
			{
				"internalType": "uint256",
				"name": "expectedProfit",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "token",
				"type": "address"
			}
		],
		"name": "withdrawTokens",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`
