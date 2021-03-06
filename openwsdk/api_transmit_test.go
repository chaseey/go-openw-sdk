package openwsdk

import (
	"encoding/json"
	"github.com/blocktree/openwallet/log"
	"github.com/google/uuid"
	"testing"
	"time"
)

func testServeTransmitNode(f func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo)) {

	api := testNewAPINode()
	err := api.ServeTransmitNode("127.0.0.1:9088")
	if err != nil {
		log.Errorf("ServeTransmitNode error: %v\n", err)
		return
	}

	tn, err := api.TransmitNode()
	if err != nil {
		log.Errorf("TransmitNode error: %v\n", err)
		return
	}

	tn.SetConnectHandler(f)

	time.Sleep(15 * time.Second)
}

func TestAPINode_ServeTransmitNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {
		log.Infof("nodeInfo: %v", nodeInfo)
	})
}

func TestTransmitNode_GetTrustNodeInfo(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		transmitNode.GetTrustNodeInfo(nodeInfo.NodeID, true,
			func(status uint64, msg string, nodeInfo *TrustNodeInfo) {
				log.Infof("nodeInfo: %v", nodeInfo)
			})
	})
}

func TestTransmitNode_CreateWalletViaTrustNode(t *testing.T) {

	alias := "openwallet333"
	password := "12345678"

	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {
		//创建钱包
		transmitNode.CreateWalletViaTrustNode(nodeInfo.NodeID, alias, password, true,
			func(status uint64, msg string, wallet *Wallet) {
				if wallet != nil {
					log.Infof("wallet: %+v\n", wallet)
				}
			})
	})
}

func TestTransmitNode_CreateAccountViaTrustNode(t *testing.T) {

	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		walletID := "WAJ7dKuES2LxJzwtJUWetsfqckiCCDQ4uy"
		alias := "openwallet_LTC_3"
		password := "12345678"
		symbol := "LTC"
		//创建账户
		transmitNode.CreateAccountViaTrustNode(nodeInfo.NodeID, walletID, alias, password, symbol, true,
			func(status uint64, msg string, account *Account, addresses []*Address) {
				if account != nil {
					log.Infof("account: %+v\n", account)
					for i, a := range addresses {
						log.Infof("address[%d]:%+v", i, a)
					}
				}
			})
	})
}

func TestTransmitNode_SendTransactionViaTrustNode(t *testing.T) {

	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		//accountID := "3i26MQmtuWVVnw8GnRCVopG3pi8MaYU6RqWVV2E1hwJx"
		//address := "mgCzMJDyJoqa6XE3RSdNGvD5Bi5VTWudRq"

		accountID := "3i26MQmtuWVVnw8GnRCVopG3pi8MaYU6RqWVV2E1hwJx"
		address := "mgCzMJDyJoqa6XE3RSdNGvD5Bi5VTWudRq"

		password := "12345678"
		sid := uuid.New().String()
		log.Infof("sid: %s", sid)
		transmitNode.SendTransactionViaTrustNode(nodeInfo.NodeID, accountID, password, sid,
			"", "0.03", address, "", "",
			true, func(status uint64, msg string, successTx []*Transaction, failedRawTxs []*FailedRawTransaction) {
				log.Infof("status: %d, msg: %s", status, msg)
				log.Info("============== success ==============")

				for _, tx := range successTx {
					log.Infof("tx: %+v", tx)
				}

				log.Info("")

				log.Info("============== fail ==============")

				for _, tx := range failedRawTxs {
					log.Infof("tx: %+v", tx.Reason)
				}

			})
	})
}

func TestTransmitNode_SetSummaryInfoViaTrustNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		setting := &SummarySetting{
			"WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA",
			"3i26MQmtuWVVnw8GnRCVopG3pi8MaYU6RqWVV2E1hwJx",
			"mkdStRouBPVrDVpYmbE5VUJqhBgxJb3dSS",
			"30",
			"0.001",
			"0",
			1,
		}

		transmitNode.SetSummaryInfoViaTrustNode(nodeInfo.NodeID, setting, true, func(status uint64, msg string) {
			log.Infof("msg:%+v", msg)
		})
	})
}

func TestTransmitNode_FindSummaryInfoByWalletIDViaTrustNode(t *testing.T) {

	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		walletID := "WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA"

		transmitNode.FindSummaryInfoByWalletIDViaTrustNode(nodeInfo.NodeID, walletID,
			true, func(status uint64, msg string, summarySettings []*SummarySetting) {
				for i, value := range summarySettings {
					log.Infof("SummarySetting[%d]: %+v", i, value)
				}
			})
	})
}

func TestTransmitNode_StartSummaryTaskViaTrustNode(t *testing.T) {

	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		plain := `

{
	"wallets": [{
		"walletID": "WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA",
		"password": "12345678",
		"accounts": [
		{
			"accountID": "A3Mxhqm65kTgS2ybHLenNrZzZNtLGVobDFYdpc1ge4eK",
			"feeRate": "0.001"
		},
		{
			"accountID": "3i26MQmtuWVVnw8GnRCVopG3pi8MaYU6RqWVV2E1hwJx",
			"feeRate": "0.001"
		}
		]
	}]
}

`
		var summaryTask SummaryTask
		err := json.Unmarshal([]byte(plain), &summaryTask)
		if err != nil {
			log.Error("json.Unmarshal error:", err)
			return
		}

		transmitNode.StartSummaryTaskViaTrustNode(nodeInfo.NodeID, 10, &summaryTask, SummaryTaskOperateTypeReset,
			true, func(status uint64, msg string) {
				log.Infof("status: %d, msg: %+v", status, msg)
			})
	})
}

func TestTransmitNode_StopSummaryTaskViaTrustNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		transmitNode.StopSummaryTaskViaTrustNode(nodeInfo.NodeID, true, func(status uint64, msg string) {
			log.Infof("msg:%+v", msg)
		})
	})
}

func TestTransmitNode_UpdateInfoViaTrustNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		transmitNode.UpdateInfoViaTrustNode(nodeInfo.NodeID, true, func(status uint64, msg string) {
			log.Infof("msg:%+v", msg)
		})
	})
}

func TestTransmitNode_AppendSummaryTaskViaTrustNode(t *testing.T) {

	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		setting := &SummarySetting{
			"WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA",
			"7ww2Gpfy8pN6HTngbMFBTEMAaVRGEpkmsiNkgAgqGQGf",
			"0x4f544cbd23c42950a5fe7f967c3e6938955a1718",
			"1",
			"0.01",
			"0",
			1,
		}

		transmitNode.SetSummaryInfoViaTrustNode(nodeInfo.NodeID, setting, true, func(status uint64, msg string) {
			log.Infof("msg:%+v", msg)
		})

		plain := `

{
	"wallets": [{
		"walletID": "WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA",
		"password": "12345678",
		"accounts": [{
			"accountID": "7ww2Gpfy8pN6HTngbMFBTEMAaVRGEpkmsiNkgAgqGQGf"
		}]
	}]
}

`
		var summaryTask SummaryTask
		err := json.Unmarshal([]byte(plain), &summaryTask)
		if err != nil {
			log.Error("json.Unmarshal error:", err)
			return
		}

		transmitNode.AppendSummaryTaskViaTrustNode(nodeInfo.NodeID, &summaryTask,
			true, func(status uint64, msg string) {
				log.Infof("msg:%+v", msg)
			})
	})
}

func TestTransmitNode_RemoveSummaryTaskViaTrustNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		transmitNode.RemoveSummaryTaskViaTrustNode(nodeInfo.NodeID,
			"WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA",
			"A3Mxhqm65kTgS2ybHLenNrZzZNtLGVobDFYdpc1ge4eK",
			true, func(status uint64, msg string) {
				log.Infof("msg:%+v", msg)
			})
	})
}

func TestTransmitNode_GetCurrentSummaryTaskViaTrustNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		transmitNode.GetCurrentSummaryTaskViaTrustNode(nodeInfo.NodeID,
			true, func(status uint64, msg string, task *SummaryTask) {
				log.Infof("msg:%+v", msg)
				for _, w := range task.Wallets {
					log.Infof("task wallet:%+v", w.WalletID)
					for _, a := range w.Accounts {
						log.Infof("task account:%+v", a.AccountID)
					}
				}

			})
	})
}

func TestTransmitNode_GetSummaryTaskLogViaTrustNode(t *testing.T) {
	testServeTransmitNode(func(transmitNode *TransmitNode, nodeInfo *TrustNodeInfo) {

		transmitNode.GetSummaryTaskLogViaTrustNode(nodeInfo.NodeID, 0, 200,
			true, func(status uint64, msg string, taskLog []*SummaryTaskLog) {
				log.Infof("msg:%+v", msg)
				for _, r := range taskLog {
					log.Infof("taskLog: %+v", r)
				}

			})
	})
}

