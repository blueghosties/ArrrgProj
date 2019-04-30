package main

// TODO: send encrypted data to remote server and decrypt
// TODO: Need to have go code listen for browser window closing so the application can exit.

// COMPLETED: App now opens a browser and index.html locally
// This comes with no guarantees. This code is my experimentation and journey of learning of learning Go.
// Feel free to use the code as you see fit and understand I am not an encryption specialist...I am still
// learning. Enjoy. 2019

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/nacl/secretbox"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	fmt.Println("Let it run, baby!")

	http.HandleFunc("/", index)
	http.HandleFunc("/process", processor)
	exec.Command("open", "http://127.0.0.1:8081/templates/index.html").Run()
	http.ListenAndServe(":8081", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "index.html", "ArrrgProj")
}

func processor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	dMsg := ""
	fMsg := r.FormValue("pmsg")
	//msg := EncIt(fMsg)
	msg := ""
	fdemsg := r.FormValue("demsg")

	// If the textboxes are empty, don't process anything
	// else execute functions
	if fMsg == "" {

	} else {
		msg = EncIt(fMsg)
	}

	if fdemsg == "" {

	} else {
		dMsg = DecIt(fdemsg)
	}

	// write the encrypted msg to file
	writeTxt(msg + "\n")

	d := struct {
		Emsg string
		Dmsg string
	}{
		Emsg: msg,
		Dmsg: dMsg,
	}

	tpl.ExecuteTemplate(w, "index.html", d)

	// fmt.Fprint(w, "<input onclick='{{GoToExit}}' type='button' value='Close App' name='btnClose'>")

	// connit(msg)
}

func EncIt(encIt string) string {
	var nonce [24]byte
	var password [32]byte

	_, _ = io.ReadAtLeast(rand.Reader, nonce[:], 24)

	_, _ = io.ReadAtLeast(rand.Reader, password[:], 32)

	// fmt.Println("Before encryptions: ", encIt)
	isEnced := secretbox.Seal(nil, []byte(encIt), &nonce, &password)
	enHexd := fmt.Sprintf("%x%x%x", password, nonce[:], isEnced)
	fmt.Println("Encrypted: ", enHexd)

	// connit(enHexd)

	return enHexd
}

func writeTxt(msg string) {

	if _, err := os.Stat("msgs.txt"); err == nil {
		// path/to/whatever exists
		f, _ := os.OpenFile("msgs.txt", os.O_APPEND|os.O_WRONLY, 0644)
		f.WriteString(msg)
		f.Close()
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		f, _ := os.Create("msgs.txt")
		f.WriteString(msg)
		f.Close()
	} else {
		fmt.Println("something funky happened.")
	}

}

// DecIt this is a fucking comment for no reason whatsoever
// not sure why
func DecIt(encdStr string) string {
	var nonce2 [24]byte
	var password [32]byte
	var message []byte

	pw := encdStr[:64]
	nonce := encdStr[64:112]
	encMsg := encdStr[112:]

	// fmt.Println(string(encdStr[111]))
	//fmt.Println("pw: " + pw)
	//fmt.Println("nonce: " + nonce)
	//fmt.Println("encMsg: " + encMsg)

	// parts := strings.SplitN(encdStr, ":", 3)
	// if len(parts) < 2 {
	// 	fmt.Errorf("expected nonce?")
	// }

	bs, err := hex.DecodeString(pw)
	if err != nil || len(bs) != 32 {
		fmt.Errorf("invalid password")
	}

	copy(password[:], bs)

	bs, err = hex.DecodeString(nonce)
	if err != nil || len(bs) != 24 {
		fmt.Errorf("invalid nonce")
	}
	copy(nonce2[:], bs)

	bs, _ = hex.DecodeString(encMsg)

	copy(message[:], bs)

	msg, ok := secretbox.Open(nil, bs, &nonce2, &password)
	if !ok {
		fmt.Errorf("invalid decoded message")
	}
	fmt.Println("MESSAGE DECRYPTED AT SERVER: ", string(msg))

	// this is here for testing purposes
	// wToFile := string(msg)
	// writeTxt(wToFile)

	return string(msg)
}

//func connit(msg string) {
//	conn, _ := net.Dial("tcp", "127.0.0.1:3333")
//
//	_, _ = fmt.Fprintf(conn, msg +"\n")
//	// listen for reply
//	message, _ := bufio.NewReader(conn).ReadString('\n')
//
//
//	fmt.Print("Message from server: "+message)
//
//	conn.Close()
//}
