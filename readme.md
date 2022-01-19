# 設計

### 最小限度的使用 RPC
相同的資料只取一次，block 持續收進 db，transaction 則是有需求時才從chain取得，並存進 db 以供日後使用

### 保持單純
indexer 定時由 chain 取得 block 並收進 db，
在由 API 取得 block 時完全以 db 內的資料為準

好處是可以讓架構較為簡單，並減少 rpc 的使用

缺點是，跟直接從 endpoint 取資料相比，透過取得的區塊會有時間差
取得的間隔設定為10秒，會產出3 ~ 4 個新的 block，也就是 api 會落後於 endpoint 的幅度


# 架構

以 MySQL 作為資料庫。

執行主體為 indexer ，以 Worker Pool pattern 用複數個 worker 來同時取得鏈上的資料並儲存 , 
worker 的數量會影響資源的使用，也要考量 endpoint 的 rate limit ，開發時使用 3 個 worker 同時處理

另外再起 apiservice ，負責處理 RESTful API 的需求

## indexer
在背景執行，持續將鏈上的資料取回，分成 fetch worker 及 confirm worker

其中 fetch worker 被設計成會單純的將 block 掃進 db

而 confirm worker 會檢查已經取得的區塊，看是否有需要替換成穩定區塊，如有則進行替換。


# 執行

WIP