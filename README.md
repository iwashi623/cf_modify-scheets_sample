# cf_spreadsheets

## これは何？
Cloud Function上のGolangコードからスプレッドシートをいじろうとしたときのサンプルコード

## Deploy
```
 $ gcloud functions deploy iwashita_test --gen2 --entry-point modifySheets --runtime go122 --trigger-http --region asia-northeast1
```
