package adsbdecoder

import (
	"strings"
	"testing"
)

func TestHexToBin(t *testing.T) {
	s := MustHex2Bin("6E406B")

	bin := strings.Join(s, "")

	if bin != "011011100100000001101011" {
		t.Errorf("Expected 011011100100000001101011 : %s--", bin)
	}
}

func TestCRCDecode(t *testing.T) {
	checksum, err := CRC(NewMessageContext("8D406B902015A678D4D220AA4BDA"), false)
	if err != nil {
		t.Error(err)
	}

	if checksum != "000000000000000000000000" {
		t.Errorf("oops")
	}
}

func TestCRCEncode(t *testing.T) {
	checksum, err := CRC(NewMessageContext("8D406B902015A678D4D220AA4BDA"), true)
	if err != nil {
		t.Error(err)
	}

	b := MustHex2Bin("AA4BDA")

	if checksum != strings.Join(b, "") {
		t.Errorf("oops")
	}
}

func TestICAO(t *testing.T) {
	testdata := map[string]string{
		"8D406B902015A678D4D220AA4BDA": "406B90",
		"A0001839CA3800315800007448D9": "400940",
		"A000139381951536E024D4CCF6B5": "3C4DD2",
		"A000029CFFBAA11E2004727281F1": "4243D0",
	}

	for hex := range testdata {
		ctx := NewMessageContext(hex)

		icao, err := ICAO(ctx)
		if err != nil {
			t.Error(err)
		}
		if icao != testdata[hex] {
			t.Errorf("%s %s", icao, testdata[hex])
		}
	}

}

func TestModeSAltcode(t *testing.T) {
	code, err := AltCode(NewMessageContext("A02014B400000000000000F9D514"))
	if err != nil {
		t.Error(err)
		return
	}

	if code != 32300 {
		t.Errorf("%d %d", code, 32300)
	}
}

func TestModeSIdCode(t *testing.T) {
	code := IDCODE(NewMessageContext("A800292DFFBBA9383FFCEB903D01"))
	if code != "1346" {
		t.Errorf("%s %s", code, "1346")
	}
}

func TestGreyCodeToAltitude(t *testing.T) {
	testData := map[string]int{
		"00000000010": -1000,
		"00000001010": -500,
		"00000011011": -100,
		"00000011010": 0,
		"00000011110": 100,
		"00000010011": 600,
		"00000110010": 1000,
		"00001001001": 5800,
		"00011100100": 10300,
		"01100011010": 32000,
		"01110000100": 46300,
		"01010101100": 50200,
		"11011110100": 73200,
		"10000000011": 126600,
		"10000000001": 126700,
	}

	for grey := range testData {
		alt := Gray2Alt(grey)
		if alt != testData[grey] {
			t.Errorf("%d %d", alt, testData[grey])
			return
		}
	}
}