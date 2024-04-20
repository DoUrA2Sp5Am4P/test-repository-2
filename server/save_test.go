package main

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	t.Run("Read and save test", TestSaveAndRead)
}

func TestAll(t *testing.T) {
	os.Remove("store.db")
	t.Run("TestRegisterFunc", TestRegisterFunc)
	t.Run("TestGenerate", TestGenerate)
	t.Run("Test401Handler", Test401Handler)
	t.Run("TestRegister", TestRegister)
	t.Run("TestLogin", TestLogin)
	t.Run("TestAddition", TestAddition)
	t.Run("TestSubstraction", TestSubstraction)
	t.Run("TestMultiplication", TestMultiplication)
	t.Run("TestDividion", TestDividion)
	t.Run("TestGetExpressions", TestGetExpressions)
	t.Run("TestOperations", TestOperations)
	t.Run("TestGetResults", TestGetResults)
	t.Run("TestDatabase", TestDatabase)
	t.Run("TestSaveAndRead", TestSaveAndRead)
	t.Run("TestSave", TestSave)
	os.Remove("store.db")
}
