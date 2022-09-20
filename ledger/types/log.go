// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"io"
	"math/big"

	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/rlp"
)

//go:generate gencodec -type Log -field-override logMarshaling -out gen_log_json.go

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`
}

type rlpLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}

// EncodeRLP implements rlp.Encoder.
func (l *Log) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpLog{Address: l.Address, Topics: l.Topics, Data: l.Data})
}

// DecodeRLP implements rlp.Decoder.
func (l *Log) DecodeRLP(s *rlp.Stream) error {
	var dec rlpLog
	err := s.Decode(&dec)
	if err == nil {
		l.Address, l.Topics, l.Data = dec.Address, dec.Topics, dec.Data
	}
	return err
}

// BalanceChange represents a contract balance transfer event.
type BalanceChange struct {
	// address of the account
	Address common.Address `json:"address"`
	// type of token changes. pando=0, ptx=1
	TokenType uint `json:"token_type"`
	// whether the delta is negative
	IsNegative bool `json:"is_negative"`
	// amount changed.
	Delta *big.Int `json:"delta"`
}

type rlpBalanceChange struct {
	Address    common.Address
	TokenType  uint
	IsNegative bool
	Delta      *big.Int
}

// EncodeRLP implements rlp.Encoder.
func (b *BalanceChange) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpBalanceChange{Address: b.Address, TokenType: b.TokenType, IsNegative: b.IsNegative, Delta: b.Delta})

}

// DecodeRLP implements rlp.Decoder.
func (b *BalanceChange) DecodeRLP(s *rlp.Stream) error {
	var dec rlpBalanceChange
	err := s.Decode(&dec)
	if err == nil {
		b.Address, b.TokenType, b.IsNegative, b.Delta = dec.Address, dec.TokenType, dec.IsNegative, dec.Delta
	}
	return err
}
