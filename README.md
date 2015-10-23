sns-subscriber
====
## Description
Amazon SNS トピックに大量のエンドポイントを追加する

## Usage
1. 改行区切りで複数のEndpoint Arnが入ったテキストファイルを用意
2. sns-subscriber -s endpoints.txt -t arn:aws:sns:xxx:xxxx -r ap-northeast-1
