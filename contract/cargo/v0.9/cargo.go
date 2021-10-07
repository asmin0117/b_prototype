package main

import(
	"encoding/json"
	"fmt"
	"strconv"
	"bytes"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {}

type PERSON struct { //수산물 보유한사람, 화물차 보유한 사람
	Name string `json:"name"`
	Account string `json:"account"` // key a/c0~
	Mileage string `json:"mileage"`
}


// TransContract 구조체
// state -> requested, registered, responded, comfirmed, prepaid, loaded, ondelivery, arrived, paid, compeleted
// TC0~
type TRANSCONTRACT struct {
	CTYPE  string `json:"ctype"`  //"S or D"
	Sender string `json:"sender"` // 대금지불자 = 화물요청자 acc
	Receiver string `json:"receiver"` // 대금수취자 = 화물운송자 acc

	// 공통
	TruckSize string `json:"trucksize"`
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	ExpectedPayment string `json:"expectedpayment"`
        // 화물운송자가 입력
	MaxCargoSize string `json:"maxcargosize"`

	// 화물요청자가 입력
	CargoSize string `json:"cargosize"`

	// 계약 상태
	ContractStatus string `json:"contractstatus"`
}


func(s * SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
"regPERSON", "sumin", "a/c111", "5000"

"regPERSON", "dr.choi", "a/c222", "5000"

"reqTransportation","TC0", "a/c111","5","Seoul","Busan","1000","10" 

"queryReqTransportation", "TC0"

"regTransportation", "TC1", "a/c222",  "5", "Seoul", "Busan", "1000", "10"

"queryRegTransportation", "TC1"

"respond", "TC0", "D", "a/c222"

"confirmContract", "TC0", "D", "a/c222" 

"load", "TC0"
"depart", "TC0",
"arrive", "TC0",
"pay", "TC0"

*/

func(s * SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function,
	args := APIstub.GetFunctionAndParameters()

	if function == "reqTransportation"{
		return s.reqTransportation(APIstub, args)
	} else if function == "regTransportation" {
		return s.regTransportation(APIstub, args)
	} else if function == "respond" {
		return s.respond(APIstub, args)
	} else if function == "confirmContract" {
		return s.confirmContract(APIstub, args)
	} else if function == "load" {
		return s.load(APIstub, args)
	} else if function == "depart" {
		return s.depart(APIstub, args)
	} else if function == "arrive" {
		return s.arrive(APIstub, args)
	} else if function == "pay" {
		return s.pay(APIstub, args)
	} else if function == "queryReqTransportation" {
		return s.queryReqTransportation(APIstub, args)
	} else if function == "queryRegTransportation" {
		return s.queryRegTransportation(APIstub, args)
	} else if function == "regPERSON" {
		return s.regPERSON(APIstub, args)
	} else if function == "history" {
		return s.history(APIstub, args)
	} else {
		return shim.Error("Invalid function name!!")
	}

	// return shim.Error("Invalid function name!!")
}


// person 등록하는 메서드 만들고
// 그러면 person의 키는 account가 되는 것 
func(s * SmartContract) regPERSON(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	// (TO DO) verify sender, expected mileage를 가지고 있는가?
	// (TO DO) args[0] -> key 가 유효한가? TCx 인가?

	fmt.Println("regPERSON start:", args[0], args[1], args[2])

	var ppp = PERSON {
		Name: 		args[0],
		Account: 	args[1], // key a/c0~
		Mileage: 	args[2],
	}
	
	pAsBytes, _ := json.Marshal(ppp)

	err := APIstub.PutState(args[1], pAsBytes)

	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[1]))
	}

	fmt.Println("regPERSON end:",  args[1])

	return shim.Success(nil)
}



// 화물주가 
// "TC0", "a/c111","5","Seoul","Busan","1000",'10"
func(s * SmartContract) reqTransportation(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	// (TO DO) verify sender, expected milage를 가지고 있는가?
	// (TO DO) args[0] -> key 가 유효한가? TCx 인가?

	var sender = args[1]

	var reqTransportInfo = TRANSCONTRACT {
		CTYPE: "S",
		Sender: sender, // 대금지불자 = 화물요청자
		
		TruckSize: args[2],
		Origin: args[3],
		Destination: args[4],
		ExpectedPayment: args[5],
		
		// 화물요청자가 입력
		CargoSize: args[6],
	
		// 계약 상태
		ContractStatus: "requested",
	}
	

	reqTransportationAsBytes, _ := json.Marshal(reqTransportInfo)
	err := APIstub.PutState(args[0], reqTransportationAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}

	return shim.Success(nil)
}


func(s * SmartContract) queryReqTransportation(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1") //query by key
	}

	reqTransportationAsBytes, _ := APIstub.GetState(args[0])
	if reqTransportationAsBytes == nil {
		return shim.Error("Could not query reqTransportation")
	}
	return shim.Success(reqTransportationAsBytes)
}

// 기사가 등록
// "regTransportation", "a/c222",  "5", "Seoul", "Busan", "1000", "10"

func(s * SmartContract) regTransportation(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// (TO DO) verify sender
	// (TO DO) args[0] -> key 가 유효한가? TCx 인가?

	var receiver = args[1]

	var regTransportInfo = TRANSCONTRACT {
		CTYPE: "D",
		Receiver: receiver,
		// 공통
		TruckSize: args[2],
		Origin: args[3],
		Destination: args[4],
		ExpectedPayment: args[5],
		
		// 화물요청자가 입력
		CargoSize: args[6],
	
		// 계약 상태
		ContractStatus: "registered",
	}

	regTransportationAsBytes, _ := json.Marshal(regTransportInfo)
	err := APIstub.PutState(args[0], regTransportationAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}

	return shim.Success(nil)
}

func(s * SmartContract) queryRegTransportation(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1") //query by key
	}

	regTransportationAsBytes, _ := APIstub.GetState(args[0])
	if regTransportationAsBytes == nil {
		return shim.Error("Could not query regTransportation")
	}
	return shim.Success(regTransportationAsBytes)
}

// 화물주나 기사가 응답
//"TC0", "D", "a/c222",

func(s * SmartContract) respond(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// GetState
	contractAsBytes, _ := APIstub.GetState(args[0])
	// (TO DO) contractAsBytes == nil? -> 등록된 TC가 없음

	// (TO DO) args[2] -> 기사가 유효한가? S or D 유효

	TransContract := TRANSCONTRACT{}
	
	json.Unmarshal(contractAsBytes, &TransContract)

	if args[1] == "D"{
		TransContract.Receiver = args[2]
	} else {
		TransContract.Sender = args[2]
		// (TO DO) milage가 expected payment보다 이상인가?
	}
		
	TransContract.ContractStatus = "responded"

	contractAsBytes, _ = json.Marshal(TransContract)
	err := APIstub.PutState(args[0], contractAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}
	return shim.Success(nil)
}
//"TC0", "D", "a/c221",
func(s * SmartContract) confirmContract(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// GetState
	contractAsBytes, _ := APIstub.GetState(args[0])
	// (TO DO) contractAsBytes == nil? -> 등록된 TC가 없음

	// (TO DO) args[2]  S or D 유효

	TransContract := TRANSCONTRACT{}
	
	json.Unmarshal(contractAsBytes, &TransContract)
		
	TransContract.ContractStatus = "confirmed"

	contractAsBytes, _ = json.Marshal(TransContract)
	err := APIstub.PutState(args[0], contractAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}
	return shim.Success(nil)
}

// "TC0"
func(s * SmartContract) load(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// GetState
	contractAsBytes, _ := APIstub.GetState(args[0])
	// (TO DO) contractAsBytes == nil? -> 등록된 TC가 없음

	// (TO DO) 계약이 완료되었느냐?

	TransContract := TRANSCONTRACT{}
	json.Unmarshal(contractAsBytes, &TransContract)

	TransContract.ContractStatus = "loaded"

	contractAsBytes, _ = json.Marshal(TransContract)
	err := APIstub.PutState(args[0], contractAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}
	return shim.Success(nil)
}

// 고쳐야 함, 미터 고쳐야함

func(s * SmartContract) depart(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// GetState
	contractAsBytes, _ := APIstub.GetState(args[0])
	// (TO DO) contractAsBytes == nil? -> 등록된 TC가 없음

	TransContract := TRANSCONTRACT{}
	json.Unmarshal(contractAsBytes, &TransContract)

	TransContract.ContractStatus = "ondelivery"

	contractAsBytes, _ = json.Marshal(TransContract)
	err := APIstub.PutState(args[0], contractAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}
	return shim.Success(nil)
}

// FIXME: 고쳐야 함, 미터 고쳐야함
func(s * SmartContract) arrive(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// GetState
	contractAsBytes, _ := APIstub.GetState(args[0])
	// (TO DO) contractAsBytes == nil? -> 등록된 TC가 없음

	TransContract := TRANSCONTRACT{}
	json.Unmarshal(contractAsBytes, &TransContract)

	TransContract.ContractStatus = "arrived"

	contractAsBytes, _ = json.Marshal(TransContract)
	err := APIstub.PutState(args[0], contractAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}
	return shim.Success(nil)
}

// "TC0"
func(s * SmartContract) pay(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	// sender와 receiver에 대한 정보를 받고, sender가 보내는 마일리지에 대해 receiver.mileage +=

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// GetState
	contractAsBytes, _ := APIstub.GetState(args[0])
	// (TO DO) contractAsBytes == nil? -> 등록된 TC가 없음

	// (TO DO) 운송이 완료되었는가

	TC := TRANSCONTRACT{}
	json.Unmarshal(contractAsBytes, &TC)
	sender_acc := TC.Sender
	reciever_acc := TC.Receiver
	transfer_amount := TC.ExpectedPayment

	 fmt.Println("pay log:", sender_acc, reciever_acc, transfer_amount)

	//TCS -> Marshal -> PutState PERSON
	//TCR -> Marshal -> PutState PERSON
	//--------------------------------------------------------------------------

	//TC.sender ->  GetState -> Unmarshal -> TCS
	senderAsBytes, _ := APIstub.GetState(sender_acc)
	sender := PERSON {}

	//TC.receiver -> GetState -> Unmarshal -> TCR
	receiverAsBytes, _ := APIstub.GetState(reciever_acc)
	receiver := PERSON {}

	// // (TO DO) senderAsBytes, receiverAsBytes nil 인지 확인

	json.Unmarshal(senderAsBytes, &sender)
	json.Unmarshal(receiverAsBytes, &receiver)

	// //input 거래금
	input, _ := strconv.Atoi(transfer_amount)
	
	//TCS.milage 빼줘 TC.ExpectedPayment
	//TCR.milage 더해줘 TC.ExpectedPayment
	// (TO DO) sender 마일리지가 input 보다 큰가?
	sender_mileage, _ := strconv.Atoi(sender.Mileage)
	sender.Mileage = strconv.Itoa(sender_mileage - input)

	receiver_mileage, _ := strconv.Atoi(receiver.Mileage)
	receiver.Mileage = strconv.Itoa(receiver_mileage + input)

	senderAsBytes, _ = json.Marshal(sender)
	receiverAsBytes, _ = json.Marshal(receiver)

	APIstub.PutState(sender_acc, senderAsBytes)
	APIstub.PutState(reciever_acc, receiverAsBytes)

	TC.ContractStatus = "paid"

	//TC -> Marshal -> PutState TRANSCONTRACT
	contractAsBytes, _ = json.Marshal(TC)
	err := APIstub.PutState(args[0], contractAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed!! : %s", args[0]))
	}
	return shim.Success(nil)
}

func (s * SmartContract) history(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	keyName := args[0]
	// 로그 남기기
	fmt.Println("readTxHistory:" + keyName)

	resultsIterator, err := stub.GetHistoryForKey(keyName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	// 로그 남기기
	fmt.Println("readTxHistory returning:\n" + buffer.String() + "\n")

	return shim.Success(buffer.Bytes())
}



func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
