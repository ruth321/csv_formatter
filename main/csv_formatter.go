package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type BQOrderRaw struct {
	//Common data
	UUID               string    `json:"uuid"`
	RoutesCount        int       `json:"rotes_count"`          //количество точек
	ServiceName        string    `json:"service_name"`         //название услуги
	Features           string    `json:"features"`             //фичи
	CreatedDatetime    time.Time `json:"created_datetime"`     //вермя создания
	Source             string    `json:"source"`               //источник заказа
	OrderState         string    `json:"order_state"`          //тип завершения заказа
	CancelReason       string    `json:"cancel_reason"`        //причина отмены
	OrderTakenTime     time.Time `json:"arrival_real_time"`    //время когда заказ был взят
	ArrivalTimeReal    time.Time `json:"arrival_real_time"`    //время прибытия по факту
	ArrivalTimePromise time.Time `json:"arrival_promise_time"` //обещанное время (сек)
	DistanceToClient   float64   `json:"distance_to_client"`   //растояние до клиента
	TripDistance       float64   `json:"trip_distance"`        //растояние поездки
	CancelTime         time.Time `json:"cancel_time"`          //
	OwnerUUID          string    `json:"owner_uuid"`           //идентификатор оунера
	UserStartName      string    `json:"user_start_name"`      //какой пользователь запустил заказа
	UserStartUUID      string    `json:"user_start_uuid"`      //тип оплаты

	//pickup data
	PickupLon      float32   `json:"pickup_lon"`
	PickupLat      float32   `json:"pickup_lat"`
	PickupDatetime time.Time `json:"pickup_datetime"` // во сколько началась поездка
	PickupArea     string    `json:"pickup_area"`     // район откуда забрали
	PickupAddress  string    `json:"pickup_address"`  //адрес

	//dropoff data
	DropoffLon      float32   `json:"dropoff_lon"`
	DropoffLat      float32   `json:"dropoff_lat"`
	DropoffDatetime time.Time `json:"dropoff_datetime"` // откуда
	DropoffArea     string    `json:"dropoff_area"`
	DropoffAddress  string    `json:"dropoff_address"`

	//tariff data
	TariffName       string `json:"tariff_name"`       //выбранный тариф
	TariffPrice      int    `json:"tariff_price"`      //цена по тарифу
	RealPrice        int    `json:"real_price"`        //цена по факту
	WaitingTime      int    `json:"waiting_time"`      //ожидание (сек)
	WaitingPrice     int    `json:"waiting_price"`     //стоимость ожидания
	BonusPayment     int    `json:"bonus_payment"`     //оплата бонусами
	GuaranteedIncome int    `json:"guaranteed_income"` //гарантированный доход
	ClientAllowance  int    `json:"client_allowance"`  //клиентская надбавка

	//driver data
	DriverUUID   string `json:"driver_uuid"`   //id водителя
	DriverCar    string `json:"driver_car"`    //машина водителя
	DriverTarrif string `json:"driver_tarrif"` //тариф водителя

	//client data
	ClientPhone string `json:"client_phone"` //номер клиента
	ClientUUID  string `json:"client_uuid"`  //
	PaymentType string `json:"payment_type"` //тип оплаты

	//products data
	StoreUUID     string `json:"store_uuid"`
	ProductsSum   int    `json:"product_sum"`
	ProductsCount int    `json:"product_count"`
	ProductsData  string `json:"product_data"`

	//insert datetime
	InsertDateTime time.Time `json:"insert_datetime"` // время вставки

	//order events list
	Events []OrderEventData `json:"events"`
}

type OrderEventData struct {
	EventTime    time.Time `json:"event_time"`
	DriverUUID   string    `json:"driver_uuid"`
	OrderUUID    string    `json:"order_uuid"`
	Publisher    string    `json:"publisher"`
	OperatorUUID string    `json:"operator_uuid"`
	State        string    `json:"state"`
	Comment      string    `json:"comment"`

	//insert datetime
	InsertDateTime time.Time `json:"insert_datetime"` // время вставки
}

func main() {
	file, err := os.Open("main.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var order BQOrderRaw
	layout := "2006-01-02 15:04:05.999999-07"
	for i := 1; i < len(records); i++ {
		createdDatetime, err := time.Parse(layout, records[i][5])
		if err != nil {
			logrus.WithFields(logrus.Fields{"order field": "CreatedDatetime", "event": "time parsing", "csv field": "createtime"}).Fatal(err)
		}
		dropoffLon, err := strconv.ParseFloat(records[i][60], 32)
		if err != nil {
			logrus.WithFields(logrus.Fields{"order field": "DropoffLon", "event": "float parsing", "csv field": "longitudeto"}).Fatal(err)
		}
		dropoffLat, err := strconv.ParseFloat(records[i][59], 32)
		if err != nil {
			logrus.WithFields(logrus.Fields{"order field": "DropoffLat", "event": "float parsing", "csv field": "latitudeto"}).Fatal(err)
		}
		pickupLon, err := strconv.ParseFloat(records[i][27], 32)
		if err != nil {
			logrus.WithFields(logrus.Fields{"order field": "PickupLon", "event": "float parsing", "csv field": "longitude"}).Fatal(err)
		}
		pickupLat, err := strconv.ParseFloat(records[i][26], 32)
		if err != nil {
			logrus.WithFields(logrus.Fields{"order field": "PickupLat", "event": "float parsing", "csv field": "latitude"}).Fatal(err)
		}
		order = BQOrderRaw{
			UUID:               records[i][0], //idx
			RoutesCount:        0,
			ServiceName:        "",
			Features:           "",
			CreatedDatetime:    createdDatetime, //createtime
			Source:             "",
			OrderState:         records[i][50], //state
			CancelReason:       "",
			OrderTakenTime:     time.Time{},
			ArrivalTimeReal:    time.Time{},
			ArrivalTimePromise: time.Time{},
			DistanceToClient:   0,
			TripDistance:       0,
			CancelTime:         time.Time{},
			OwnerUUID:          "",
			UserStartName:      "",
			UserStartUUID:      "",
			PickupLon:          float32(pickupLon), //longitude
			PickupLat:          float32(pickupLat), //latitude
			PickupDatetime:     time.Time{},
			PickupArea:         "",
			PickupAddress:      records[i][2],       //addressfrom
			DropoffLon:         float32(dropoffLon), //longitudeto
			DropoffLat:         float32(dropoffLat), //latitudeto
			DropoffDatetime:    time.Time{},
			DropoffArea:        "",
			DropoffAddress:     records[i][4], //addressto
			TariffName:         "",
			TariffPrice:        0,
			RealPrice:          0,
			WaitingTime:        0,
			WaitingPrice:       0,
			BonusPayment:       0,
			GuaranteedIncome:   0,
			ClientAllowance:    0,
			DriverUUID:         "",
			DriverCar:          "",
			DriverTarrif:       "",
			ClientPhone:        records[i][24], //aclientphone
			ClientUUID:         "",
			PaymentType:        "",
			StoreUUID:          "",
			ProductsSum:        0,
			ProductsCount:      0,
			ProductsData:       "",
			InsertDateTime:     time.Time{},
			Events:             nil,
		}
		ordersDirPath := "orders"
		createdYear := strconv.Itoa(createdDatetime.Year())
		createdMonth := strconv.Itoa(int(createdDatetime.Month()))
		createdDay := strconv.Itoa(createdDatetime.Day())
		savingPath := fmt.Sprintf("%s/%s/%s/%s/", ordersDirPath, createdYear, createdMonth, createdDay)
		err = os.MkdirAll(savingPath, os.ModePerm)
		if err != nil {
			logrus.WithField("event", "making directory").Fatal(err)
		}
		orderFile, err := json.Marshal(order)
		if err != nil {
			logrus.WithField("event", "encoding json").Fatal(err)
		}
		fileName := order.UUID + ".json"
		err = ioutil.WriteFile(savingPath+fileName, orderFile, os.ModePerm)
		if err != nil {
			logrus.WithField("event", "saving order file").Fatal(err)
		}
	}
}

func checkErr(err error, orderField string, event string, csvField string) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"order field": orderField, "event": event, "csv field": csvField}).Fatal(err)
	}
}
