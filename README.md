CEDEC2019のセッション情報JSONをiCal形式にするやつ
=================================================

# 使い方

* Goをインストール
* CEDEC2019公式サイトからダウンロードしたjson_file.jsonファイルを同じディレクトリに置く
* `go run main.go > calendar.ics 2>&1`を実行
* calendar.icsファイルが出来るので、Google Calendarにインポートする
* Linuxでしか動作確認していないのでWindowsとかだと改行文字まわりで問題が出るかもしれません。修正の必要あればIssue出すかPullRequestください

# LICENSE

MIT
