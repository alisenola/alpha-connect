package utils

import (
	"crypto/md5"
	"io"
	"os"
)

func FileToMD5(filePath string) ([]byte, error) {
	//Initialize variable returnMD5String now in case an error has to be returned

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	//Tell the program to call the following function when the current function returns

	//Open a new hash types to write to
	hash := md5.New()

	//Copy the file in the hash types and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	return hashInBytes, nil
}
