/**
  * Hyperledger Fabric implementation of the SRAM-based PUF Authentication and Integrity (SPAI) protocol
  * Code adapted from Udemy Course: Design and develop Fabric 2.1 applications from end-to-end using GoLang & Fabric Node SDK
  * presented by Rajeev Sakhuja.
 **/
 package main

 import (
	"fmt"
	// "io/ioutil"
	"math/rand"
	// goBytes "bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"
	"strings"

	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"

	peer "github.com/hyperledger/fabric-protos-go/peer"

	// Conversion functions
	"strconv"

	// JSON Encoding
	"encoding/json"
)

// SPAIChaincode Represents our chaincode object
type SPAIChaincode struct {
}

// Setup struct to return field sensor setup values
type Setup struct {
	Challenge string `json:"challenge"`
	Hash      string `json:"hash"`
	TimeStamp string `json:"ts"`
}

// Request struct to keep track of active requests
type Request struct {
	Status        string  `json:"status"`
	Crp           CRP     `json:"crp"`
	TimeStamp     uint32  `json:"timestamp"`
}

// CRP tmp storage
type CRP struct {
	Challenge string `json:"challenge"`
	Response  string `json:"response"`
}

// Sensor data struct of field sensor data
type Sensor struct {
	SensorID      string  `json:"sensorID"`
	DeviceID      string  `json:"rtuID"`
	Data          string  `json:"data"`
	Status        string  `json:"status"`
	Crp           []CRP   `json:"crp"`
	ActRequest    Request `json:"request"`
}



type Dependency struct {
	PrvEvent string `json:"prvEvent"`
	NxtEvent string `json:"nxtEvent"`
}

type Device struct {
	DeviceID string   `json:"deviceID"`
	State    string   `json:"state"`
	PrvState string   `json:"prvstate"`
	Message  []string `json:"message"`
	Event    []string `json:"event"`
}

type Process struct {
	ProcessID  string       `json:"processID"`
	Device     []Device     `json:"device"`
	Message    []string     `json:"message"`
	Dependency []Dependency `json:"dependency"`
}

// OwnerPrefix is used for creating the key for balances
// const OwnerPrefix="owner."

// Init Implements the Init method, dummy init method
func (token *SPAIChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// dummy init method
	fmt.Println("Init method")
	return successResponse("Init Successful!!!")
}

// Invoke method
func (token *SPAIChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Get the function name and parameters
	function, args := stub.GetFunctionAndParameters()

	fmt.Println("Invoke executed : ", function, " args=", args)

	switch {
	// Query function
	case function == "statusOf":
		return statusOf(stub, args)
	case function == "enroll":
		return enroll(stub, args)
	case function == "addCRP":
		return addCRP(stub, args)
	case function == "auth":
		return authSetup(stub, args)
	case function == "verify":
		return verify(stub, args)
	}

	return errorResponse("Invalid function",1)
}

func enroll(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	sProcess := args[0]
	// fmt.Println(sProcess)

	// Initialize process based on config file
	var process Process
	// Store in local struct
	json.Unmarshal([]byte(sProcess), &process)
	// Encode to JSON
	processJSON, _ := json.Marshal(process)
	// Save config to stateDB
	stub.PutState(process.ProcessID, []byte(processJSON))
	return shim.Success([]byte(processJSON))
}

func statusOf(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Check if owner id is in the arguments
	if len(args) < 1   {
		return errorResponse("Needs processID!!!", 6)
	}
	processID := args[0]
	bytes, err := stub.GetState(processID)
	if err != nil {
		return errorResponse(err.Error(), 7)
	}
	
	response := processJSON(processID, string(bytes))
	
	return successResponse(response)
}

func addCRP(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 3 {
		return errorResponse("Needs challenge and response!!!", 700)
	}

	sensorID := args[0]
	challenge := args[1]
	response := args[2]

	// get sensor status
	bytes, err := stub.GetState(sensorID)
	if err != nil {
		return errorResponse(err.Error(), 7)
	}

	var sensor Sensor
	err = json.Unmarshal(bytes, &sensor)
	if err != nil {
		return errorResponse("Failed to read sensor status", 700)
	}

	crp := CRP{Challenge: challenge, Response: response}
	sensor.Crp = append(sensor.Crp, crp)

	//convert to JSON and store the sensor struct
	sensorJSON, _ := json.Marshal(sensor)
	stub.PutState(sensorID, []byte(sensorJSON))
	
	return shim.Success([]byte(sensorJSON))
}

func authSetup(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var responseBytes []uint8 = nil
	var sensor Sensor
	
	if len(args) < 2 {
		return errorResponse("Needs Sensor ID and RTU ID!!!", 700)
	}

	sensorID := args[0]
	rtuID := args[1]

	// get sensor status
	bytes, err := stub.GetState(sensorID)
	if err != nil {
		return errorResponse(err.Error(), 7)
	}
	 
	err = json.Unmarshal(bytes, &sensor)
	if err !=nil {
		return errorResponse("Failed to read sensor status", 700)
	}

	// check if rtu can fetch CRP
	if sensor.DeviceID != rtuID {
		return errorResponse("RTU does not have the privileges to pull CRPs", 7)
	}

	// keep track of active request
	activeCRP := len(sensor.Crp)
	request := Request{}
	request.Crp = sensor.Crp[rand.Intn(activeCRP)]
	request.Status = "active"

	// get timestamp (1)
	t := time.Now().UTC()
	ts := uint32(t.UnixNano())
	request.TimeStamp = ts
	tsBytes := make([]byte, 4)
	// ts_size := binary.PutVarint(tsBytes, ts)
	binary.LittleEndian.PutUint32(tsBytes, ts)

	// get hmac (2)
	strChallenge := strings.Split(request.Crp.Challenge, ",")
	strResponse := strings.Split(request.Crp.Response, ",")
	 
	// convert challenge to uint16
	var challengeBytes []byte = nil
	for _, c := range strChallenge {
		if c != "" {
			value, _ := strconv.Atoi(c)
			// get bytes
			valueBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(valueBytes, uint16(value))
			for _, b := range valueBytes {
				challengeBytes = append(challengeBytes, b)
			}
		}
		// break
	}

	// convert response to uint8
	for _, r := range strResponse {
		if r != "" {
			value, _ := strconv.Atoi(r)
			responseBytes = append(responseBytes, uint8(value))
		}
	}

	// concat challenge and TS
	// var message []byte = nil
	for _, t := range tsBytes {
		challengeBytes = append(challengeBytes, t)
	}
	mac := hmac.New(sha256.New, responseBytes)
	mac.Write(challengeBytes)
	h1 := mac.Sum(nil)

	// get TS' (3)
	responseHash := sha256.Sum256(responseBytes)
	
	fmt.Println("string response: ", strResponse)
	fmt.Println("uint8 response: ", responseBytes)
	// fmt.Println("hash response: ", responseHash)
	// fmt.Println("hash response: ", hex.EncodeToString(responseHash[:]))

	fmt.Println("ts bytes: ", tsBytes)
	tsBytes[0] ^= responseHash[0]
	tsBytes[1] ^= responseHash[1]
	tsBytes[2] ^= responseHash[2]
	tsBytes[3] ^= responseHash[3]
	ts = binary.LittleEndian.Uint32(tsBytes)
	
	fmt.Println("challenge: ", challengeBytes[:10])
	fmt.Println("hmac: ", hex.EncodeToString(h1))
	fmt.Println("hmac bytes: ", h1)
	fmt.Println("ts': ", fmt.Sprint(ts))
	fmt.Println("ts prime bytes: ", tsBytes)

	sensor.ActRequest = request
	//convert to JSON and store the sensor struct
	sensorJSON, _ := json.Marshal(sensor)
	stub.PutState(sensorID, []byte(sensorJSON))

	// return setup to device
	// setupResponse := append(challengeBytes, h1...)
	// setupResponse = append(setupResponse, tsBytes[:tSize]...)
	// return shim.Success(setupResponse) 
	setupResponse := request.Crp.Challenge + "-" + hex.EncodeToString(h1) + "-" + fmt.Sprint(ts)
	return shim.Success([]byte(setupResponse)) 
}

func findEvent(devices []Device, deviceID string, state string) (int, bool){
	for i, device := range devices {
		if device.DeviceID == deviceID {
			// fmt.Println("Found Device")
			for _, event := range device.Event {
				if event == state {
					return i, true
				}	
			}
		}
	}
	return -1, false
}

func checkDependency(process Process, deviceID string, state string) (bool) {
	var prvEvent string
	for _, device := range process.Device {
		if device.DeviceID == deviceID {
			prvEvent = device.PrvState
			break
		}
	}

	for _, dependency := range process.Dependency {
		if dependency.PrvEvent == prvEvent && dependency.NxtEvent == state {
			return true
		}
	}

	return false
}

func verify(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var process Process

	if len(args) < 4 {
		return errorResponse("Missing Args!!!", 700)
	}

	processID := args[0]
	deviceID := args[1]
	state, _ := hex.DecodeString(args[2])
	// proofID := args[3]

	// get process
	bytes, err := stub.GetState(processID)
	if err != nil {
		return errorResponse(err.Error(), 7)
	}

	err = json.Unmarshal(bytes, &process)
	if err !=nil {
	 	return errorResponse("Failed to read process status", 700)
	}

	//verify device ID
	// construct proof ID
	// update devie state
	var id []byte = nil
	var key = []byte{'S', 'E', 'C', 'R'}
	for _, device := range process.Device {
		if device.DeviceID == deviceID {
			idBytes := []byte(device.DeviceID)
			id = append(id, idBytes...)
			stateBytes := []byte(device.State)
			id = append(id, stateBytes...)
			// XoR
			id[0] ^= key[0]
			id[1] ^= key[1]
			id[2] ^= key[2]
			// id[3] ^= key[3]
			break
		}
	}
	// if proofID != id {
	// 	fmt.Println("Send alarm data to SCADA!!!!!")
	// 	return errorResponse("Failed to authenticate!!!!", 700)
	// }
	// fmt.Println(string(id))
	fmt.Println("**Identity Valid**")
	
	// verify valid event
	// fmt.Println(processID)
	// fmt.Println(deviceID)
	// fmt.Println(state)
	// fmt.Println(hex.DecodeString(state))
	_, validEvent := findEvent(process.Device, deviceID, string(state))
	if !validEvent {
		fmt.Println("Send alarm to SCADA!!!!!")
		return errorResponse("Invalid Event!!!!", 700)
	}
	fmt.Println("**Event Valid**")

	// verify dependency
	validDep := checkDependency(process, deviceID, string(state))
	if !validDep {
		fmt.Println("Send alarm data to SCADA!!!!!")
		return errorResponse("Invalid Dependency!!!!", 700)
	}
	fmt.Println("**Dependency Valid**")

	// update devie state
	for idx, device := range process.Device {
		if device.DeviceID == deviceID {
			process.Device[idx].PrvState = process.Device[idx].State
			process.Device[idx].State = string(state)
			// device.State = state
		}
	}

	//convert to JSON and store the sensor struct
	processJSON, _ := json.Marshal(process)
	stub.PutState(process.ProcessID, []byte(processJSON))
	return shim.Success([]byte(processJSON))

	// return successResponse("sucess")



	// var responseBytes []uint8 = nil
	// var sensorData []uint8 = nil
	// var sensor Sensor

	// sensorID := args[0]

	// if len(args) < 3 {
	// 	return errorResponse("Needs data and dataHash!!!", 700)
	// }

	// // get sensor status
	// bytes, err := stub.GetState(sensorID)
	// if err != nil {
	// 	return errorResponse(err.Error(), 7)
	// }
	 
	// err = json.Unmarshal(bytes, &sensor)
	// if err !=nil {
	// 	return errorResponse("Failed to read sensor status", 700)
	// }

	// if sensor.ActRequest.Status != "active" {
	// 	return errorResponse("Sensor does not have an active Request!!!!!", 700)
	// }

	// // pull response from Active request
	// // convert response to uint8
	// for _, r := range strings.Split(sensor.ActRequest.Crp.Response, ",") {
	// 	if r != "" {
	// 		value, _ := strconv.Atoi(r)
	// 		responseBytes = append(responseBytes, uint8(value))
	// 	}
	// }

	// // read data and get re-create hash
	// data, _ := strconv.Atoi(args[1])
	// sensorData = append(sensorData, uint8(data))

	// mac := hmac.New(sha256.New, responseBytes)
	// mac.Write(sensorData)
	// masterHash := mac.Sum(nil)

	// // read data hmac
	// sensorHash, _ := hex.DecodeString(args[2])

	// // verify hmac
	// if !hmac.Equal(sensorHash, masterHash) {
	// 	fmt.Println("")
	// 	// return errorResponse("Unable to verify data integrity!!!", 700)
	// }

	// // reset active request
	// sensor.ActRequest.Crp.Challenge = ""
	// sensor.ActRequest.Crp.Response = ""
	// sensor.ActRequest.TimeStamp = 0
	// sensor.ActRequest.Status = ""
	
	// //convert to JSON and store the sensor struct
	// sensorJSON, _ := json.Marshal(sensor)
	// stub.PutState(sensorID, []byte(sensorJSON))

	// fmt.Println("Forward data to SCADA: ", sensorData)

	// // verify Hmac
	// return shim.Success([]byte(sensorJSON))
}

 // processJSON creates a JSON for representing the sensor status
func processJSON(ProcessID, Status string) (string) {
	 return "{\"process\":\""+ProcessID+"\", \"status\":"+Status+ "}"
}

func errorResponse(err string, code  uint ) peer.Response {
	codeStr := strconv.FormatUint(uint64(code), 10)
	// errorString := "{\"error\": \"" + err +"\", \"code\":"+codeStr+" \" }"
	errorString := "{\"error\":" + err +", \"code\":"+codeStr+" \" }"
	return shim.Error(errorString)
}

func successResponse(dat string) peer.Response {
	success := "{\"response\": " + dat +", \"code\": 0 }"
	return shim.Success([]byte(success))
}

// Chaincode registers with the Shim on startup
func main() {
	fmt.Println("Started....")
	err := shim.Start(new(SPAIChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
