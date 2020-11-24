package main

import (
	"context"
	"os/signal"
	"time"

	"fmt"

	"github.com/stianeikeland/go-rpio"
	cvnet2 "github.com/DatanoiseTV/cvnet2-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"syscall"
)

type server struct {
	cvnet2.UnimplementedCVServer
}

const (
	PIXI_WRITE byte = iota
	PIXI_READ
)

const (
	PIXI_DEVICE_ID byte = iota
	PIXI_INTERRUPT
	PIXI_ADC_DATA_STATUS_0_15
	PIXI_ADC_DATA_STATUS_16_19
	PIXI_OVERCURRENT_STATUS_0_15
	PIXI_OVERCURRENT_STATUS_16_19
	PIXI_GPI_STATUS_0_15
	PIXI_GPI_STATUS_16_19
	PIXI_INT_TEMP_DATA
	PIXI_EXT1_TEMP_DATA
	PIXI_EXT2_TEMP_DATA
	PIXI_GPI_DATA_0_15
	PIXI_GPI_DATA_16_19
	PIXI_GPO_DATA_0_15
	PIXI_GPO_DATA_16_19
	PIXI_PLACEHOLDER_0
	PIXI_DEVICE_CTRL
	PIXI_INTERRUPT_MASK
	PIXI_GPI_IRQ_MODE_0_7
	PIXI_GPI_IRQ_MODE_8_15
	PIXI_GPI_IRQ_MODE_16_19
	PIXI_DAC_PRESET_DATA1
	PIXI_DAC_PRESET_DATA2
	PIXI_TEMP_MON_CONFIG
	PIXI_TEMP_INT_HIGH_THRESHOLD
	PIXI_TEMP_INT_LOW_THRESHOLD
	PIXI_TEMP_EXT1_HIGH_THRESHOLD
	PIXI_TEMP_EXT1_LOW_THRESHOLD
	PIXI_TEMP_EXT2_HIGH_THRESHOLD
	PIXI_TEMP_EXT2_LOW_THRESHOLD
	PIXI_PLACEHOLDER_1
	PIXI_PORT0_CONFIG
	PIXI_PORT1_CONFIG
	PIXI_PORT2_CONFIG
	PIXI_PORT3_CONFIG
	PIXI_PORT4_CONFIG
	PIXI_PORT5_CONFIG
	PIXI_PORT6_CONFIG
	PIXI_PORT7_CONFIG
	PIXI_PORT8_CONFIG
	PIXI_PORT9_CONFIG
	PIXI_PORT10_CONFIG
	PIXI_PORT11_CONFIG
	PIXI_PORT12_CONFIG
	PIXI_PORT13_CONFIG
	PIXI_PORT14_CONFIG
	PIXI_PORT15_CONFIG
	PIXI_PORT16_CONFIG
	PIXI_PORT17_CONFIG
	PIXI_PORT18_CONFIG
	PIXI_PORT19_CONFIG // 0x33
	PIXI_PLACEHOLDER_2 // 0x34
	PIXI_PLACEHOLDER_3 // 0x35
	PIXI_PLACEHOLDER_4 // 0x36
	PIXI_PLACEHOLDER_5 // 0x37
	PIXI_PLACEHOLDER_6 // 0x38
	PIXI_PLACEHOLDER_7 // 0x39
	PIXI_PORT0_ADC_DATA
	PIXI_PORT1_ADC_DATA
	PIXI_PORT2_ADC_DATA
	PIXI_PORT3_ADC_DATA
	PIXI_PORT4_ADC_DATA
	PIXI_PORT5_ADC_DATA
	PIXI_PORT6_ADC_DATA
	PIXI_PORT7_ADC_DATA
	PIXI_PORT8_ADC_DATA
	PIXI_PORT9_ADC_DATA
	PIXI_PORT10_ADC_DATA
	PIXI_PORT11_ADC_DATA
	PIXI_PORT12_ADC_DATA
	PIXI_PORT13_ADC_DATA
	PIXI_PORT14_ADC_DATA
	PIXI_PORT15_ADC_DATA
	PIXI_PORT16_ADC_DATA
	PIXI_PORT17_ADC_DATA
	PIXI_PORT18_ADC_DATA
	PIXI_PORT19_ADC_DATA
	PIXI_PLACEHOLDER_8  // 0x53
	PIXI_PLACEHOLDER_9  // 0x54
	PIXI_PLACEHOLDER_10 // 0x55
	PIXI_PLACEHOLDER_11 // 0x56
	PIXI_PLACEHOLDER_12 // 0x57
	PIXI_PLACEHOLDER_13 // 0x58
	PIXI_PLACEHOLDER_14 // 0x59
	PIXI_PORT0_DAC_DATA
	PIXI_PORT1_DAC_DATA
	PIXI_PORT2_DAC_DATA
	PIXI_PORT3_DAC_DATA
	PIXI_PORT4_DAC_DATA
	PIXI_PORT5_DAC_DATA
	PIXI_PORT6_DAC_DATA
	PIXI_PORT7_DAC_DATA
	PIXI_PORT8_DAC_DATA
	PIXI_PORT9_DAC_DATA
	PIXI_PORT10_DAC_DATA
	PIXI_PORT11_DAC_DATA
	PIXI_PORT12_DAC_DATA
	PIXI_PORT13_DAC_DATA
	PIXI_PORT14_DAC_DATA
	PIXI_PORT15_DAC_DATA
	PIXI_PORT16_DAC_DATA
	PIXI_PORT17_DAC_DATA
	PIXI_PORT18_DAC_DATA
	PIXI_PORT19_DAC_DATA
)

const (
	CH_MODE_HIZ = iota
	CH_MODE_PI
	CH_MODE_DIDIR_LT_TERM
	CH_MODE_GPO_REG
	CH_MODE_GPO_UNI
	CH_MODE_DAC
	CH_MODE_DAC_ADC_MON
	CH_MODE_ADC_P
	CH_MODE_ADC_DIFF_P
	CH_MODE_ADC_DIFF_N
	CH_MODE_DAC_ADC_DIFF_N
	CH_MODE_TERM_GPI_SW
	CH_MODE_TERM_REG_SW
)

const (
	CH_NO_RANGE = iota
	CH_0_TO_10P
	CH_5N_TO_5P
	CH_10N_TO_0
	CH_0_TO_2P5_5N_TO_5P
	CH_RES
	CH_0_TO_2P5_0_TO_10P
	CH_RES_2
)

const (
	TEMP_CHANNEL_INT = iota
	TEMP_CHANNEL_EXT0
	TEMP_CHANNEL_EXT1
)

const (
	DEVID = 0xFFFF
	ADCCTL = 0x0003
	DACCTL = 0x000C
	ADCCONV = 0x0030
	DACREF = 0x0040
	THSHDN = 0x0080
	TMPCTL = 0x0700
	TMPCTLINT = 0x0100
	TMPCTLEXT1 = 0x0200
	TMPCTLEXT2 = 0x0400
	TMPPER = 0x0800
	RS_CANCEL = 0x1000
	LPEN = 0x2000
	BRST = 0x4000
	RESET = 0x8000
)

const (
	ADC_MODE_IDLE = iota
	ADC_MODE_SWEEP
	ADC_MODE_CONV
	ADC_MODE_CONT
)

const (
	TMPINTMONCFG = 0x0003
	TMPEXT1MONCFG = 0x000C
	TMPEXT2MONCFG = 0x0030
	TMPINTHI = 0x0FFF
)

const (
	FUNCPRM = 0x0FFF
	FUNCID = 0xF000
	FUNCPRM_ASSOCIATED_PORT = 0x001F
	FUNCPRM_NR_OF_SAMPLES = 0x00E0
	FUNCPRM_RANGE = 0x0700
	FUNCPRM_AVR_INV = 0x0800
	FUNCID_MODE0_HIGHZ = 0x0000
	ADCDAT = 0x0FFF
)

func writeRegister(address byte, value int){
	data := []byte{byte(address << 0x01) | PIXI_WRITE, byte(value >> 8),  byte(value & 0xFF)}
	rpio.SpiTransmit(data...)
}

func readRegister(address byte) int {
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		panic(err)
	}
	data := []byte{byte(address << 0x01) | PIXI_READ, 0x00, 0x00}
	rpio.SpiExchange(data)
	return int(data[1])<<8 | int(data[2])
}

func configChannel(channel byte, channel_mode int, dac_dat int, rng int, adc_ctl int){
	info := 0
	info = readRegister(PIXI_DEVICE_CTRL)
	writeRegister (PIXI_DEVICE_CTRL, info | DACREF | 0 );

	info = readRegister(PIXI_DEVICE_CTRL)
	writeRegister(0x60 + channel, dac_dat)
	writeRegister ( 0x20 + channel, ( ( (channel_mode << 12 ) & FUNCID ) |
		( (rng << 8 ) & FUNCPRM_RANGE ) ) )
}

func readAnalog(channel byte) int{
	return readRegister(0x40 + channel)
}

func writeAnalog(channel byte, value int) int{
	result := 0
	channel_func := 0
	channel_func = readRegister(0x20 + channel)
	channel_func = ( channel_func & FUNCID ) >> 12
	if(channel_func == 5){
		writeRegister(0x60 + channel, value)
		result = readRegister(0x60 + channel)
	}

	return result
}


func readTemperature(channel byte) float64 {
	result := 0.0
	rawresult := 0
	sign := 0
	rawresult = readRegister(PIXI_INT_TEMP_DATA + channel)
	sign = (rawresult & 0x0800) >> 1;
	if(sign == 1){
		rawresult = ( ( rawresult & 0x07FF ) ^ 0x07FF ) + 1
	}
	result = float64(0.125) * float64( rawresult & 0x0007 )
	result = float64(result) + float64( ( rawresult >> 3) & 0x01FF )
	if (sign == 1){
		result = result * -1
	}
	return result

}

func (s *server ) PinMode(ctx context.Context, in *cvnet2.ConfigMessage) (*cvnet2.ConfigMessage, error){


	channelMode := 0
	switch in.Mode {
		case cvnet2.ConfigMessage_CV_IN:
			channelMode = 0x07
			break
		case cvnet2.ConfigMessage_CV_OUT:
			channelMode = 0x05
			break
		case cvnet2.ConfigMessage_GATE_IN:
			channelMode = 0x01
			break
		case cvnet2.ConfigMessage_GATE_OUT:
			channelMode = 0x04
			break
		default:
			return nil, nil;
	}

	channelRange := 0
	switch in.Range {
	case cvnet2.ConfigMessage_NEG_FIVE_TO_FIVE:
		channelRange = 0x03
		break
	case cvnet2.ConfigMessage_ZERO_TO_FIVE:
		channelRange = 0x01
		break
	case cvnet2.ConfigMessage_ZERO_TO_TEN:
		channelRange = 0x01
		break
	default:
		return nil, nil
	}

	adcCtl := 0
	if in.Mode == cvnet2.ConfigMessage_CV_IN {
		adcCtl = 0x03 // continuous mode
	}

	configChannel (byte(in.Channel), channelMode, 0, channelRange, adcCtl )

	return nil, nil
}


func (s *server ) ReadCV(ctx context.Context, in *cvnet2.CVMessage) (*cvnet2.CVMessage, error){
	return nil, nil
}

func (s *server ) WriteCV(ctx context.Context, in *cvnet2.CVMessage) (*cvnet2.CVMessage, error){
	return nil, nil
}

func (s *server ) ReadGate(ctx context.Context, in *cvnet2.GateMessage) (*cvnet2.GateMessage, error){
	return nil, nil
}

func (s *server ) WriteGate(ctx context.Context, in *cvnet2.GateMessage) (*cvnet2.GateMessage, error){
	return nil, nil
}

func (s *server ) ReadCVStream(in *cvnet2.CVMessage, src cvnet2.CV_ReadCVStreamServer) (error){
	return nil
}

func (s *server ) WriteCVStream(src cvnet2.CV_WriteCVStreamServer) (error){
	return nil
}

func (s *server ) ReadGateStream(in *cvnet2.GateMessage, src cvnet2.CV_ReadGateStreamServer) (error){
	return nil
}

func (s *server ) WriteGateStream(src cvnet2.CV_WriteGateStreamServer) (error){
	return nil
}

func main() {
	// We need to be root to access GPIO
	if os.Getuid() != 0 {
		fmt.Println("Sorry, root required.")
		os.Exit(1)
	}

	// Set nice priority to -20 to allow low-latency output on PREEMPT_RT Linux Kernel
	syscall.Setpgid(0, 0); syscall.Setpriority(syscall.PRIO_PGRP, 0, -20)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Printf("You pressed ctrl + C. User interrupted infinite loop.")
		os.Exit(0)
	}()

	if err := rpio.Open(); err != nil {
		panic(err)
	}

	rpio.SpiSpeed(16000000)
	rpio.SpiChipSelect(0) // Select CE0 slave

	version := readRegister(0x00);
	if version == 0x0424 {
		fmt.Println("Found PIXI.")
	} else {
		fmt.Println("Error: Pixi not found.")
	}

	writeRegister (PIXI_DEVICE_CTRL, (0xFFFF-BRST) | THSHDN);

	control := readRegister(PIXI_DEVICE_CTRL)
	writeRegister(PIXI_DEVICE_CTRL, control | 0xffff - RS_CANCEL)
	writeRegister(PIXI_DEVICE_CTRL, 0x0230)
	control = readRegister(PIXI_DEVICE_CTRL)
	writeRegister (PIXI_DEVICE_CTRL, control | TMPCTLINT | TMPCTLEXT1 | TMPCTLEXT2);
	control = readRegister(PIXI_DEVICE_CTRL)
	fmt.Println(control)

	writeRegister(PIXI_DEVICE_CTRL, control | ADC_MODE_CONT | DACREF | 0)


	temperature := readTemperature (0)
	fmt.Println(temperature)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	cvnet2.RegisterCVServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
