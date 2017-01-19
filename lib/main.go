package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

//MessagePayload payload structure to send to client
type MessagePayload struct {
	FileContent  []byte
	Lastmodified string
}

//FileWatcher watches given file and send change on browser
func FileWatcher(filename string, conn *websocket.Conn) {
	lastModTime := time.Now()

	fileWatch, openErr := os.Open(filename)
	if openErr != nil {
		log.Fatal("Error", openErr)
	}
	fmt.Println("file to watch opened")
	fileInfo, statErr := fileWatch.Stat()
	if statErr != nil {
		log.Fatal("Error", statErr)
	}
	actualLastModTime := fileInfo.ModTime()
	UpdateBrowser(filename, actualLastModTime, conn)
	for {
		time.Sleep(2 * time.Second)
		fileInfo, statErr = fileWatch.Stat()
		if statErr != nil {
			log.Fatal("Error", statErr)
		}
		actualLastModTime := fileInfo.ModTime()
		if lastModTime.Before(actualLastModTime) {
			fmt.Println("File has been modified", actualLastModTime)
			UpdateBrowser(filename, actualLastModTime, conn)
			lastModTime = actualLastModTime
		} else {
			fmt.Println("-------not touched----------")
		}
	}
}

//UpdateBrowser take file pointer as argument
// and print out file content
func UpdateBrowser(f string, lastmodtime time.Time, conn *websocket.Conn) {
	fileData, readErr := ioutil.ReadFile(f)
	if readErr != nil {
		log.Fatal("Error", readErr)
	}
	newPayload := MessagePayload{fileData, lastmodtime.String()}
	if err := conn.WriteJSON(newPayload); err != nil {
		log.Fatal("Error", err)
	}
}
