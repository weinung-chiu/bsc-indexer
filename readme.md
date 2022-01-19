# 設計

### 最小限度的使用 RPC
相同的資料只取一次，block 持續收進 db，transaction 則是有需求時才從chain取得，並存進 db 以供日後使用

### 保持單純
indexer 定時由 chain 取得 block 並收進 db， 呼叫 API 取得 block 時完全以 db 內的資料為準

好處是可以讓架構較為簡單，並減少 rpc 的使用

缺點是，跟直接從 endpoint 取資料相比，透過取得的區塊會有時間差

取得的間隔設定為10秒，期間會產出3 ~ 4 個新的 block，也就是 api 會落後於 endpoint 的幅度


# 架構

以 MySQL 作為資料庫。

執行主體為 indexer ，以 Worker Pool pattern 用複數個 worker 來同時取得鏈上的資料並儲存 , 
worker 的數量會影響資源的使用，也要考量 endpoint 的 rate limit ，開發時使用 3 個 worker 同時處理

另外再起 apiservice ，負責處理 RESTful API 的需求

## indexer
在背景執行，持續將鏈上的資料取回，分成 fetch worker 及 confirm worker

其中 fetch worker 被設計成會單純的將 block 掃進 db

而 confirm worker 會檢查已經取得的區塊，看是否有需要替換成穩定區塊，如有則進行替換。

在開發時，避免從空 db 開始執行時需要花費很久才能取到最新block的狀況，indexer 會從鏈的中間開始掃到最新 block (而非從0開始)，
這個範圍可以透過 IndexLimit 進行調整。

# 執行

MySQL db
```shell
docker compose up
```


db migrate
```shell
go run cmd/db_migrate/main.go
```


indexer and RESTful API Service
```shell
go run cmd/dev/main.go
```

# APIs
```shell
GET 127.0.0.1/blocks

GET 127.0.0.1/blocks/:id

GET 127.0.0.1/transaction/:hash
```