package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//SmartContract
type SmartContract struct {
}

//Value
type Value struct {
	SensorID     string `json:"sensorID"`
	Temp         string `json:"temp"`
	Time         string `json:"time"`
	Outlier      string `json:"outlier"`
	StuckatFault string `json:"stuckatfault"`
	SCFault      string `json:"scfault"`
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error when starting SmartContract", err)
	}
}

func sum(numbers []float64) (total float64) {
	for _, x := range numbers {
		total += x
	}
	return total
}

func stdDev(numbers []float64, mean float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += math.Pow(number-mean, 2)
	}
	variance := total / float64(len(numbers)-1)
	return math.Sqrt(variance)
}

func Abs(x float64) float64 {
	if x < 0 {
		return -1 * x
	}
	return x
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	values := []Value{
		Value{SensorID: "0", Temp: " ", Time: " ", Outlier: " "},
		Value{SensorID: "1", Temp: " ", Time: " ", Outlier: " "},
		Value{SensorID: "2", Temp: " ", Time: " ", Outlier: " "},
	}

	i := 0
	for i < len(values) {
		fmt.Println("i is ", i)
		valueAsBytes, _ := json.Marshal(values[i])
		stub.PutState(strconv.Itoa(i), valueAsBytes)
		fmt.Println("Added", values[i])
		i = i + 1
	}

	fmt.Println("Chaincode instantiated")
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "registerSensor" {
		return s.registerSensor(stub, args)
	} else if function == "addTemp" {
		return s.addTemp(stub, args)
	} else if function == "getHistory" {
		return s.getHistory(stub, args)
	} else if function == "getState" {
		return s.getState(stub, args)
	} else if function == "getAllStates" {
		return s.getAllStates(stub, args)
	}

	return shim.Success(nil)
}

func (s *SmartContract) detectStuckatFault(stub shim.ChaincodeStubInterface, sensorID string, measurement float64) string {
	type Response struct {
		TxId      string
		Value     Value
		Timestamp string
		IsDelete  string
	}

	var response []Response
	counter := 0
	// b := `[{"TxId":"0d6157c800fbe1904b63a7fe81a096f65b91d6c426da74c3299de9f9a70e1dbc", "Value":{"sensorID":"3","temp":" ","time":" ","outlier":""}, "Timestamp":"2020-04-14 22:57:16.602 +0000 UTC", "IsDelete":"false"},{"TxId":"b0c042787d8ea0083915c0e5c8cae48d09fb8b839bb6d9e5b13db9892d6e6824", "Value":{"sensorID":"3","temp":"21.13","time":"04/15/2020, 10:47:40 AM","outlier":""}, "Timestamp":"2020-04-15 07:47:40.386 +0000 UTC", "IsDelete":"false"},{"TxId":"2e37d9feb6460e9b400b410f11e678df5bad8bea8ada1c1ebe68c00854e692cc", "Value":{"sensorID":"3","temp":"20.12","time":"04/15/2020, 10:47:45 AM","outlier":""}, "Timestamp":"2020-04-15 07:47:45.282 +0000 UTC", "IsDelete":"false"}]`
	var temps []float64
	var sensorParam []string
	sensorParam = append(sensorParam, sensorID)
	b := s.getHistory(stub, sensorParam).Payload

	json.Unmarshal([]byte(b), &response)

	for _, element := range response {
		if (element.Value.Temp == "") || (element.Value.Temp == " ") {
			continue
		}
		if floatedTemp, err := strconv.ParseFloat(element.Value.Temp, 64); err == nil {
			temps = append(temps, floatedTemp)
		} else {
			fmt.Println(err)
		}
	}
	if len(temps) == 0 {
		return "0%"
	}

	for i := len(temps) - 1; i >= 0; i-- {
		if temps[i] == measurement {
			counter++
		} else {
			break
		}
		if counter == 10 {
			break
		}
	}

	return strconv.Itoa(counter*10) + "%"

}

func (s *SmartContract) detectSCFault(stub shim.ChaincodeStubInterface, sensorID string, measurement float64) string {
	type Response struct {
		Value Value
	}

	var response []Response
	counter := 0
	size := 0
	// b := `[{"TxId":"0d6157c800fbe1904b63a7fe81a096f65b91d6c426da74c3299de9f9a70e1dbc", "Value":{"sensorID":"3","temp":" ","time":" ","outlier":""}, "Timestamp":"2020-04-14 22:57:16.602 +0000 UTC", "IsDelete":"false"},{"TxId":"b0c042787d8ea0083915c0e5c8cae48d09fb8b839bb6d9e5b13db9892d6e6824", "Value":{"sensorID":"3","temp":"21.13","time":"04/15/2020, 10:47:40 AM","outlier":""}, "Timestamp":"2020-04-15 07:47:40.386 +0000 UTC", "IsDelete":"false"},{"TxId":"2e37d9feb6460e9b400b410f11e678df5bad8bea8ada1c1ebe68c00854e692cc", "Value":{"sensorID":"3","temp":"20.12","time":"04/15/2020, 10:47:45 AM","outlier":""}, "Timestamp":"2020-04-15 07:47:45.282 +0000 UTC", "IsDelete":"false"}]`
	var sensorParam []string
	sensorParam = append(sensorParam, "0")
	sensorParam = append(sensorParam, "999")

	b := s.getAllStates(stub, sensorParam).Payload

	json.Unmarshal([]byte(b), &response)

	if (response == nil) || (len(response) == 0) {
		return "-"
	}

	for _, element := range response {
		size++
		if element.Value.SensorID == sensorID {
			continue
		}
		if (element.Value.Temp == "") || (element.Value.Temp == " ") {
			continue
		}
		if floatedTemp, err := strconv.ParseFloat(element.Value.Temp, 64); err == nil {
			if Abs(measurement-floatedTemp) > float64(2) {
				counter++
			}
		} else {
			fmt.Println(err)
		}
	}
	if counter == 0 {
		return "0%"
	}

	if counter > 10 {
		return "100%"
	}

	return fmt.Sprintf("%.2f", (float64(counter)/(float64(size)-float64(1)))*float64(100)) + "%"

}

func (s *SmartContract) detectOutlier(stub shim.ChaincodeStubInterface, sensorID string, measurement float64) string {
	type Response struct {
		TxId      string
		Value     Value
		Timestamp string
		IsDelete  string
	}

	var response []Response
	// b := `[{"TxId":"0d6157c800fbe1904b63a7fe81a096f65b91d6c426da74c3299de9f9a70e1dbc", "Value":{"sensorID":"3","temp":" ","time":" ","outlier":""}, "Timestamp":"2020-04-14 22:57:16.602 +0000 UTC", "IsDelete":"false"},{"TxId":"b0c042787d8ea0083915c0e5c8cae48d09fb8b839bb6d9e5b13db9892d6e6824", "Value":{"sensorID":"3","temp":"21.13","time":"04/15/2020, 10:47:40 AM","outlier":""}, "Timestamp":"2020-04-15 07:47:40.386 +0000 UTC", "IsDelete":"false"},{"TxId":"2e37d9feb6460e9b400b410f11e678df5bad8bea8ada1c1ebe68c00854e692cc", "Value":{"sensorID":"3","temp":"20.12","time":"04/15/2020, 10:47:45 AM","outlier":""}, "Timestamp":"2020-04-15 07:47:45.282 +0000 UTC", "IsDelete":"false"}]`
	var temps []float64
	var sensorParam []string
	sensorParam = append(sensorParam, sensorID)
	b := s.getHistory(stub, sensorParam).Payload

	json.Unmarshal([]byte(b), &response)

	for _, element := range response {
		if (element.Value.Temp == "") || (element.Value.Temp == " ") {
			continue
		}
		if floatedTemp, err := strconv.ParseFloat(element.Value.Temp, 64); err == nil {
			temps = append(temps, floatedTemp)
		} else {
			fmt.Println(err)
		}
	}
	if len(temps) == 0 {
		return "No"
	}

	mean := sum(temps) / float64(len(temps))
	stdev := stdDev(temps, mean)

	// mean := stat.Mean(temps, nil)
	// stdev := stat.StdDev(temps, nil)
	UpperLimit := mean + stdev*3
	LowerLimit := mean - stdev*3

	fmt.Println(UpperLimit)
	fmt.Println(LowerLimit)

	if (measurement < LowerLimit) || (measurement > UpperLimit) {
		return "Yes"
	} else {
		return "No"
	}

}

func (s *SmartContract) registerSensor(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}
	sensorID := args[0]
	value := Value{SensorID: args[0], Temp: " ", Time: " "}
	valueAsBytes, _ := json.Marshal(value)
	stub.PutState(sensorID, valueAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) addTemp(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("invalid number of arguments")
	}
	sensorID := args[0]
	valueAsBytes, _ := stub.GetState(sensorID)
	value := Value{}
	json.Unmarshal(valueAsBytes, &value)
	fmt.Println(value)

	value.Temp = args[1]
	value.Time = args[2]
	if ms, err := strconv.ParseFloat(args[1], 64); err != nil {
		fmt.Println(err)
	} else {
		value.Outlier = s.detectOutlier(stub, sensorID, ms)
		value.StuckatFault = s.detectStuckatFault(stub, sensorID, ms)
		value.SCFault = s.detectSCFault(stub, sensorID, ms)
	}

	fmt.Println(value)
	valueByBytes, _ := json.Marshal(value)
	stub.PutState(sensorID, valueByBytes)

	return shim.Success(nil)
}

// GetState returns the value of the specified asset key
func (s *SmartContract) getState(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	sensorID := args[0]
	valueAsBytes, _ := stub.GetState(sensorID)
	value := Value{}
	json.Unmarshal(valueAsBytes, &value)
	buffer.WriteString(string(valueAsBytes))

	buffer.WriteString("]")

	fmt.Println(buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}
	sensorID := args[0]

	iterator, err := stub.GetHistoryForKey(sensorID)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer iterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			shim.Error(err.Error())
			fmt.Printf(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		if queryResponse.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(queryResponse.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(queryResponse.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("getHistory:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) getAllStates(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startID := args[0]
	endID := args[1]

	resultsIterator, err := stub.GetStateByRange(startID, endID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Value\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getAllStates queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
