package tftp

import (
	"crypto/rand"
	"log"
)

type Test struct {
}

// Request TestTFTPRequestAck Test Ack
func (t Test) Request() {
	// Create a new request packet
	// random byte key
	optMap := make(map[string][]byte)
	optMap["key"] = []byte("akjhlk;ashdjfh1394uy49r82y3984fhuihjasjkhdf")
	optMap["door"] = []byte("door")
	request, err := NewReq([]byte("test.txt"), []byte("octet"), 512, optMap)
	if err != nil {
		log.Fatal(err)
	}
	// Convert the request packet to a byte slice
	packet, err := request.ToBytes()
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := request.Parse(packet)
	if err != nil {
		return
	}

	bsize, _ := bytes.ToBytes()
	log.Printf("Request Packet: %d", len(bsize))
	TestEncryptDecrypt(bsize)
}

func (t Test) Data() {
	data, err := NewData(1, make([]byte, 512))
	if err != nil {
		log.Fatal(err)
	}
	packet := data.ToBytes()
	var data2 Data
	err = data2.Parse(packet)
	if err != nil {
		log.Fatal(err)
	}

	bsize := data2.ToBytes()
	log.Printf("Data Packet: %d", len(bsize))
	TestEncryptDecrypt(bsize)

}

func (t Test) Error() {
	// make 2 byte error code
	errPack := NewErr(02, []byte("File not found"))
	packet := errPack.ToBytes()
	var err2 Error
	err := err2.Parse(packet)
	if err != nil {
		return
	}
	bsize := err2.ToBytes()
	log.Printf("Error Packet: %d", len(bsize))
	TestEncryptDecrypt(bsize)
}

func (t Test) Ack() {
	ack := NewAck(02)
	packet := ack.ToBytes()
	var ack2 Ack
	err := ack2.Parse(packet)
	if err != nil {
		return
	}
	bsize := ack2.ToBytes()
	log.Printf("Ack Packet: %d", len(bsize))
	TestEncryptDecrypt(bsize)

}

func (t Test) Oack() {
	options := make(map[string][]byte)
	options["blksize"] = []byte("512")
	options["key"] = []byte("akjhlk;ashdjfh1394uy49r82y3984fhuihjasjkhdf")
	options["door"] = []byte("door")
	oack := NewOpt(options)
	packet := oack.ToBytes()
	var oack2 OptionAcknowledgement
	err := oack2.Parse(packet)
	if err != nil {
		return
	}
	bsize := oack2.ToBytes()
	log.Printf("Oack Packet: %d", len(bsize))
	TestEncryptDecrypt(bsize)

}

func (t Test) Test() {
	t.Request()
	t.Data()
	t.Error()
	t.Ack()
	t.Oack()
}

func xor(data []byte, key []byte) []byte {
	ciphertext := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ciphertext[i] = data[i] ^ key[i%len(key)]
	}
	return ciphertext
}

func decryptXOR(ciphertext []byte, key []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i++ {
		plaintext[i] = ciphertext[i] ^ key[i%len(key)]
	}
	return plaintext
}

func TestEncryptDecrypt(data []byte) {
	key := make([]byte, 128)
	_, err := rand.Read(key)
	if err != nil {
		return
	}

	ciphertext := xor(data, key)

	plaintext := decryptXOR(ciphertext, key)

	if string(plaintext) != string(data) {
		log.Fatal("Error: plaintext != data")
	} else {
		log.Printf("Success: plaintext == data")
	}
}