package coder

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// RaydiumAmmInstructionCoder implements the Coder interface.
type RaydiumAmmInstructionCoder struct{}

func NewRaydiumAmmInstructionCoder() *RaydiumAmmInstructionCoder {
	return &RaydiumAmmInstructionCoder{}
}

// Decode decodes the given byte array into an instruction.
func (coder *RaydiumAmmInstructionCoder) Decode(data []byte) (interface{}, error) {
	return decodeData(data)
}

func (coder *RaydiumAmmInstructionCoder) DecodeCompute(data []byte) (Compute, error) {
	return decodeCompute(data)
}

func (coder *RaydiumAmmInstructionCoder) DecodeTransfer(data []byte) (Transfer, error) {
	return decodeTransfer(data)
}

// Decoding function.
func decodeData(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)
	var instructionID byte
	binary.Read(buf, binary.LittleEndian, &instructionID)

	switch instructionID {
	case 1:
		return decodeInitialize2(buf)
	case 4:
		return decodeWithdraw(buf)
	case 9:
		return decodeSwapBaseIn(buf)
	case 11:
		return decodeSwapBaseOut(buf)

	default:
		return nil, errors.New("invalid instruction ID")
	}
}

func decodeCompute(data []byte) (Compute, error) {
	var instruction Compute

	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &instruction.Instruction)
	binary.Read(buf, binary.LittleEndian, &instruction.Value)

	return instruction, nil
}

func decodeTransfer(data []byte) (Transfer, error) {
	var instruction Transfer

	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &instruction.Instruction)
	binary.Read(buf, binary.LittleEndian, &instruction.Amount)
	return instruction, nil
}

func decodeInitialize2(buf *bytes.Reader) (Initialize2, error) {
	var instruction Initialize2
	binary.Read(buf, binary.LittleEndian, &instruction.Nonce)
	binary.Read(buf, binary.LittleEndian, &instruction.OpenTime)
	binary.Read(buf, binary.LittleEndian, &instruction.InitPcAmount)
	binary.Read(buf, binary.LittleEndian, &instruction.InitCoinAmount)

	return instruction, nil
}

func decodeWithdraw(buf *bytes.Reader) (Withdraw, error) {
	var instruction Withdraw
	binary.Read(buf, binary.LittleEndian, &instruction.Amount)

	return instruction, nil
}

func decodeSwapBaseIn(buf *bytes.Reader) (SwapBaseIn, error) {
	var instruction SwapBaseIn
	binary.Read(buf, binary.LittleEndian, &instruction.AmountIn)
	binary.Read(buf, binary.LittleEndian, &instruction.MinimumAmountOut)

	return instruction, nil
}

func decodeSwapBaseOut(buf *bytes.Reader) (SwapBaseOut, error) {
	var instruction SwapBaseOut
	binary.Read(buf, binary.LittleEndian, &instruction.MaxAmountIn)
	binary.Read(buf, binary.LittleEndian, &instruction.AmountOut)

	return instruction, nil
}
