package main

import (
	"github.com/nats-io/stan.go"
	"log"
)

var message1 = []byte("{\n  \"order_uid\": \"b563feb7b2b84b6test\",\n  \"track_number\": \"WBILMTESTTRACK\",\n  " +
	"\"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Test Testov\",\n    \"phone\": \"+9720000000\",\n    " +
	"\"zip\": \"2639809\",\n    \"city\": \"Moscow\",\n    \"address\": \"Lenina 2\",\n    \"region\": " +
	"\"Moscow\",\n    \"email\": \"test1@gmail.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"b563feb7b2b84b6test\",\n    " +
	"\"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"wbpay\",\n    \"amount\": 1817,\n    " +
	"\"payment_dt\": 1637907727,\n    \"bank\": \"alpha\",\n    \"delivery_cost\": 1500,\n    \"goods_total\": 317,\n    " +
	"\"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 9934930,\n      \"track_number\": \"WBILMTESTTRACK\",\n" +
	"\"price\": 453,\n      \"rid\": \"ab4219087a764ae0btest\",\n      \"name\": \"Mascaras\",\n      \"sale\": 30,\n" +
	"\"size\": \"0\",\n      \"total_price\": 317,\n      \"nm_id\": 2389212,\n      \"brand\": \"Vivienne Sabo\",\n " +
	"\"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"test1\",\n  " +
	"\"delivery_service\": \"meest\",\n  \"shardkey\": \"9\",\n  \"sm_id\": 99,\n  \"date_created\": \"2021-11-26T06:22:19Z\",\n  " +
	"\"oof_shard\": \"1\"\n}")

var message2 = []byte("{\n  \"order_uid\": \"b563feb7b2b84b7test\",\n  \"track_number\": \"WBILMTESTTRACK\",\n  " +
	"\"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Oleg\",\n    \"phone\": \"+9720000000\",\n    " +
	"\"zip\": \"2639809\",\n    \"city\": \"Ryazan\",\n    \"address\": \"Stalina 1\",\n    \"region\": " +
	"\"Ryazan\",\n    \"email\": \"test2@gmail.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"b563feb7b2b84b6test\",\n    " +
	"\"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"wbpay\",\n    \"amount\": 1817,\n    " +
	"\"payment_dt\": 1637907727,\n    \"bank\": \"alpha\",\n    \"delivery_cost\": 1500,\n    \"goods_total\": 317,\n    " +
	"\"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 9934930,\n      \"track_number\": \"WBILMTESTTRACK\",\n" +
	"\"price\": 453,\n      \"rid\": \"ab4219087a764ae0btest\",\n      \"name\": \"Mascaras\",\n      \"sale\": 30,\n" +
	"\"size\": \"0\",\n      \"total_price\": 317,\n      \"nm_id\": 2389212,\n      \"brand\": \"Vivienne Sabo\",\n " +
	"\"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"test2\",\n  " +
	"\"delivery_service\": \"meest\",\n  \"shardkey\": \"9\",\n  \"sm_id\": 99,\n  \"date_created\": \"2021-11-26T06:22:19Z\",\n  " +
	"\"oof_shard\": \"1\"\n}")

var message3 = []byte("{\n  \"order_uid\": \"b563feb7b2b84b8test\",\n  \"track_number\": \"WBILMTESTTRACK\",\n  " +
	"\"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Egor\",\n    \"phone\": \"+9720000000\",\n    " +
	"\"zip\": \"2639809\",\n    \"city\": \"Saint P\",\n    \"address\": \"Nevskaya 13\",\n    \"region\": " +
	"\"Leningradskaya obl\",\n    \"email\": \"test3@gmail.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"b563feb7b2b84b6test\",\n    " +
	"\"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"wbpay\",\n    \"amount\": 1817,\n    " +
	"\"payment_dt\": 1637907727,\n    \"bank\": \"alpha\",\n    \"delivery_cost\": 1500,\n    \"goods_total\": 317,\n    " +
	"\"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 9934930,\n      \"track_number\": \"WBILMTESTTRACK\",\n" +
	"\"price\": 453,\n      \"rid\": \"ab4219087a764ae0btest\",\n      \"name\": \"Mascaras\",\n      \"sale\": 30,\n" +
	"\"size\": \"0\",\n      \"total_price\": 317,\n      \"nm_id\": 2389212,\n      \"brand\": \"Vivienne Sabo\",\n " +
	"\"status\": 202\n    },\n{\n      \"chrt_id\": 9934930,\n      \"track_number\": \"WBILMTEST2TRACK\",\n      \"price\": 500,\n      " +
	"\"rid\": \"ab4219087a764ae0btest\",\n      \"name\": \"Sneakers\",\n      \"sale\": 10,\n      \"size\": \"44\",\n      \"total_price\": 500,\n " +
	"     \"nm_id\": 2389212,\n      \"brand\": \"Adidas\",\n      \"status\": 202\n    } ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"test3\",\n  " +
	"\"delivery_service\": \"meest\",\n  \"shardkey\": \"9\",\n  \"sm_id\": 99,\n  \"date_created\": \"2021-11-26T06:22:19Z\",\n  " +
	"\"oof_shard\": \"1\"\n}")

var fakemessage1 = []byte("Hello world")

func main() {
	sc, err := stan.Connect("test-cluster", "testSendID")
	if err != nil {
		log.Println(err.Error())
	}

	sc.Publish("foo", message1)
	sc.Publish("foo", message2)
	sc.Publish("foo", message3)
	sc.Publish("foo", fakemessage1)
	sc.Close()
}
