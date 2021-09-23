package main

import (
    "fmt"
    "bufio"
    "net/url"
    "strings"
    "path"
    "net"
    "net/http"
    "os"
    "time"
    "sync"
    "crypto/tls"
)

func printResponseCode(client *http.Client, url string) {
    req, err := http.NewRequest("GET", url, nil)
    if err!= nil {
        fmt.Println("Couldn't get " + url)
        os.Exit(1)
    }

    req.Header.Add("Connection", "close")
    req.Close = true

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Couldn't get" + url)
        return
    }
    fmt.Printf("%s [%d]\n", url, resp.StatusCode)
}

func main() {

    timeout := time.Duration(10000000000)

    var tr = &http.Transport{
        MaxIdleConns: 30,
        IdleConnTimeout: time.Second,
        DisableKeepAlives: true,
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        DialContext: (&net.Dialer{
            Timeout: timeout,
            KeepAlive: time.Second,
        }).DialContext,
    }

    re := func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }

    client := &http.Client{
        Transport: tr,
        CheckRedirect: re,
        Timeout: timeout,
    }

    allURLs := make(chan string)

    var httpWG sync.WaitGroup
    for i := 0; i < 20; i++ {
        httpWG.Add(1)

        go func() {
            for url := range allURLs {
                printResponseCode(client, url)
                continue
            }
            httpWG.Done()
        }()
    }

    // read from stdin
    sc := bufio.NewScanner(os.Stdin)
    for sc.Scan() {
        domain := strings.ToLower(sc.Text())
        u, err := url.Parse(domain)
        if err != nil {
            fmt.Println("Couldn't parse")
            os.Exit(1)
        }
        split := strings.Split(u.Hostname(), ".")
        u.Path = path.Join(u.Path, split[0])
        allURLs <- domain
        allURLs <- u.String()
        allURLs <- u.String()+"/"
    }

    close(allURLs)
    httpWG.Wait()

}
