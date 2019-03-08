package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type OperVoteFund struct {
	VoteType  WaykiVoteType
	PubKey    *PubKeyId
	VoteValue int64
}

type WaykiDelegateTx struct {
	WaykiBaseSignTx
	OperVoteFunds []OperVoteFund
	Fees          uint64
}

func (tx WaykiDelegateTx) SignTx() string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(len(tx.OperVoteFunds)))
	for _, fund := range tx.OperVoteFunds {
		writer.WriteByte(byte(fund.VoteType))
		writer.WritePubKeyId(*fund.PubKey)
		writer.WriteVarInt(fund.VoteValue)
	}
	writer.WriteVarInt(int64(tx.Fees))
	//signedBytes := tx.doSignTx()
	//writer.WriteBytes(signedBytes)
	writer.WriteBytes([]byte{})

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiDelegateTx) doSignTx() []byte {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(len(tx.OperVoteFunds)))
	for _, fund := range tx.OperVoteFunds {
		writer.WriteByte(byte(fund.VoteType))
		writer.WriteBytes(*fund.PubKey)
		writer.WriteVarInt(fund.VoteValue)
	}
	writer.WriteVarInt(int64(tx.Fees))
	hash, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(tx.PrivateKey)
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
