#!/usr/bin/env bash
go mod download
go build -o ./cmd/stunning-couscous

cd cmd/partnerseeder
go build -o ../seeder


echo "########### HI THERE ###########"
echo ""
echo "=> To run the cli app, run: "
echo "./cmd/seeder"
echo ""
echo "############## END #############"