#!/bin/sh
FILE_GLOB=/home/siuyin/go/src/siuyin/dra_processed/sample.log \
NUM_FILES=3 \
NUM_DAYS=3 \
REPORT_LEVELS=INFO,WARN,ERROR \
DATE_REGEXP='(\d{8})-\d{6}:' \
SCAN_TEXT=PullList,ChangeList,work_request \
	go run main.go
