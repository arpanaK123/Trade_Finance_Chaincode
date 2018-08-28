package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("query"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestTradefinance_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	// Init a=123 b=234
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})

	checkState(t, stub, "a", "100")
	checkState(t, stub, "b", "200")
}

func TestTradefinance_Query(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("exporter", scc)

	// Init a=345 b=456
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("a"), []byte("300"), []byte("b"), []byte("400")})

	// Query a
	checkQuery(t, stub, "a", "300")

	// Query b
	checkQuery(t, stub, "b", "400")
}

func TestTradefinance_Invoke(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("exporter", scc)

	// Init a=567 b=678
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("a"), []byte("567"), []byte("b"), []byte("678")})

	// Invoke a->b for 123
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("a"), []byte("b"), []byte("123")})
	checkQuery(t, stub, "a", "444")
	checkQuery(t, stub, "b", "801")

	// Invoke b->a for 234
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("b"), []byte("a"), []byte("234")})
	checkQuery(t, stub, "a", "678")
	checkQuery(t, stub, "b", "567")
	checkQuery(t, stub, "a", "678")
	checkQuery(t, stub, "b", "567")
}