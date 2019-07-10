package loralib

import "github.com/stianeikeland/go-rpio"
import "time"
import "os"
import "log"

//Vars for go-rpio
var ssPin rpio.Pin = 25
var dio0 rpio.Pin = 4
var RST rpio.Pin = 17

var CHANNEL uint8 = 0

var sx1272 bool

const REG_VERSION = 0x42
const OPMODE_SLEEP = 0x00
const freq = 868100000
const REG_FRF_MSB = 0x06
const REG_FRF_MID = 0x07
const REG_FRF_LSB = 0x08
const REG_SYNC_WORD = 0x39

const (
	SF7  = 7
	SF8  = 8
	SF9  = 9
	SF10 = 10
	SF11 = 11
	SF12 = 12
)

const sf = SF7

const REG_MODEM_CONFIG1 = 0x1D
const REG_MODEM_CONFIG2 = 0x1E
const REG_MODEM_CONFIG3 = 0x26
const REG_SYMB_TIMEOUT_LSB = 0x1F
const REG_MAX_PAYLOAD_LENGTH = 0x23
const REG_PAYLOAD_LENGTH = 0x22
const PAYLOAD_LENGTH = 0x40
const REG_HOP_PERIOD = 0x24
const REG_FIFO_ADDR_PTR = 0x0D
const REG_FIFO_RX_BASE_AD = 0x0F
const REG_LNA = 0x0C
const LNA_MAX_GAIN = 0x23

const REG_OPMODE = 0x01
const not_OPMODE_MASK = 0xF8

const OPMODE_LORA = 0x80
const OPMODE_STANDBY = 0x01
const OPMODE_RX = 0x05
const OPMODE_TX = 0x03

const RegPaRamp = 0x0A
const RegPaConfig = 0x09
const RegPaDac = 0x5A

const RegDioMapping1 = 0x40
const MAP_DIO0_LORA_TXDONE = 0x40
const MAP_DIO1_LORA_NOP = 0x30
const MAP_DIO2_LORA_NOP = 0xC0
const REG_IRQ_FLAGS = 0x12
const REG_IRQ_FLAGS_MASK = 0x11
const not_IRQ_LORA_TXDONE_MASK = 0xF7
const REG_FIFO_TX_BASE_AD = 0x0E
const REG_FIFO = 0x00
const REG_FIFO_RX_CURRENT_ADDR = 0x10
const REG_RX_NB_BYTES = 0x13

const IRQ_LORA_TXDONE_MASK = 0x08

const REG_PKT_SNR_VALUE = 0x19

func InitiateRPIO() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	rpio.PinMode(ssPin, rpio.Output)
	rpio.PinMode(dio0, rpio.Input)
	rpio.PinMode(RST, rpio.Output)
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		panic(err)
	}
	rpio.SpiChipSelect(CHANNEL)
	rpio.SpiSpeed(500000)
}

func ConfigSend() {
	// Prepare to send block
	opmodeLora()
	opmode(OPMODE_STANDBY)
	writeReg(RegPaRamp, (readReg(RegPaRamp)&0xF0)|0x08) // set PA ramp-up time 50 uSec
	configPower(23)
}

func ReceiveMode() {
	opmode(OPMODE_RX)
	// log.Printf("Listening at SF%d on %f Mhz.\n", sf, float64(float64(freq)/1000000))
}

func SetupLoRa() {
	rpio.WritePin(RST, rpio.High)
	time.Sleep(100 * time.Millisecond)
	rpio.WritePin(RST, rpio.Low)
	time.Sleep(100 * time.Millisecond)
	var version byte = readReg(REG_VERSION)
	if version == 0x22 {
		// sx1272
		log.Println("SX1272 detected, starting.")
		sx1272 = true
	} else {
		// sx1276?
		rpio.WritePin(RST, rpio.Low)
		time.Sleep(100 * time.Millisecond)
		rpio.WritePin(RST, rpio.High)
		time.Sleep(100 * time.Millisecond)
		version = readReg(REG_VERSION)
		if version == 0x12 {
			// sx1276
			log.Println("SX1276 detected, Starting.")
			sx1272 = false
		} else {
			log.Println("Unrecognized transceiver.")
			//log.Printf("Transceiver version %x",version)
			os.Exit(1)
		}
	}

	opmode(OPMODE_SLEEP)

	//set frequency
	var frf uint64 = uint64(freq<<19) / 32000000
	writeReg(REG_FRF_MSB, byte(frf>>16))
	writeReg(REG_FRF_MID, byte(frf>>8))
	writeReg(REG_FRF_LSB, byte(frf>>0))

	writeReg(REG_SYNC_WORD, 0x34) //LoRaWAN public sync word

	if sx1272 {
		if sf == SF11 || sf == SF12 {
			writeReg(REG_MODEM_CONFIG1, 0x0B)
		} else {
			writeReg(REG_MODEM_CONFIG1, 0x0A)
		}
		writeReg(REG_MODEM_CONFIG2, (sf<<4)|0x04)
	} else {
		if sf == SF11 || sf == SF12 {
			writeReg(REG_MODEM_CONFIG3, 0x0C)
		} else {
			writeReg(REG_MODEM_CONFIG3, 0x04)
		}
		writeReg(REG_MODEM_CONFIG1, 0x72)
		writeReg(REG_MODEM_CONFIG2, (sf<<4)|0x04)
	}

	if sf == SF10 || sf == SF11 || sf == SF12 {
		writeReg(REG_SYMB_TIMEOUT_LSB, 0x05)
	} else {
		writeReg(REG_SYMB_TIMEOUT_LSB, 0x08)
	}
	writeReg(REG_MAX_PAYLOAD_LENGTH, 0x80)
	writeReg(REG_PAYLOAD_LENGTH, PAYLOAD_LENGTH)
	writeReg(REG_HOP_PERIOD, 0xFF)
	writeReg(REG_FIFO_ADDR_PTR, readReg(REG_FIFO_RX_BASE_AD))

	writeReg(REG_LNA, LNA_MAX_GAIN)
}

func readReg(addr byte) byte {
	var spibuf []byte
	spibuf = make([]byte, 2)
	selectreceiver()
	spibuf[0] = addr & 0x7F
	spibuf[1] = 0x00
	rpio.SpiExchange(spibuf)
	unselectreceiver()
	return byte(spibuf[1])
}

func opmode(mode byte) {
	writeReg(REG_OPMODE, readReg(REG_OPMODE)&not_OPMODE_MASK|mode)
}

func writeReg(addr byte, value byte) {
	var spibuf byte = addr | 0x80
	selectreceiver()
	rpio.SpiTransmit(spibuf, value)
	unselectreceiver()
}

func selectreceiver() {
	rpio.WritePin(ssPin, rpio.Low)
}

func unselectreceiver() {
	rpio.WritePin(ssPin, rpio.High)
}

func opmodeLora() {
	var u byte = OPMODE_LORA
	if sx1272 == false {
		u |= 0x8 // TBD: sx1276 high freq
	}
	writeReg(REG_OPMODE, u)
}

func configPower(pw int8) {
	if sx1272 == false {
		// no boost used for now
		if pw >= 17 {
			pw = 15
		} else if pw < 2 {
			pw = 2
		}
		// check board type for BOOST pin
		writeReg(RegPaConfig, byte(0x80|byte(pw&0xf)))
		writeReg(RegPaDac, readReg(RegPaDac)|0x4)
	} else {
		// set PA config (2-17 dBm using PA_BOOST)
		if pw > 17 {
			pw = 17
		} else if pw < 2 {
			pw = 2
		}
		writeReg(RegPaConfig, byte(0x80|byte(pw-2)))
	}
}

func Send(send_array []byte) {
	// set the IRQ mapping DIO0=TxDone DIO1=NOP DIO2=NOP
	writeReg(RegDioMapping1, MAP_DIO0_LORA_TXDONE|MAP_DIO1_LORA_NOP|MAP_DIO2_LORA_NOP)
	// clear all radio IRQ flags
	writeReg(REG_IRQ_FLAGS, 0xFF)
	// mask all IRQs but TxDone
	writeReg(REG_IRQ_FLAGS_MASK, not_IRQ_LORA_TXDONE_MASK)

	// initialize the payload size and address pointers
	writeReg(REG_FIFO_TX_BASE_AD, 0x00)
	writeReg(REG_FIFO_ADDR_PTR, 0x00)
	writeReg(REG_PAYLOAD_LENGTH, byte(len(send_array)))

	// download buffer to the radio FIFO
	writeBuf(REG_FIFO, send_array)
	// now we actually start the transmission
	opmode(OPMODE_TX)

	//log.Printf("send: %s\n", string(send_array))
	//log.Printf("Send packets at SF%d on %f Mhz.\n", sf, float64(float64(freq)/1000000))
}

func ClearReceiver() {
	// return transciever to receive mode
	// set the IRQ mapping DIO0=TxDone DIO1=NOP DIO2=NOP
	writeReg(RegDioMapping1, MAP_DIO0_LORA_TXDONE|MAP_DIO1_LORA_NOP|MAP_DIO2_LORA_NOP)
	// clear all radio IRQ flags
	writeReg(REG_IRQ_FLAGS, 0xFF)
	// ari
	// mask all IRQs WITH TxDone
	writeReg(REG_IRQ_FLAGS_MASK, IRQ_LORA_TXDONE_MASK)
}

func writeBuf(addr byte, send_array []byte) {
	selectreceiver()
	rpio.SpiTransmit(append([]byte{addr | 0x80}, send_array...)...)
	unselectreceiver()
}

func CheckReceivedBuffer() (bool, []byte) {
	var SNR int
	var rssicorr int

	if int(rpio.ReadPin(dio0)) == 1 {
		if a, message, receivedbytes := receive(); a {
			value := byte(readReg(REG_PKT_SNR_VALUE))
			if (value & 0x80) == 1 { //The SNR sign bit is 1
				// Invert and divide by 4
				value = ((value ^ 0xFF + 1) & 0xFF) >> 2
				SNR = int(-value)
			} else {
				//divide by 4
				SNR = int((value & 0xFF) >> 2)
			}

			if sx1272 {
				rssicorr = 139
			} else {
				rssicorr = 157
			}

			log.Printf("Packet RSSI: %d, ", int(readReg(0x1A))-rssicorr)
			log.Printf("RSSI: %d, ", int(readReg(0x1B))-rssicorr)
			log.Printf("SNR: %v, ", SNR)
			log.Printf("Encrypted length: %v", int(receivedbytes))
			log.Printf("\n")
			//log.Printf("Payload: %s\n", string(message))
			return true, message
			// Received message
		} else {
			return false, nil
		}
	} // dio0=1
	return false, nil
}

func receive() (bool, []byte, byte) {
	var payload []byte
	var receivedCount byte

	// clear rxDone
	writeReg(REG_IRQ_FLAGS, 0x40)

	irqflags := int(readReg(REG_IRQ_FLAGS))

	// payload crc: 0x20
	if (irqflags & 0x20) == 0x20 {
		log.Println("CRC error")
		writeReg(REG_IRQ_FLAGS, 0x20)
		return false, payload[:], 0
	} else {
		currentAddr := byte(readReg(REG_FIFO_RX_CURRENT_ADDR))
		receivedCount = byte(readReg(REG_RX_NB_BYTES))

		writeReg(REG_FIFO_ADDR_PTR, currentAddr)

		for i := 0; i < int(receivedCount); i++ {
			payload = append(payload, byte(readReg(REG_FIFO)))
		}

	}
	return true, payload[:], receivedCount
}
