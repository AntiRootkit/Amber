package main

import "crypto/rc4"
import "math/rand"
import "io/ioutil"
import "os/exec"
import "strconv"
import "time"
import "os"

func crypt() {

	verbose("Ciphering payload...", "*")

	if len(target.key) != 0 {
		payload, err := ioutil.ReadFile("Payload")
		ParseError(err, "Can't open payload file.")
		progress()
		if (len(target.key)%8) != 0 && target.staged == true {
			//tmp := make([]byte,(len(target.key)+(8-len(target.key)%8)))
			target.key = append(target.key, GenerateKey(8-(len(target.key)%8))...)
			verbose(string("Key size rounded to "+strconv.Itoa(len(target.key))), "*")
		}
		progress()
		payload = RC4(payload, target.key)
		payload_rc4, err2 := os.Create("Payload.rc4")
		ParseError(err2, "Can't create payload.rc4 file.")
		progress()
		payload_key, err3 := os.Create("Payload.key")
		ParseError(err3, "Can't create payload.rc4 file.")
		payload_rc4.Write(payload)
		payload_rc4.Write(target.key)
		payload_rc4.Close()
		payload_key.Close()
	} else {
		key := GenerateKey(target.KeySize)
		if target.KeySize != len(key) {
			verbose(string("Key size rounded to "+strconv.Itoa(len(key))), "*")
		}
		progress()
		payload, err := ioutil.ReadFile("Payload")
		ParseError(err, "Can't open payload file.")
		progress()
		payload = RC4(payload, key)
		payload_rc4, err2 := os.Create("Payload.rc4")
		ParseError(err2, "Can't create payload.rc4 file.")
		progress()
		payload_key, err3 := os.Create("Payload.key")
		ParseError(err3, "Can't create payload.rc4 file.")
		payload_rc4.Write(payload)
		payload_key.Write(key)
		payload_rc4.Close()
		payload_key.Close()
	}
	progress()
	verbose("Payload encrypted with RC4 algorithm\n", "*")

	hex, _ := exec.Command("sh", "-c", "xxd -i Payload.key").Output()
	verbose(string(hex), "B")

	remove("Payload")
}

func xor(Data []byte, Key []byte) []byte {
	for i := 0; i < len(Data); i++ {
		Data[i] = (Data[i] ^ (Key[(i % len(Key))]))
	}
	return Data
}

func RC4(data []byte, key []byte) []byte {
	c, e := rc4.NewCipher(key)
	ParseError(e, "While RC4 encryption !")
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)
	return dst
}

func GenerateKey(Size int) []byte {

	if target.staged == true && (Size%8) != 0 && Size >= 8 {
		Size += (8 - (Size % 8))
	}

	Key := make([]byte, Size)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < Size; i++ {
		Key[i] = byte(rand.Intn(255))
	}
	return Key
}
