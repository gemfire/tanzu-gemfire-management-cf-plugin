package main

import "testing"

func TestValidatePCCInstance(t *testing.T) {
	ourPCCInstance := "pcc1"
	pccInstancesAvailable := []string{"pcc1", "pcc2", "pcc3"}
	err := ValidatePCCInstance(ourPCCInstance, pccInstancesAvailable)
	if err != nil{
		t.Error("Unexpected Error")
	}
}
