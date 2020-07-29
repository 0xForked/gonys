package pkg

/*
 * Gonys - A Notification Service for SMS
 *
 * GSM Core Functionality
 *
 * @author A. A. Sumitro <hello@aasumitro.id>
 * https://aasumitro.id
 */

import (
	"errors"
	"github.com/aasumitro/gonys/src/utils"
	"github.com/tarm/serial"
	"log"
	"strings"
	"time"
)

/*
 *--------------------------------------------------------------------------
 * GSM attributes
 *--------------------------------------------------------------------------
 *
 * GMS struct defines the config attributes of the external device,
 * to make a connection between the core-engine/system and external device.
 *
 * @var string ComPort | Serial port communication (hardware interface)
 * @var int BaudRate | Serial port transferring -
 * (rate information transferred in maximum 'x' bits per second)
 * @var os Port | External device port
 * @var string DeviceId | External device id/name
 */
type GSM struct {
	ComPort  string
	BaudRate int
	Port     *serial.Port
	DeviceId string
}

/*
 *--------------------------------------------------------------------------
 * Set The Config
 *--------------------------------------------------------------------------
 *
 * Here we will set the environment configuration.
 *
 * @return struct GSM
 */
func NewGSMModem(ComPort string, BaudRate int, DeviceId string) (modem *GSM) {
	return &GSM{
		ComPort: ComPort,
		BaudRate: BaudRate,
		DeviceId: DeviceId,
	}
}

/*
 *--------------------------------------------------------------------------
 * Make a connection
 *--------------------------------------------------------------------------
 *
 * Here we will be trying to connect with the core device,
 * the serial port configuration is needed and will be applying
 * to opens a serial port with the specified configuration.
 *
 * @var tarm/serial.Config config
 *
 * @do initialized
 * @return error if !initialized
 */
func (modem *GSM) Connect() (err error) {
	config := &serial.Config{
		Name: modem.ComPort,
		Baud: modem.BaudRate,
		ReadTimeout: time.Second,
	}

	modem.Port, err = serial.OpenPort(config)
	if err == nil {
		modem.init()
	}

	return err
}

/*
 *--------------------------------------------------------------------------
 * init this class
 *--------------------------------------------------------------------------
 *
 * Here we will trying to send a command to the external device/core
 *
 * @do initialized
 * @return error if !initialized
 */
func (modem *GSM) init() {
	modem.WriteCommand("ATE0\r\n", true)      // echo off
	modem.WriteCommand("AT+CMEE=1\r\n", true) // useful error messages
	modem.WriteCommand("AT+WIND=0\r\n", true) // disable notifications
	modem.WriteCommand("AT+CMGF=1\r\n", true) // switch to TEXT mode
}

/*
 *--------------------------------------------------------------------------
 *
 *--------------------------------------------------------------------------
 *
 *
 *
 * @return string
 */
func (modem *GSM) Expect(possibilities []string) (string, error) {
	readMax := 0

	for _, possibility := range possibilities {
		length := len(possibility)

		if length > readMax {
			readMax = length
		}
	}

	readMax = readMax + 2

	var status = ""

	buf := make([]byte, readMax)

	for i := 0; i < readMax; i++ {
		conn, err := modem.Port.Read(buf)
		if err != nil {
			panic(err)
		}

		if conn > 0 {
			status = string(buf[:conn])

			for _, possibility := range possibilities {
				if strings.HasSuffix(status, possibility) {
					log.Println("--- Expect:",
						utils.Transpose(strings.Join(possibilities, "|")),
						"Got:", utils.Transpose(status))
					return status, nil
				}
			}
		}
	}

	log.Println("--- Expect:",
		utils.Transpose(strings.Join(possibilities, "|")),
		"Got:", utils.Transpose(status),
		"(match not found!)")

	return status, errors.New("match not found")
}

/*
 *--------------------------------------------------------------------------
 *
 *--------------------------------------------------------------------------
 *
 *
 *
 * @return string
 */
func (modem *GSM) Send(command string)  {
	log.Println("---Send:", utils.Transpose(command))

	_ = modem.Port.Flush()

	_, err := modem.Port.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
}

/*
 *--------------------------------------------------------------------------
 *
 *--------------------------------------------------------------------------
 *
 *
 *
 * @return string
 */
func (modem *GSM) Read(number int) string {
	var output = ""

	buf := make([]byte, number)

	for i := 0; i < number; i++ {
		conn, err := modem.Port.Read(buf)
		if err != nil {
			panic(err)
		}

		if conn > 0 {
			output = string(buf[:conn])
		}
	}

	log.Printf("---Read(%d): %v", number, utils.Transpose(output))

	return output
}

/*
 *--------------------------------------------------------------------------
 *
 *--------------------------------------------------------------------------
 *
 *
 *
 * @return string
 */
func (modem *GSM) WriteCommand(command string, waitCallback bool) string  {
	modem.Send(command)

	if waitCallback {
		exc, err := modem.Expect([]string{"OK\r\n", "ERROR\r\n"})
		if err != nil {
			panic(err)
		}

		return exc
	}

	return modem.Read(1)
}

/*
 *--------------------------------------------------------------------------
 *
 *--------------------------------------------------------------------------
 *
 *
 *
 * @return string
 */
func (modem *GSM) WriteMessage(number string, message string) string {
	modem.Send("AT+CMGS=\"" + number + "\"\r")

	modem.Read(3)

	return modem.WriteCommand(message+string(26), true)
}

