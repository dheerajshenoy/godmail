package mail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"encoding/json"
)

var PACKAGE_NAME string = "godmail"

// This function sends the mail
func sendMail(fromAddr string, password string, toAddr[] string, subject string, body bytes.Buffer, files[] string) (bool, error) {
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    // Set up authentication information.
    auth := smtp.PlainAuth("", fromAddr, password, smtpHost)

    // Create the body of the message
	writer := multipart.NewWriter(&body)

	// If body is empty, make it empty
	if body.String() == "" {
		body.WriteString("")
	}

	// Add email text body
	textWriter, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	textWriter.Write([]byte(body.String()))

    // Add attachment

	for i := 0; i < len(files); i++ {
		attachFile(writer, files[i])
	}

    writer.Close()

    // Set headers
    headers := make(map[string]string)
    headers["From"] = fromAddr
    headers["To"] = toAddr[0]
    headers["Subject"] = subject
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary())

    // Combine headers and body
    var msg bytes.Buffer
    for k, v := range headers {
        msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
    }
    msg.WriteString("\r\n")

	msg.Write(body.Bytes())

    // Send email
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromAddr, toAddr, msg.Bytes())

    if err != nil {
        fmt.Println("Error:", err)
		return false, err
	}

	return true, nil
}

// Function to attach a file
func attachFile(writer *multipart.Writer, filePath string) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    // Create attachment header
    part, err := writer.CreatePart(map[string][]string{
        "Content-Disposition": {fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath))},
        "Content-Type":        {"application/octet-stream"},
        "Content-Transfer-Encoding": {"base64"},
    })

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Read file and write base64 encoded content
    buf := make([]byte, 3*1024) // Read in chunks of 3KB
    for {
        n, err := file.Read(buf)
        if err != nil {
            break
        }
        encoded := make([]byte, base64.StdEncoding.EncodedLen(n))
        base64.StdEncoding.Encode(encoded, buf[:n])
        part.Write(encoded)
    }
}

// Function to display details provided to send mail and ask for confirmation from the user
func confirmSendingMail(fromAddr string, toAddr[] string, password string, subject string, body string, files[] string) {

    // Create a reader for standard input
    reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nSender Email Address: ", fromAddr)
	fmt.Println("Receipient Email Address: ", strings.Join(toAddr, ","))
	fmt.Println("Body: ", body)
	fmt.Println("Files: ", strings.Join(files, ","))

    // Prompt the user for input
    fmt.Print("Is this okay (y/n): ")

	// Read the input from the user
	confirmation, err := reader.ReadString('\n') // Read until a newline character

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

    confirmation = strings.TrimSpace(strings.ToLower(confirmation))
	if confirmation == "y" {
		var bodyBytesBuffer bytes.Buffer
		bodyBytesBuffer.WriteString(body)
		_, err := sendMail(fromAddr, password, toAddr, subject, bodyBytesBuffer, files)
		if err != nil {
			fmt.Println("Error while sending mail", err)
		}
		fmt.Println("Email Sent Successfully")
	}

}

type Config struct {
    Password    string        `json:"password"`
}

func pathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

// Function to handle and parse arguments passed to the program
func Mail() {

	fromAddrFlag := flag.String("from", "", "Senders email address")
	toAddrFlag := flag.String("to", "", "Receipient email address or addresses")
	subjectFlag := flag.String("subject", "", "Subject of the email")
	bodyFlag := flag.String("body", "", "Body of the email")
	_ = flag.String("editor", "false", "Open editor for writing the body of the email")
	filesFlag := flag.String("files", "", "Files to attach")
	passwordFlag := flag.String("password", "", "Senders email password (Use this only as a last resort)")

	flag.Parse()

	// Get list of all flags specified on command line
	specifiedFlags := make(map[string]bool)
	flag.CommandLine.Visit(func(f *flag.Flag) {
		specifiedFlags[f.Name] = true
	})

	// Read config file, if it exists
	var conf Config

	configDir, err := os.UserConfigDir()

	if err == nil {
		configFile := filepath.Join(configDir, PACKAGE_NAME, "config.json")
		exists, err := pathExists(configFile)

		if exists && err == nil {

			file, err := os.Open(configFile)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer file.Close()

			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&conf); err != nil {
				fmt.Println("Error decoding JSON:", err)
				return
			}
		}
	}

	if *fromAddrFlag == "" {
		fmt.Println("Error: -name flag is required")
		os.Exit(-1)
	}

	if *toAddrFlag == "" {
		fmt.Println("Error: -to flag is required")
		os.Exit(-1)
	}

	if *subjectFlag == "" {
		fmt.Println("No subject found, using default subject value")
		*subjectFlag = "Subject"
	}

	// TODO: Use vim/nano/something else for asking body of the email
	if specifiedFlags["editor"] {
		
	}

	password := ""

	if *passwordFlag == "" {
		password = conf.Password
	} else {
		password = *passwordFlag
	}

	files := []string{}

	if *filesFlag != "" {
		files = strings.Split(*filesFlag, ",")
	} 
	// confirmSendingMail(*fromAddrFlag, [] string { *toAddrFlag } , *passwordFlag, *bodyFlag, [] string {*filesFlag})

	confirmSendingMail(*fromAddrFlag, [] string { *toAddrFlag } , password, *subjectFlag, *bodyFlag, files)
}
