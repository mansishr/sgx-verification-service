/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package parser

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"github.com/pkg/errors"
	"gopkg.in/restruct.v1"
	"math/big"
	"strconv"
)

type EcdsaQuoteData struct {
	PckCertSize uint32
}

type SwQuoteData struct {
	Dummy uint32
}

type SkcBlobHeader struct {
	MajorNum  uint32
	MinorNum  uint32
	QuoteSize uint32
	QuoteType uint32
	KeyType   uint32
}

type KeyDetailsRSA struct {
	ExponentLen uint32
	ModulusLen  uint32
}

type KeyDetailsEC struct {
	Dummy uint32
}

type SkcBlobParsed struct {
	Header         SkcBlobHeader
	RsaKeyDetails  KeyDetailsRSA
	ECKeyDetails   KeyDetailsEC
	EcdsaQuoteInfo EcdsaQuoteData
	SwQuoteInfo    SwQuoteData
	RawBlobLen     int
	RawBlob        []byte
	QuoteBlob      []byte
	PubKeyBlob     []byte
}

const (
	KeyTypeRsa = 1
	KeyTypeEc  = 2
)

const (
	QuoteTypeEcdsa = 1
	QuoteTypeSw    = 2
)

func ParseSkcQuoteBlob(rawBlob string) *SkcBlobParsed {
	if len(rawBlob) < 1 {
		log.Error("ParseSkcQuoteBlob: SKC Blob is Empty")
		return nil
	}

	parsedObj := new(SkcBlobParsed)
	_, err := parsedObj.parseSkcBlobData(rawBlob)
	if err != nil {
		log.Error("parseSkcBlobData: SKC Blob Parsing Error: ", err.Error())
		return nil
	}
	return parsedObj
}

func ParseQVLQuoteBlob(rawBlob string) *SkcBlobParsed {
	log.Trace("parser/skc_quote_parser:ParseQVLQuoteBlob() Entering")
	defer log.Trace("parser/skc_quote_parser:ParseQVLQuoteBlob() Leaving")
	if len(rawBlob) < 1 {
		log.Error("ParseQVLQuoteBlob: SKC Blob is Empty")
		return nil
	}

	parsedObj := new(SkcBlobParsed)
	decodedBlob, err := base64.StdEncoding.DecodeString(rawBlob)
	if err != nil {
		log.Error("Failed to Base64 Decode Quote")
		return nil
	}
	quoteSize := len(decodedBlob)
	parsedObj.QuoteBlob = make([]byte, quoteSize)
	copy(parsedObj.QuoteBlob, decodedBlob)
	return parsedObj
}

func (e *SkcBlobParsed) getKeyType() uint32 {
	return e.Header.KeyType
}

func (e *SkcBlobParsed) GetQuoteType() uint32 {
	return e.Header.QuoteType
}

func (e *SkcBlobParsed) GetQuoteBlob() []byte {
	return e.QuoteBlob
}

func (e *SkcBlobParsed) GetPubKeyBlob() []byte {
	return e.PubKeyBlob
}

func (e *SkcBlobParsed) parseSkcBlobData(blob string) (bool, error) {
	decodedBlob, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		log.Error("Failed to Base64 Decode Quote")
		return false, errors.Wrap(err, "ParseSkcBlob: Failed to Base64 Decode Quote")
	}

	var keyDetailsLen int
	var quoteDetailsLen int
	var pubKeySize int = 0

	e.RawBlob = make([]byte, len(decodedBlob))
	copy(e.RawBlob, decodedBlob)

	e.RawBlobLen = len(e.RawBlob)
	restruct.Unpack(e.RawBlob, binary.LittleEndian, &e.Header)

	if e.getKeyType() == KeyTypeRsa {
		restruct.Unpack(e.RawBlob[20:], binary.LittleEndian, &e.RsaKeyDetails)
		keyDetailsLen = 8
	} else if e.getKeyType() == KeyTypeEc {
		restruct.Unpack(e.RawBlob[20:], binary.LittleEndian, &e.ECKeyDetails)
		keyDetailsLen = 4
	} else {
		return false, errors.Wrap(err, "ParseSkcBlob: Invalid Key Type Received")
	}

	// first 20 bytes are added by skc_library to identify type of quote and signing key type
	quoteDetailsOffset := 20 + keyDetailsLen

	if e.GetQuoteType() == QuoteTypeEcdsa {
		restruct.Unpack(e.RawBlob[quoteDetailsOffset:], binary.LittleEndian, &e.EcdsaQuoteInfo)
	} else if e.GetQuoteType() == QuoteTypeSw {
		restruct.Unpack(e.RawBlob[quoteDetailsOffset:], binary.LittleEndian, &e.SwQuoteInfo)
	} else {
		return false, errors.Wrap(err, "parseSkcBlobData: Invalid Quote Type Received: ")
	}

	quoteDetailsLen = 4
	pubKeyStrOfset := quoteDetailsOffset + quoteDetailsLen
	quoteStrOffset := quoteDetailsOffset + quoteDetailsLen
	if e.GetQuoteType() == QuoteTypeEcdsa {
		quoteStrOffset = quoteStrOffset + int(e.EcdsaQuoteInfo.PckCertSize)
		pubKeyStrOfset = pubKeyStrOfset + int(e.EcdsaQuoteInfo.PckCertSize)
	}

	if e.getKeyType() == KeyTypeRsa {
		pubKeySize = (int(e.RsaKeyDetails.ModulusLen) + int(e.RsaKeyDetails.ExponentLen))
		quoteStrOffset = quoteStrOffset + pubKeySize
	} else {
		quoteStrOffset += 8 //Because of union member
	}

	pubKeyEndOffset := pubKeyStrOfset + pubKeySize
	e.PubKeyBlob = make([]byte, pubKeySize)
	copy(e.PubKeyBlob, e.RawBlob[pubKeyStrOfset:pubKeyEndOffset])

	quoteEndOffset := quoteStrOffset + int(e.Header.QuoteSize)
	log.Debug("TotalBlobSize: ", e.RawBlobLen, ", QuoteStrOffset: ", quoteStrOffset, ", QuoteEndOffet: ", quoteEndOffset)

	e.QuoteBlob = make([]byte, e.Header.QuoteSize)
	copy(e.QuoteBlob, e.RawBlob[quoteStrOffset:])

	e.DumpSkcBlobHeader()
	return true, nil
}

func (e *SkcBlobParsed) DumpSkcBlobHeader() {
	log.Debugf("===================================>SkcQuoteBlobHeader<=======================================")
	log.Debug("Header->MajorNum = ", e.Header.MajorNum)
	log.Debug("Header->MinorNum = ", e.Header.MinorNum)
	log.Debug("Header->QuoteSize = ", e.Header.QuoteSize)
	log.Debug("Header->QuoteType = ", e.Header.QuoteType)
	log.Debug("Header->KeyType = ", e.Header.KeyType)

	if e.getKeyType() == KeyTypeRsa {
		log.Debug("RSAKeyDetails->ModulusLen = ", e.RsaKeyDetails.ModulusLen)
		log.Debug("RSAKeyDetails->ExponentLen = ", e.RsaKeyDetails.ExponentLen)
	} else {
		log.Debug("ECKeyDetails->ModulusLen = ", e.ECKeyDetails.Dummy)
	}

	if e.GetQuoteType() == QuoteTypeEcdsa {
		log.Debug("EcdsaQuoteInfo->PckCertSize = ", e.EcdsaQuoteInfo.PckCertSize)
	} else if e.GetQuoteType() == QuoteTypeSw {
		log.Debug("SwQuoteInfo->Dummy = ", e.SwQuoteInfo.Dummy)
	}
}

func (e *SkcBlobParsed) GetRsaExponentLen() uint32 {
	if e.getKeyType() == KeyTypeRsa {
		return e.RsaKeyDetails.ExponentLen
	}
	log.Error("GetRSAModulusLen: Invalid Key type")
	return 0
}

func (e *SkcBlobParsed) GetRsaPubKey() ([]byte, error) {
	var err error
	if e.getKeyType() != KeyTypeRsa {
		return nil, errors.Wrap(err, "GetRsaPubKey: Invalid Public Key Type")
	}

	pubKeyBlob := e.GetPubKeyBlob()
	if len(pubKeyBlob) == 0 {
		return nil, errors.Wrap(err, "GetRsaPubKey: Invalid Public Key length")
	}

	exponentLen := int(e.GetRsaExponentLen())
	exponentArr := pubKeyBlob[:exponentLen]
	modulusStrOffset := exponentLen

	n := big.Int{}
	n.SetBytes(pubKeyBlob[modulusStrOffset:])
	eb := big.Int{}
	eb.SetBytes(exponentArr)

	ex, err := strconv.Atoi(eb.String())
	if err != nil {
		return nil, errors.Wrap(err, "GetRsaPubKey: Strconv to int")
	}

	pubKey := rsa.PublicKey{N: &n, E: int(ex)}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "GetRsaPubKey: Marshal error")
	}

	rsaPem := pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}
	rsaBytes := pem.EncodeToMemory(&rsaPem)
	if rsaBytes == nil {
		return nil, errors.Wrap(err, "GetRsaPubKey: Pem Encode failed")
	}

	return rsaBytes, nil
}
