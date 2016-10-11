#!/bin/sh
FILE_GLOB=/home/siuyin/go/src/siuyin/dra_processed/sample.log \
NUM_FILES=2 \
NUM_DAYS=1 \
REPORT_LEVELS=INFO,WARN,ERROR \
SCAN_TEXT=PullList,ChangeList,work_request \
	go run main.go
