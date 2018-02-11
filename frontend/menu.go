package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"../block"
	"../network"
	"../pool"
)

func main() {
	// Initialization
	pool.GenesisPool()
	network.startServer(make(chan bool))

	offset := true
	for offset {
		fmt.Println("")
		fmt.Println("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
		fmt.Printf("★★★★★★ Welcome to Forest! You can begin your chatting now ★★★★★★\n")
		fmt.Println("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
		fmt.Printf("1. Write new message;\n")
		fmt.Printf("2. Store a public key;\n")
		fmt.Printf("3. Exit Forest. \n")

		//get input order from user
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("★★★(Enter any char of options above)\n>>")
		order, _ := reader.ReadString('\n')
		order = strings.TrimSuffix(order, "\n")
		fmt.Print(order)
		fmt.Println(reflect.TypeOf(order))
		orderInt, _ := strconv.Atoi(order)
		//fmt.Println(reflect.TypeOf(order))

		switch orderInt {
		case 1:
			fmt.Printf("You entered 1\n")
			newM()

		case 2:
			reader2 := bufio.NewReader(os.Stdin)
			fmt.Printf("You entered 2\n")
			fmt.Print("Enter the public key you want to store: \n>>")
			pubKey, _ := reader2.ReadString('\n')
			pubKey = strings.TrimSuffix(pubKey, "\n")
			fmt.Print("Enter the name of this public key: \n>>")
			uName, _ := reader2.ReadString('\n')
			fmt.Print("The public key you entered is: " + pubKey + ", the name you entered is: " + uName + "\n")
			storeAPublicKey(pubKey, uName)

		case 3:
			fmt.Printf("★Thank you for using Forest, you already exited ★\n")
			offset = false
		}
	}
}

/*To-Do: write a new message*/
func newM() {
	fmt.Println("func!!!!!!!!!!!")

	// Fetch message from user, choose public key, read public key into pubkey
	blk := block.CreateBlock(message, pubkey)
	network.forwardBlock(blk)
}

/* Used by storeAPublicKey function */
var (
	fileInfo *os.FileInfo
	err      error
)

/*To-Do: storeAPublicKey is storing public_key-user_Name pairs into a txt file,
so that newM() function could use to select receivers from one of the list.*/
func storeAPublicKey(publicKey string, userName string) {
	/*CHECK if the test.txt is already existed, if not, create the file. */
	fileInfo, err := os.Stat("test.txt")
	if err != nil {
		newFile, err := os.Create("test.txt")
		if err != nil {
			fmt.Println(err)
		}
		newFile.Close()
	} else {
		fmt.Println("File does exist. File information:")
		fmt.Println(fileInfo)
	}

	/*write public key and user name into the file. */
	file, err := os.OpenFile(
		"test.txt",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	li, err := file.WriteString(publicKey + "," + userName + "\n")
	fmt.Printf("wrote %d into the file.\n", li)

}
