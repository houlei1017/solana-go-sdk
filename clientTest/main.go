package main

import (
	"bufio"
	"container/list"
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"io"
	"os"
	"strconv"
	"strings"
)

const rpcUrl = "http://122.248.192.182:8899"

var arg = map[string]interface{}{
	"encoding": "json",
	//"transactionDetails": "full",
	"transactionDetails": "signatures",
	"rewards":            false,
	"commitment":         "finalized",
}

// 两个线程分别从节点获取块数据，对比数据是否一致，验证空块情况
func getAndSaveFile() {

	go func() {
		curSlot := uint64(82939961)

		filePath := "nodelay.txt"
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		defer file.Close()

		client := client.NewClient(rpcUrl)
		for {
			slot, err := client.GetSlot(context.Background())
			if err != nil {
				fmt.Errorf("GetSlot", err)
				break
			}

			for i := curSlot - 1; i < slot; i++ {
				block, err := client.GetBlock(context.Background(), i)
				if err != nil {
					fmt.Errorf("GetConfirmedBlock", err)
					break
				}
				var str string
				if block.Blockhash == "" {
					str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", i, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
				} else {
					str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d\n", i, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
				}
				fmt.Print("no ->" + str)
				file.WriteString(str)
			}
		}
	}()

	curSlot := uint64(82939961)
	filePath := "somedelay.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()

	client := client.NewClient(rpcUrl)
	for {
		slot, err := client.GetSlot(context.Background())
		if err != nil {
			fmt.Errorf("GetSlot", err)
			break
		}
		slot -= 20

		for i := curSlot - 1; i < slot; i++ {
			block, err := client.GetBlock(context.Background(), i)
			if err != nil {
				fmt.Errorf("GetConfirmedBlock", err)
				break
			}
			var str string
			if block.Blockhash == "" {
				str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", i, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
			} else {
				str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d\n", i, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
			}
			fmt.Print("delay ->" + str)
			file.WriteString(str)
		}
	}
}

func getString(str string) string {
	result := strings.Index(str, "EMPTY slot")
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)

		var r = []rune(str)
		length := len(r)
		return string(r[result+11 : length-2])
	}
	return ""
}

func getListFromFile(filePath string) *list.List {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//读原来文件的内容，并且显示在终端
	l := list.New()
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		//fmt.Print(str)

		num := getString(str)
		//fmt.Println(num)
		l.PushBack(num)
	}
	return l
}

func checkEmptyBlock() {
	l := getListFromFile("empty.txt")
	client := client.NewClient(rpcUrl)

	for i := l.Front(); i != nil; i = i.Next() {
		//fmt.Println(i.Value)

		intNum, _ := strconv.Atoi(i.Value.(string))
		height := uint64(intNum)

		block, err := client.GetBlock(context.Background(), height)
		if err != nil {
			fmt.Println("GetConfirmedBlock", err)
			continue
		}
		var str string
		if block.Blockhash == "" {
			str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", height, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
		} else {
			str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d\n", height, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
		}
		fmt.Print(str)
	}
	fmt.Println("Total size: ", l.Len())
}

func checkBlock() {
	client11 := client.NewClient(rpcUrl)
	for {
		slot, err := client11.GetSlot(context.Background())
		if err != nil {
			fmt.Errorf("GetSlot", err)
			break
		}

		block, err := client11.GetBlock(context.Background(), slot)
		if err != nil || block.Blockhash == "" {
			if rpcErr, ok := err.(*rpc.JsonRpcError); ok {
				//Slot 被跳过打印日志记录，按正常逻辑处理
				//目前发现这两种是可以被忽略的情况,
				if rpcErr.Code == -32007 || rpcErr.Code == -32009 {
					fmt.Println("GetConfirmedBlock", err)
				}
			}
		}

		var str string
		if block.Blockhash == "" {
			str = fmt.Sprintf("slot:%d hash:%s ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", slot, block.Blockhash, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
		} else {
			str = fmt.Sprintf("slot:%d hash:%s ParentSLot:%d PreviousHash:%s size:%d\n", slot, block.Blockhash, block.ParentSlot, block.PreviousBlockhash, len(block.Transactions))
		}
		fmt.Print("no ->" + str)
	}
}

type Date111 struct {
	Parsed struct {
		Info struct {
			IsNative    bool   `json:"isNative"`
			Mint        string `json:"mint"`
			Owner       string `json:"owner"`
			State       string `json:"state"`
			TokenAmount struct {
				Amount         string  `json:"amount"`
				Decimals       int     `json:"decimals"`
				UiAmount       float64 `json:"uiAmount"`
				UiAmountString string  `json:"uiAmountString"`
			} `json:"tokenAmount"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Program string `json:"program"`
	Space   int    `json:"space"`
}

func tet11() {
	//rawurl := "https://api.devnet.solana.com"
	////rawurl := "https://api.testnet.solana.com"
	//rawurl := "https://api.mainnet-beta.solana.com"
	//
	//toPublicKey := "APhyMCpYjQ9RdEBn8cs4ifyBXjxAS5JtM3wYpWMJjsY5"
	////toPublicKey := "2bj53paPfbLXFBruju2XHEfdrdfQjdD1d1iVwAKyGRCS"
	//cs := client.NewClient(rawurl)
	////充值 验证地址必须为系统地址才可以
	//info, _ := cs.GetAccountInfo(context.Background(), toPublicKey)
	//tt, ok := info.(client.AccountInfo)
	//if ok {
	//	fmt.Println(tt.Parsed.Info.Owner, tt.Parsed.Info.Mint)
	//} else {
	//	fmt.Println("Fuck")
	//}
	//
	//resByre, resByteErr := json.Marshal(info.Data)
	//if resByteErr != nil {
	//	fmt.Println("读取信息失败")
	//	return
	//}
	//var newData Date111
	//jsonRes := json.Unmarshal(resByre, &newData)
	//if jsonRes != nil {
	//	fmt.Println("读取信息失败")
	//	return
	//}
	//
	////fmt.Println(newData)
	//fmt.Println(newData.Parsed.Info.Owner, newData.Parsed.Info.Mint)
}

func TestGetTokenAccountsByOwner() {
	rawurl := "https://api.mainnet-beta.solana.com"
	cs := client.NewClient(rawurl)

	//list := []string{
	//	"GQp5ZoNoNHNJq7ZvabygqV6513uY8WP7SL5Q3XVZ326T", //两个token 帐户
	//	"GA7HWfCEZ2GrQk1N3ceGV83oUz2KUDeDaXLsGVsz18aM", //一个token	帐户
	//	"2K3RjUsnNGr3awtRGek3USe4gAZ2E6fSYdz6t88ii2Cm"} //没有token帐户
	//mint := "kinXdEcpDQeHPEuQnqmUgtYykqKGVFq6CeVX5iAHJq6"
	list := []string{
		"48AEer19GHJohjnnDKMVfn2y6MBJGtxwvASkTeoxEaJC"} //没有token帐户
	mint := "7atgF8KQo4wJrD5ATGX7t1V2zVvykPJbFfNeVf1icFv1"
	for _, account := range list {
		info, err := cs.GetTokenAccountsByOwnerByMint(context.Background(), account, mint)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var chainAmount uint64 = 0
		for _, token := range info {
			fmt.Printf("Token account: %s  balance: %d \n", token.PublicKey.String(), token.Amount)
			chainAmount += token.Amount
		}

		fmt.Println(chainAmount)
	}
}

func getHealth() {
	list := []string{
		"https://api.mainnet-beta.solana.com",
		"http://122.248.192.182:8899",
		"http://18.140.48.191:8899",
	}

	for _, rawurl := range list {
		cs := client.NewClient(rawurl)
		info, err := cs.GetHealth(context.Background())
		if err != nil {
			fmt.Println(rawurl, err)
		} else {
			fmt.Println(rawurl, info)
		}
	}
}

func main() {
	//tet11()
	TestGetTokenAccountsByOwner()
	//getHealth()

}
