# vcoin-third-api
透過WebSocket監聽BTCUSDT最新交易的資訊

### 啟動服務
```
go run main.go server [-p = port] 
```
預設 port = 3000

### 透過訪問取得最新交易資訊
訪問 http://localhost:3000/api/record/newest
### 透過監聽取的最新交易資訊
透過此[網站測試](http://coolaf.com/tool/chattest)
輸入 ws://localhost:3000/api/record/observe 按下連接即可獲得最新的BTCUSDT交易資訊

