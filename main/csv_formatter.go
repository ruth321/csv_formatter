package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
	OrderTakenTime     time.Time `json:"order_taken_time"`     //время когда заказ был взят
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
		logrus.WithField("event", "csv file opening").Fatal(err)
	}
	r := csv.NewReader(file)
	csvOrders, err := r.ReadAll()
	if err != nil {
		logrus.WithField("event", "reading from csv reader").Fatal(err)
	}
	err = file.Close()
	if err != nil {
		logrus.WithField("event", "csv file closing").Fatal(err)
	}
	var order BQOrderRaw
	for i := 1; i < len(csvOrders); i++ {
		createdDatetime, err := TimeParser(csvOrders[i][5])
		errorHandler(err, "CreatedDatetime", "time parsing", "createtime")
		dropoffLon, err := strconv.ParseFloat(csvOrders[i][60], 32)
		errorHandler(err, "DropoffLon", "float parsing", "longitudeto")
		dropoffLat, err := strconv.ParseFloat(csvOrders[i][59], 32)
		errorHandler(err, "DropoffLat", "float parsing", "latitudeto")
		pickupLon, err := strconv.ParseFloat(csvOrders[i][27], 32)
		errorHandler(err, "PickupLon", "float parsing", "longitude")
		pickupLat, err := strconv.ParseFloat(csvOrders[i][26], 32)
		errorHandler(err, "PickupLat", "float parsing", "latitude")
		orderTakenTime, err := TimeParser(csvOrders[i][73])
		errorHandler(err, "OrderTakenTime", "time parsing", "appointtime")
		paymentType := "Картой"
		if csvOrders[i][64] == "f" {
			paymentType = "Наличные"
		}
		waitingTime, err := IntegerParser(csvOrders[i][16])
		errorHandler(err, "WaitingTime", "int parsing", "waiting")
		waitingTime *= 60
		dropoffDatetime, err := TimeParser(csvOrders[i][79])
		errorHandler(err, "DropoffDatetime", "time parsing", "s_time_stop_taxometr")
		//TODO спросить про csv поле stoimost
		//TODO спросить про название папки day
		order = BQOrderRaw{
			UUID:               csvOrders[i][0], //idx
			RoutesCount:        0,
			ServiceName:        csvOrders[i][40], //orderoptionid
			Features:           csvOrders[i][48], //feauteres
			CreatedDatetime:    createdDatetime,  //createtime
			Source:             "",
			OrderState:         csvOrders[i][50], //state
			CancelReason:       "",
			OrderTakenTime:     orderTakenTime, //appointtime
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
			PickupAddress:      csvOrders[i][2],     //addressfrom
			DropoffLon:         float32(dropoffLon), //longitudeto
			DropoffLat:         float32(dropoffLat), //latitudeto
			DropoffDatetime:    dropoffDatetime,     //s_time_stop_taxometr
			DropoffArea:        "",
			DropoffAddress:     csvOrders[i][33], //addresstofull
			TariffName:         "",
			TariffPrice:        0,
			RealPrice:          0,
			WaitingTime:        waitingTime, //*waiting
			WaitingPrice:       0,
			BonusPayment:       0,
			GuaranteedIncome:   0,
			ClientAllowance:    0,
			DriverUUID:         "",
			DriverCar:          "",
			DriverTarrif:       "",
			ClientPhone:        csvOrders[i][24], //aclientphone
			ClientUUID:         csvOrders[i][1],  //clientid
			PaymentType:        paymentType,      //*withcardpayment
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
		if int(createdDatetime.Month()) < 10 {
			createdMonth = "0" + createdMonth
		}
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

func TimeParser(strTime string) (time.Time, error) {
	if strTime == "" {
		return time.Time{}, nil
	}
	layout := "2006-01-02 15:04:05.999999-07"
	timeTime, err := time.Parse(layout, strTime)
	return timeTime, err
}

func IntegerParser(intStr string) (int, error) {
	if intStr == "" {
		return 0, nil
	}
	intInt, err := strconv.Atoi(intStr)
	return intInt, err
}

func errorHandler(err error, jsonField string, event string, csvField string) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"json field": jsonField, "event": event, "csv field": csvField}).Fatal(err)
	}
}
