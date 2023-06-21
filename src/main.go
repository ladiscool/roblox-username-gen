package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Username struct {
	ThreadCount     int
	CharCount       int
	PathToProxyFile string
}

type Response struct {
	Message string `json:"message"`
}

func (u *Username) GenProxy() string {
	proxyFile, err := os.ReadFile(u.PathToProxyFile)
	if err != nil {
		panic(err)
	}

	proxies := strings.Split(string(proxyFile), "\n")
	return proxies[rand.Intn(len(proxies))]
}

func sanitizeProxyURL(proxyURL string) string {
	proxyURL = strings.TrimSpace(proxyURL)
	proxyURL = strings.ReplaceAll(proxyURL, "\r", "")
	return proxyURL
}

func (u *Username) UserCheck(wg *sync.WaitGroup) {
	defer wg.Done()

	proxyURL := u.GenProxy()
	proxyURL = sanitizeProxyURL(proxyURL)

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxyURL)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
		Timeout: 5 * time.Second,
	}

	letters := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"

	for {
		username := generateRandomString(letters+digits, u.CharCount)

		resp, err := client.Get(fmt.Sprintf("https://auth.roblox.com/v1/usernames/validate?Username=%s&Birthday=20%%2C%%20May%%2C%%202000&Context=0", username))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var res Response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			err = json.Unmarshal(body, &res)
			if err != nil {
				continue
			}

			if res.Message == "Username is valid" {
				fmt.Printf("(\033[92m+\033[0m) Valid -> \033[92m%s\033[0m\n", username)
				file, err := os.OpenFile("./cogs/usernames.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					continue
				}

				if _, err := file.WriteString(username + "\n"); err != nil {
					file.Close()
					continue
				}

				file.Close()
			} else {
				fmt.Printf("(\033[91m-\033[0m) Invalid -> \033[91m%s\033[0m\n", username)
			}
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func (u *Username) Start() {
	var wg sync.WaitGroup

	for i := 0; i < u.ThreadCount; i++ {
		wg.Add(1)
		go u.UserCheck(&wg)
	}

	wg.Wait()
}

func generateRandomString(charset string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	threadCount, charCount := 0, 0
    
        clearScreen()
    
	fmt.Print("(\033[95m?\033[0m) Threads -> ")
	fmt.Scanln(&threadCount)
	fmt.Print("(\033[95m?\033[0m) Character Size -> ")
	fmt.Scanln(&charCount)

	clearScreen()

	username := Username{
		ThreadCount:     threadCount,
		CharCount:       charCount,
		PathToProxyFile: "./proxies.txt",
	}

	username.Start()
}
