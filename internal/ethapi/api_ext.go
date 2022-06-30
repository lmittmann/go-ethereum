package ethapi

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type ExecutionResultArgs struct {
	GasUsed     hexutil.Uint64
	MinGasLimit hexutil.Uint64
	Output      hexutil.Bytes
	AccessList  types.AccessList
	Logs        []*types.Log
	Err         error
}

func (s *TransactionAPI) Multicall(ctx context.Context, args []TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride) ([]ExecutionResultArgs, error) {
	results, err := DoMulticall(ctx, s.b, args, blockNrOrHash, overrides, s.b.RPCEVMTimeout(), s.b.RPCGasCap())
	if err != nil {
		return nil, err
	}

	var firstErr error
	for i, result := range results {
		if err := result.Err; err != nil {
			if reason, unpackErr := abi.UnpackRevert(result.Output); unpackErr == nil {
				result.Err = fmt.Errorf("%w: %s", err, reason)
			}
			firstErr = fmt.Errorf("call %d reverted: %w", i+1, result.Err)
		}
	}
	return results, firstErr
}

func DoMulticall(ctx context.Context, b Backend, args []TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, timeout time.Duration, globalGasCap uint64) ([]ExecutionResultArgs, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	if err := overrides.Apply(state); err != nil {
		return nil, err
	}
	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	var (
		evm     *vm.EVM
		vmError func() error
	)
	go func() {
		// Wait for the context to be done and cancel the evm. Even if the
		// EVM has finished, cancelling may be done (repeatedly)
		<-ctx.Done()
		if evm != nil {
			evm.Cancel()
		}
	}()

	results := make([]ExecutionResultArgs, len(args))
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	for i, arg := range args {
		msg, err := arg.ToMessage(math.MaxUint64, header.BaseFee)
		if err != nil {
			return nil, err
		}

		tracer := logger.NewAccessListTracer(nil, msg.From(), *msg.To(), vm.PrecompiledAddressesBerlin)
		evm, vmError, err = b.GetEVM(ctx, msg, state, header, &vm.Config{NoBaseFee: true, Tracer: tracer})
		if err != nil {
			return nil, err
		}

		// Execute the message.
		result, err := core.ApplyMessage(evm, msg, gp)
		if err := vmError(); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		results[i] = ExecutionResultArgs{
			GasUsed:     hexutil.Uint64(result.UsedGas),
			MinGasLimit: hexutil.Uint64(result.UsedGas + state.GetRefund()),
			Output:      result.ReturnData,
			AccessList:  tracer.AccessList(),
			Logs:        state.Logs(),
			Err:         result.Err,
		}
		if i < len(args)-1 {
			state.Commit(false)
		}
	}
	return results, nil
}
