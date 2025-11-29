package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

// HashPassword hashes a plain text password using Argon2id.
// Argon2id is a memory-hard password hashing algorithm that provides
// protection against GPU and side-channel attacks.
//
// Returns the hashed password or an error if hashing fails.
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return hash, err
}

// CheckPasswordHash compares a plain text password with a hashed password.
//
// Returns true if the password matches the hash, false otherwise.
// Returns an error if the comparison process fails.
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
		return match, err
	}

	return match, err
}
