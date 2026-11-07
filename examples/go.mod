module github.com/andywang2011/go-alipay/examples

go 1.27.3

require github.com/andywang2011/go-alipay v0.0.0

require (
	github.com/andywang2011/go-httplib v0.1.0 // indirect
	github.com/andywang2011/go-logging v0.1.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace (
	github.com/andywang2011/go-httplib => ../../go-httplib
	github.com/andywang2011/go-logging => ../../go-logging
)

replace github.com/andywang2011/go-alipay => ..
