package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type WaykiCallContractTx struct {
	WaykiBaseSignTx
	AppId    *UserIdWraper //user regid or user key id or app regid
	Fees     uint64
	Values   uint64 //transfer amount
	Contract []byte
}

func (tx WaykiCallContractTx) SignTx() string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteUserId(tx.AppId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteBytes(tx.Contract)

	signedBytes := tx.doSignTx()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiCallContractTx) doSignTx() []byte {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteUserId(tx.AppId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteBytes(tx.Contract)

	hash, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(tx.PrivateKey)
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
