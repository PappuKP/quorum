// Copyright 2017 The go-ethereum Authors
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

package params

import (
	"math/big"
	"reflect"
	"testing"
)

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new *ChainConfig
		head        uint64
		wantErr     *ConfigCompatError
	}
	var storedMaxCodeConfig0, storedMaxCodeConfig1 []MaxCodeConfigStruct
	defaultRec := MaxCodeConfigStruct{big.NewInt(0), 24}
	rec1 := MaxCodeConfigStruct{big.NewInt(5), 32}
	rec2 := MaxCodeConfigStruct{big.NewInt(10), 40}
	rec3 := MaxCodeConfigStruct{big.NewInt(8), 40}

	storedMaxCodeConfig0 = append(storedMaxCodeConfig0, defaultRec)
	storedMaxCodeConfig1 = append(storedMaxCodeConfig1, defaultRec)
	storedMaxCodeConfig1 = append(storedMaxCodeConfig1, rec1)
	storedMaxCodeConfig1 = append(storedMaxCodeConfig1, rec2)

	var passedValidMaxConfig0 []MaxCodeConfigStruct
	passedValidMaxConfig0 = append(passedValidMaxConfig0, defaultRec)
	passedValidMaxConfig0 = append(passedValidMaxConfig0, rec1)

	var passedValidMaxConfig1 []MaxCodeConfigStruct
	passedValidMaxConfig1 = append(passedValidMaxConfig1, defaultRec)
	passedValidMaxConfig1 = append(passedValidMaxConfig1, rec1)
	passedValidMaxConfig1 = append(passedValidMaxConfig1, rec3)

	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &ChainConfig{EIP150Block: big.NewInt(10)},
			new:     &ChainConfig{EIP150Block: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{HomesteadBlock: nil},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{HomesteadBlock: big.NewInt(1)},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    big.NewInt(1),
				RewindTo:     0,
			},
		},
		{
			stored: &ChainConfig{HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:    &ChainConfig{HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			head:   25,
			wantErr: &ConfigCompatError{
				What:         "EIP150 fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored:  &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(10)}},
			new:     &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(20)}},
			head:    4,
			wantErr: nil,
		},
		{
			stored: &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(10)}},
			new:    &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(20)}},
			head:   30,
			wantErr: &ConfigCompatError{
				What:         "Ceil 2N/3 fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored: &ChainConfig{MaxCodeSizeChangeBlock: big.NewInt(10)},
			new:    &ChainConfig{MaxCodeSizeChangeBlock: big.NewInt(20)},
			head:   30,
			wantErr: &ConfigCompatError{
				What:         "max code size change fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored:  &ChainConfig{MaxCodeSizeChangeBlock: big.NewInt(10)},
			new:     &ChainConfig{MaxCodeSizeChangeBlock: big.NewInt(20)},
			head:    4,
			wantErr: nil,
		},
		{
			stored: &ChainConfig{QIP714Block: big.NewInt(10)},
			new:    &ChainConfig{QIP714Block: big.NewInt(20)},
			head:   30,
			wantErr: &ConfigCompatError{
				What:         "permissions fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored:  &ChainConfig{QIP714Block: big.NewInt(10)},
			new:     &ChainConfig{QIP714Block: big.NewInt(20)},
			head:    4,
			wantErr: nil,
		},
		{
			stored: &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig0},
			new:    &ChainConfig{MaxCodeSizeConfig: nil},
			head:   4,
			wantErr: &ConfigCompatError{
				What:         "genesis file missing max code size information",
				StoredConfig: big.NewInt(4),
				NewConfig:    big.NewInt(4),
				RewindTo:     3,
			},
		},
		{
			stored:  &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig0},
			new:     &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig0},
			head:    4,
			wantErr: nil,
		},
		{
			stored: &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig0},
			new:    &ChainConfig{MaxCodeSizeConfig: passedValidMaxConfig0},
			head:   10,
			wantErr: &ConfigCompatError{
				What:         "maxCodeSize data incompatible. updating maxCodeSize for past",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(10),
				RewindTo:     9,
			},
		},
		{
			stored:  &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig0},
			new:     &ChainConfig{MaxCodeSizeConfig: passedValidMaxConfig0},
			head:    4,
			wantErr: nil,
		},
		{
			stored:  &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig1},
			new:     &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig1},
			head:    12,
			wantErr: nil,
		},
		{
			stored:  &ChainConfig{MaxCodeSizeConfig: storedMaxCodeConfig1},
			new:     &ChainConfig{MaxCodeSizeConfig: passedValidMaxConfig1},
			head:    12,
			wantErr: &ConfigCompatError{
				What:         "maxCodeSize data incompatible. maxCodeSize historical data does not match",
				StoredConfig: big.NewInt(12),
				NewConfig:    big.NewInt(12),
				RewindTo:     11,
			},
		},
	}

	for _, test := range tests {
		err := test.stored.CheckCompatible(test.new, test.head, false)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}
