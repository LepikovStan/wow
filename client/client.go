package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func isValidProof(challenge, response string, difficulty int) bool {
	targetPrefix := strings.Repeat("0", difficulty)
	combined := challenge + response
	hash := sha256.Sum256([]byte(combined))
	return strings.HasPrefix(hex.EncodeToString(hash[:]), targetPrefix)
}

// solveProofOfWork finds a valid proof-of-work response
func solveProofOfWork(challenge string, difficulty int) string {
	var nonce uint64
	for {
		nonceStr := strconv.FormatUint(nonce, 10)
		if isValidProof(challenge, nonceStr, difficulty) {
			return nonceStr
		}
		nonce++
	}
}

func main() {
	serverAddr := "wisdom-server-1:8080"

	for {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			text := strings.Split(scanner.Text(), ":")
			challenge := text[0]
			difficulty, err := strconv.Atoi(text[1])
			if err != nil {
				fmt.Println("Error get difficulty:", err)
				return
			}

			fmt.Println("Received Challenge:", challenge)
			fmt.Println("Received difficulty:", difficulty)

			proofOfWorkResponse := solveProofOfWork(challenge, difficulty)
			fmt.Println("Sending Proof of Work Response:", proofOfWorkResponse)

			_, err = fmt.Fprintf(conn, proofOfWorkResponse+"\n")
			if err != nil {
				fmt.Println("Error sending proof-of-work response:", err)
				return
			}

			if scanner.Scan() {
				response := scanner.Text()
				fmt.Println("Server Response:", response)
			}
		} else {
			fmt.Println("Error reading challenge from server:", scanner.Err())
		}
		time.Sleep(time.Second * 5)
	}
}
