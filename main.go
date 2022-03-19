package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
	"time"
)

var Cache = make(map[string]Orders) // map для хранения кеша

type Delivery struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction"`
	RequestId    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Items struct {
	ChrtId      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	Rid         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmId        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
}

type Orders struct {
	OrderUid          string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Items   `json:"items"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerId        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	Shardkey          string    `json:"shardkey" db:"shardkey"`
	SmId              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

func WriteCache() { // функция записи данных в кеш из бд
	db, err := sqlx.Open("postgres", "user=postgres password=b1O6ZPqsX7 dbname=wb_L0 sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// чтение данных из таблиц бд в переменные
	order := []Orders{}
	delivery := []Delivery{}
	payment := []Payment{}
	item := []Items{}
	err = db.Select(&order, `select order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, 
       shardkey, sm_id, date_created, oof_shard from orders`)
	if err != nil {
		log.Println(err.Error())
	}
	err = db.Select(&delivery, `select name, phone, zip, city, address, region, email from deliveries;`)
	if err != nil {
		log.Println(err.Error())
	}
	err = db.Select(&payment, `select transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, 
       goods_total, custom_fee from payments;`)
	if err != nil {
		log.Println(err.Error())
	}
	err = db.Select(&item, `select chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status from items;`)
	if err != nil {
		log.Println(err.Error())
	}

	// создание массива с order_uid для каждого элемента items
	var itemIds []string
	rows, err := db.Query(`select order_uid from items`)
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var itemId string
		err = rows.Scan(&itemId)
		if err != nil {
			log.Println(err.Error())
		}
		itemIds = append(itemIds, itemId)
	}

	j := 0
	for i := 0; i < len(order); i++ { // в структуру Orders записываем остальные структуры
		order[i].Delivery = delivery[i]
		order[i].Payment = payment[i]
		for j < len(itemIds) {
			var itemStruct []Items
			if itemIds[j] == order[i].OrderUid { // сравниваем order_uid из order и item, в случае совпадения записываем элемент item в массив структур Items
				itemStruct = append(itemStruct, item[j])
				j++ //увеличиваем счетчик, чтобы не начинать каждый раз с начала массива
			} else {
				order[i].Items = itemStruct // Как только order_uid из двух структур перестают совпадать, записываем полученный массив Items в order
				break
			}
		}
		Cache[order[i].OrderUid] = order[i] // полученный элемент order записываем в map с ключом order_uid
	}
}

func ReadFromChannel() { // функция чтения из канала
	sc, _ := stan.Connect("test-cluster", "ClientID")
	_, err := sc.Subscribe("foo", func(msg *stan.Msg) { // подписываемся на канал и при получении сообщения вызываем WriteData
		err := WriteData(msg)
		if err == nil {
			log.Println("Received a message")
		} else {
			log.Println(err.Error())
		}
	})
	if err != nil {
		log.Println(err.Error())
	}

}

func WriteData(m *stan.Msg) error { // функция записи данных в кеш и бд
	// сначала идет запись из json в кеш
	var order Orders
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		return err
	}
	if order.OrderUid == "" {
		return fmt.Errorf("Null OrderUid")
	}
	Cache[order.OrderUid] = order

	// затем из кеша данные переносятся в бд
	db, err := sql.Open("postgres", "user=postgres password=b1O6ZPqsX7 dbname=wb_L0 sslmode=disable")
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec(`insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, 
        sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService,
		order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	if err != nil {
		delete(Cache, order.OrderUid)
		db.Exec(`delete from orders where order_uid = $1`, order.OrderUid)
		return fmt.Errorf("Insert into orders", err)
	}
	_, err = db.Exec(`insert into deliveries (customer_id, name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7, $8);`,
		order.CustomerId, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		delete(Cache, order.OrderUid)
		db.Exec(`delete from orders where order_uid = $1`, order.OrderUid)
		return fmt.Errorf("Insert into deliveries", err)
	}
	_, err = db.Exec(`insert into payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, 
        goods_total, custom_fee) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		order.OrderUid, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		delete(Cache, order.OrderUid)
		db.Exec(`delete from orders where order_uid = $1`, order.OrderUid)
		return fmt.Errorf("Insert into payments", err)
	}
	for i := 0; i < len(order.Items); i++ {
		_, err = db.Exec(`insert into items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUid, order.Items[i].ChrtId, order.Items[i].TrackNumber, order.Items[i].Price, order.Items[i].Rid,
			order.Items[i].Name, order.Items[i].Sale, order.Items[i].Size, order.Items[i].TotalPrice, order.Items[i].NmId, order.Items[i].Brand, order.Items[i].Status)
		if err != nil {
			delete(Cache, order.OrderUid)
			db.Exec(`delete from orders where order_uid = $1`, order.OrderUid)
			return fmt.Errorf("Insert into items", err)
		}
	}
	return err
}

func HomePage(w http.ResponseWriter, r *http.Request) { // домашняя страница
	tmpl, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}
}

func IdPage(w http.ResponseWriter, r *http.Request) { // страница с записью с заданным id
	needId := r.URL.Query().Get("id")
	if _, ok := Cache[needId]; ok {
		b, _ := json.Marshal(Cache[needId])
		_, err := w.Write(b)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		_, err := w.Write([]byte("Запись не найдена"))
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func DataListPage(w http.ResponseWriter, r *http.Request) { // страница вывода списка всех записей
	outputArray := make([]Orders, 0)
	for _, elem := range Cache {
		outputArray = append(outputArray, elem)
	}

	b, _ := json.Marshal(outputArray)
	_, err := w.Write(b)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomePage)
	mux.HandleFunc("/record", IdPage)
	mux.HandleFunc("/list/", DataListPage)

	WriteCache()

	go ReadFromChannel()

	log.Println("Запуск сервера...")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
