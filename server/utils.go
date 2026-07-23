package main

import "math/rand/v2"


func getRandomNickname() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	nickname := make([]byte, 10)
	for i := range nickname {
		nickname[i] = letterBytes[rand.IntN(len(letterBytes))] 
	}

	return string(nickname)
}


func getRandomId() int {
	return rand.IntN(1000000000)
}

