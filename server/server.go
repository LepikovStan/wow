package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	difficulty = 2 // Number of leading zeros required in the proof-of-work hash
)

// generateChallenge generates a random challenge for the client
func generateChallenge() string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}

// isValidProof checks if the proof-of-work is valid
func isValidProof(challenge, response string) bool {
	targetPrefix := strings.Repeat("0", difficulty)
	combined := challenge + response
	hash := sha256.Sum256([]byte(combined))
	return strings.HasPrefix(hex.EncodeToString(hash[:]), targetPrefix)
}

// solveProofOfWork finds a valid proof-of-work response
func solveProofOfWork(challenge string) string {
	var nonce uint64
	for {
		nonceStr := strconv.FormatUint(nonce, 10)
		if isValidProof(challenge, nonceStr) {
			return nonceStr
		}
		nonce++
	}
}

var wordsOfWisdom = []string{
	"\"Wisdom is not a product of schooling but of the lifelong attempt to acquire it.\" - Albert Einstein",
	"\"The only true wisdom is in knowing you know nothing.\" - Socrates",
	"\"The wise man doesn't give the right answers, he poses the right questions.\" - Claude Levi-Strauss",
}

func getWordOfWisdom() (string, error) {
	if len(wordsOfWisdom) == 0 {
		return "", fmt.Errorf("empty string slice")
	}

	mathrand.Seed(time.Now().UnixNano())
	randomIndex := mathrand.Intn(len(wordsOfWisdom))
	return wordsOfWisdom[randomIndex], nil
}

// handleConnection handles the client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	challenge := generateChallenge()
	conn.Write([]byte(challenge + ":" + fmt.Sprint(difficulty) + "\n"))

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		response := scanner.Text()
		if isValidProof(challenge, response) {
			// Proof-of-work is valid, allow the connection
			wisdom, err := getWordOfWisdom()
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
			}
			conn.Write([]byte(wisdom + "\n"))
		} else {
			// Invalid proof-of-work, reject the connection
			conn.Write([]byte("Invalid proof-of-work. Connection rejected.\n"))
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Word of Wisdom TCP Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
