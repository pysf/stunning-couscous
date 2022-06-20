#!/usr/bin/env bash

cd cmd/partnerseeder
go mod download 
go build 
chmod a+x ./cmd/partnerseeder/partnerseeder

echo "########### HI THERE ###########"
echo ""
echo "=> To run the cli app, run: "
echo "./cmd/partnerseeder/partnerseeder"
echo ""
echo "############## END #############"