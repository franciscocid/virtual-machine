package chip8

import "fmt"

type Graphics interface {
	Clear()
	TogglePixel(x, y int) (isAlreadyToggled bool)
}

type Chip8 struct {
	CurrState    Chip8State
	StateHistory []Chip8State
	UI           Graphics
}

func (c *Chip8) LoadGame(gameData []uint8) {
	c.StateHistory = make([]Chip8State, 0)

	c.CurrState = Chip8State{}
	for i, data := range gameData {
		c.CurrState.Memory[i] = data
	}
}

func (c *Chip8) Tick() {
	opcode := c.CurrState.FetchOpcode()

	newState := c.ExecuteOpcode(opcode)

	c.StateHistory = append(c.StateHistory, c.CurrState)
	c.CurrState = newState
}

func (c *Chip8) ExecuteOpcode(opcode uint16) Chip8State {
	fmt.Printf("OP %04x\t", opcode)

	switch opcode {
	case 0x00E0:
		return c.clearScreen()
	case 0x00EE:
		return c.returnFromSubroutine()
	}

	addr := opcode & 0x0FFF
	x := uint8(opcode & 0x0F00 >> 8)
	y := uint8(opcode & 0x00F0 >> 4)
	value := uint8(opcode & 0x00FF)
	nibble := uint8(opcode & 0x000F)

	firstOpcodeByte := opcode >> 12
	switch firstOpcodeByte {
	case 0x0:
		return c.syscall(addr)
	case 0x1:
		return c.jumpToAddress(addr)
	case 0x2:
		return c.callSubroutine(addr)
	case 0x3:
		return c.skipIfVxEqualValue(x, value)
	case 0x4:
		return c.skipIfVxNotEqualValue(x, value)
	case 0x5:
		return c.skipIfVxEqualVy(x, y)
	case 0x6:
		return c.loadIntoVx(x, value)
	case 0x7:
		return c.addToVx(x, value)
	case 0x8:
		switch opcode & 0x000F {
		case 0x0:
			return c.loadVxIntoVy(x, y)
		case 0x1:
			return c.loadBitwiseVxOrVyIntoVx(x, y)
		case 0x2:
			return c.loadBitwiseVxAndVyIntoVx(x, y)
		case 0x3:
			return c.loadBitwiseVxExclusiveOrVyIntoVx(x, y)
		case 0x4:
			return c.addVyToVx(x, y)
		case 0x5:
			return c.subtractVxByVy(x, y)
		case 0x6:
			return c.shiftVxRight(x)
		case 0x7:
			return c.loadVySubtractedByVxIntoVx(x, y)
		case 0xE:
			return c.shiftVxLeft(x)
		}
	case 0x9:
		return c.skipIfVxNotEqualVy(x, y)
	case 0xA:
		return c.loadAddressIntoI(addr)
	case 0xB:
		return c.jumpToAddressPlusV0(addr)
	case 0xC:
		return c.loadRandomValueBitwiseAndValueIntoVx(addr)
	case 0xD:
		return c.drawSprite(x, y, nibble)
	case 0xE:
		switch opcode & 0x00FF {
		case 0x9E:
			return c.skipIfVxKeyIsPressed(x, y)
		case 0xA1:
			return c.skipIfVxKeyIsNotPressed(x, y)
		}
	case 0xF:
		switch opcode & 0x00FF {
		case 0x07:
			return c.loadDelayTimerIntoVx(x)
		case 0x0A:
			return c.waitButtonPressAndLoadIntoVx(x)
		case 0x15:
			return c.loadVxIntoDelayTimer(x)
		case 0x18:
			return c.loadVxIntoSoundTimer(x)
		case 0x1E:
			return c.addVxToI(x)
		case 0x29:
			return c.loadVxDigitSpriteAddressIntoI(x)
		case 0x33:
			return c.loadVxDigitsIntoI(x)
		case 0x55:
			return c.loadRangeV0ToVxIntoMemoryStartingFromI(x)
		case 0x65:
			return c.loadMemoryStartingFromIIntoRangeV0ToVx(x)
		}
	}
	return c.invalidOpcode()
}

func New() *Chip8 {
	return &Chip8{
		CurrState:    Chip8State{},
		StateHistory: []Chip8State{},
	}
}
