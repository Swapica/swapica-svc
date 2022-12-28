package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"reflect"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func NewProxy(rpc string, signer signature.Signer, swapperContract string, confirmations int) (types.Proxy, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}

	b, err := swapica.NewSwapica(common.HexToAddress(swapperContract), client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create swapper contract for address %s", swapperContract)
	}

	return &evmProxy{
		client:          client,
		signer:          signer,
		chainID:         chainID,
		swapperContract: common.HexToAddress(swapperContract),
		swapper:         b,
		confirmations:   confirmations,
	}, nil
}

type evmProxy struct {
	client          *ethclient.Client
	signer          signature.Signer
	chainID         *big.Int
	swapperContract common.Address
	swapper         *swapica.Swapica
	confirmations   int
}

func (e *evmProxy) getThreshold() (int64, error) {
	threshold, err := e.swapper.SignaturesThreshold(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}

	return threshold.Int64(), nil
}

func (e *evmProxy) checkTxDataAndSign(opts *bind.TransactOpts, tx *ethTypes.Transaction, rawData []byte) (*ethTypes.Transaction, int64, error) {
	abi, err := swapica.SwapicaMetaData.GetAbi()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to parse swapica ABI")
	}

	newParams, newMethod, err := decodeTxParams(*abi, tx.Data())
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to decode tx params")
	}

	oldParams, oldMethod, err := decodeTxParams(*abi, rawData)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to decode tx params")
	}

	if len(oldParams) != len(newParams) {
		return nil, 0, types.ErrWrongSignedTx
	}

	if newMethod.Name != oldMethod.Name {
		return nil, 0, types.ErrWrongSignedTx
	}

	// Check if all params except signature is equal
	if !reflect.DeepEqual(oldParams[:len(oldParams)-1], newParams[:len(newParams)-1]) {
		return nil, 0, types.ErrWrongSignedTx
	}

	signatures := oldParams[len(oldParams)-1].([][]byte)
	newSig := newParams[len(newParams)-1].([][]byte)[0]

	for _, sig := range signatures {
		if reflect.DeepEqual(sig, newSig) {
			return nil, int64(len(signatures)), errors.New("double signature")
		}
	}

	// Add signature to params and encode new transaction
	newParams[len(newParams)-1] = append(signatures, newSig)
	if signs, ok := newParams[len(newParams)-1].([][]byte); ok {
		contract := bind.NewBoundContract(e.swapperContract, *abi, nil, e.client, nil)
		newTx, err := contract.Transact(opts, newMethod.Name, newParams...)
		return newTx, int64(len(signs)), err
	}

	return nil, 0, errors.New("failed to cast signatures param")
}
