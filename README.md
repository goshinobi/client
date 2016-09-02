# client
net/http wrapper

# example

```
c := client.New(true)

proxyURL, err := url.Parse(proxyString)
transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
client := &http.Client{Transport: transport}
c.Add(client, nil)

fmt.Println(c.GetIPs())
fmt.Println(c.Get("http://google.co.jp"))
```


