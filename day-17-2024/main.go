package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type OpCodeType uint64

const (
	ADV OpCodeType = iota
	BXL
	BST
	JNZ
	BXC
	OUT
	BDV
	CDV
)

type OperandType uint64

const (
	LITERAL OperandType = iota
	COMBO
)

type Operand struct {
	operandType OperandType
	value       uint64
}

func (o Operand) getOperandValue(vm *VM) uint64 {
	switch o.operandType {
	case LITERAL:
		return o.value
	case COMBO:
		if o.value <= 3 {
			return o.value
		}
		if o.value == 4 {
			return vm.aReg
		}
		if o.value == 5 {
			return vm.bReg
		}
		if o.value == 6 {
			return vm.cReg
		}
		panic(fmt.Sprintf("Invalid value for COMBO operand: %d", o.value))
	default:
		panic(fmt.Sprintf("Invalid operand type: %d", o.operandType))
	}
}

type Instruction struct {
	operator OpCodeType
	operand  Operand
}

func (i Instruction) applyInstruction(vm *VM) {
	switch i.operator {
	case JNZ:
		if (*vm).aReg != 0 {
			(*vm).ip = int(i.operand.getOperandValue(vm)) / 2
			return
		}
	case ADV:
		opVal := datastructures.Pow(uint64(2), i.operand.getOperandValue(vm))
		aReg := (*vm).aReg
		(*vm).aReg = aReg / opVal
	case BDV:
		opVal := datastructures.Pow(uint64(2), i.operand.getOperandValue(vm))
		aReg := (*vm).aReg
		(*vm).bReg = aReg / opVal
	case CDV:
		opVal := datastructures.Pow(uint64(2), i.operand.getOperandValue(vm))
		aReg := (*vm).aReg
		(*vm).cReg = aReg / opVal
	case BST:
		bReg := i.operand.getOperandValue(vm) % 8
		(*vm).bReg = bReg
	case BXL:
		(*vm).bReg = (*vm).bReg ^ i.operand.getOperandValue(vm)
	case BXC:
		(*vm).bReg = (*vm).bReg ^ (*vm).cReg
	case OUT:
		opVal := i.operand.getOperandValue(vm) % 8
		(*vm).outputBuffer = append((*vm).outputBuffer, opVal)
	}

	(*vm).ip++
}

type VM struct {
	aReg         uint64
	bReg         uint64
	cReg         uint64
	ip           int // instruction pointer
	instructions []Instruction
	outputBuffer []uint64
}

func (vm VM) runProgram() []uint64 {
	for vm.ip < len(vm.instructions) {
		vm.instructions[vm.ip].applyInstruction(&vm)
	}

	return vm.outputBuffer
}

func (vm VM) instructionsToArr() []uint64 {
	program := []uint64{}
	for idx := 0; idx < len(vm.instructions); idx++ {
		instr := vm.instructions[idx]
		program = append(program, uint64(instr.operator))
		program = append(program, uint64(instr.operand.value))
	}
	return program
}

func readFileAsString(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, nil
}

func buildVM(data []string) VM {
	vm := VM{}
	regAStr := strings.Split(data[0], " ")[2]
	regBStr := strings.Split(data[1], " ")[2]
	regCStr := strings.Split(data[2], " ")[2]

	regA, _ := strconv.Atoi(regAStr)
	regB, _ := strconv.Atoi(regBStr)
	regC, _ := strconv.Atoi(regCStr)

	vm.aReg = uint64(regA)
	vm.bReg = uint64(regB)
	vm.cReg = uint64(regC)

	programStr := strings.Split(data[4], " ")[1]
	split := strings.Split(programStr, ",")

	instructions := []Instruction{}

	for idx := 0; idx < len(split); idx += 2 {
		operatorStr := split[idx]
		operandStr := split[idx+1]

		operator, _ := strconv.Atoi(operatorStr)
		opr, _ := strconv.Atoi(operandStr)
		operand := uint64(opr)

		opType := OpCodeType(operator)

		var oper Operand

		switch opType {
		case ADV:
			oper = Operand{operandType: COMBO, value: operand}
		case BXL:
			oper = Operand{operandType: LITERAL, value: operand}
		case BST:
			oper = Operand{operandType: COMBO, value: operand}
		case JNZ:
			oper = Operand{operandType: LITERAL, value: operand}
		case BXC:
			oper = Operand{operandType: LITERAL, value: operand}
		case OUT:
			oper = Operand{operandType: COMBO, value: operand}
		case BDV:
			oper = Operand{operandType: COMBO, value: operand}
		case CDV:
			oper = Operand{operandType: COMBO, value: operand}
		default:
			panic("Invalid optype")
		}

		instructions = append(instructions, Instruction{operator: opType, operand: oper})
	}

	vm.instructions = instructions
	return vm
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	vm := buildVM(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		res := vm.runProgram()
		strs := []string{}
		for _, val := range res {
			strs = append(strs, strconv.Itoa(int(val)))
		}
		println(strings.Join(strs, ","))
	} else {
		type Tuple struct {
			matchCount uint64
			guess      uint64
		}

		expectedProgramOutput := vm.instructionsToArr()
		potentialSols := []Tuple{}
		potentialSols = append(potentialSols, Tuple{0, 0})
	out:
		for len(potentialSols) > 0 {
			posSol := potentialSols[0]
			potentialSols = potentialSols[1:]

			for pows := uint64(0); pows < 8; pows++ {
				vm.aReg = posSol.guess + pows
				vm.bReg = 0
				vm.cReg = 0
				vm.ip = 0
				vm.outputBuffer = []uint64{}
				progOutput := vm.runProgram()

				if len(progOutput) > 0 && expectedProgramOutput[uint64(len(expectedProgramOutput))-posSol.matchCount-1] == progOutput[0] {
					if posSol.matchCount+1 == uint64(len(expectedProgramOutput)) {
						println(posSol.guess + pows)
						break out
					}
					potentialSols = append(potentialSols, Tuple{posSol.matchCount + 1, (posSol.guess + pows) * 8})
				}
			}
		}
	}
}
