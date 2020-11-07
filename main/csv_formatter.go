package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
}

func main() {
	orderoptionFile, err := os.Open("public_orderoption.csv")
	if err != nil {
		logrus.WithField("event", "opening orderoptionid description file").Fatal(err)
	}
	csvReader := csv.NewReader(orderoptionFile)
	orderoption, err := csvReader.ReadAll()
	if err != nil {
		logrus.WithField("event", "reading orderoptionid description file").Fatal(err)
	}

	fileName := flag.String("n", "main.csv", "file name")
	fileDir := flag.String("d", "./", "file directory")
	ordersDirPath := flag.String("sd", "orders", "saving directory")
	flag.Parse()
	var filePath string
	if (*fileDir)[len(*fileDir)-1] == '/' {
		filePath = fmt.Sprintf("%s%s", *fileDir, *fileName)
	} else {
		filePath = fmt.Sprintf("%s/%s", *fileDir, *fileName)
	}
	if (*ordersDirPath)[len(*fileDir)-1] == '/' {
		*ordersDirPath = removeCharByIndex(*ordersDirPath, len(*fileDir)-1)
	}
	file, err := os.Open(filePath)
	if err != nil {
		logrus.WithField("event", "opening csv file").Fatal(err)
	}
	var csvOrders [][]string
	var csvString string
	var order BQOrderRaw
	csvString, err = readString(file)
	if err != nil {
		logrus.WithField("event", "reading string from csv file").Fatal(err)
	}
	for {
		csvString, err = readString(file)
		if err != nil {
			if err == io.EOF {
				logrus.WithField("event", "file read successfully").Info(err)
				break
			}
			logrus.WithField("event", "reading string from csv file").Fatal(err)
		}
		strReader := strings.NewReader(csvString)
		csvReader := csv.NewReader(strReader)
		csvOrders, err = csvReader.ReadAll()
		if err != nil {
			logrus.WithField("event", "reading from csv reader").Fatal(err)
		}
		createdDatetime, err := timeParser(csvOrders[0][5])
		errorHandler(err, "CreatedDatetime", "time parsing", "createtime")
		dropoffLon, err := floatParser(csvOrders[0][60])
		errorHandler(err, "DropoffLon", "float parsing", "longitudeto")
		dropoffLat, err := floatParser(csvOrders[0][59])
		errorHandler(err, "DropoffLat", "float parsing", "latitudeto")
		pickupLon, err := floatParser(csvOrders[0][27])
		errorHandler(err, "PickupLon", "float parsing", "longitude")
		pickupLat, err := floatParser(csvOrders[0][26])
		errorHandler(err, "PickupLat", "float parsing", "latitude")
		orderTakenTime, err := timeParser(csvOrders[0][73])
		errorHandler(err, "OrderTakenTime", "time parsing", "appointtime")
		paymentType := "Картой"
		if csvOrders[0][64] == "f" {
			paymentType = "Наличные"
		}
		waitingTime, err := intParser(csvOrders[0][16])
		errorHandler(err, "WaitingTime", "int parsing", "waiting")
		waitingTime *= 60
		dropoffDatetime, err := timeParser(csvOrders[0][79])
		errorHandler(err, "DropoffDatetime", "time parsing", "s_time_stop_taxometr")
		tariffPrice, err := intParser(csvOrders[0][36])
		errorHandler(err, "TariffPrice", "int parsing", "stoimost_tarif")
		realPrice, err := intParser(csvOrders[0][11])
		errorHandler(err, "RealPrice", "int parsing", "stoimost")
		serviceName := serviceNameHandler(csvOrders[0][40], orderoption)
		//TODO спросить про json псевдонимы(не соответствуют названиям полей из Google Диск)
		//TODO спросить про описание ArrivalTimePromise(тип time, а в описании секунды)
		//TODO спросить, нужно ли преобразовать значения из state
		//TODO спросить, что значат OwnerUUID, UserStartName, UserStartUUID(описание: тип оплаты)
		//TODO спросить про пустые поля PickupArea и DropoffArea на Google Диске
		//TODO спросить, что значат BonusPayment, GuaranteedIncome, ClientAllowance, InsertDateTime
		//TODO спросить, нужно ли выводить логи сохранения файлов
		//TODO спросить, можно ли брать описание различных id из файлов
		order = BQOrderRaw{
			UUID:               csvOrders[0][0], //idx
			RoutesCount:        0,
			ServiceName:        serviceName,      //orderoptionid
			Features:           csvOrders[0][48], //feauteres
			CreatedDatetime:    createdDatetime,  //createtime
			Source:             "crm",
			OrderState:         csvOrders[0][50], //state
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
			PickupLon:          pickupLon, //longitude
			PickupLat:          pickupLat, //latitude
			PickupDatetime:     time.Time{},
			PickupArea:         "",
			PickupAddress:      csvOrders[0][2], //addressfrom
			DropoffLon:         dropoffLon,      //longitudeto
			DropoffLat:         dropoffLat,      //latitudeto
			DropoffDatetime:    dropoffDatetime, //s_time_stop_taxometr
			DropoffArea:        "",
			DropoffAddress:     csvOrders[0][33], //addresstofull
			TariffName:         "",
			TariffPrice:        tariffPrice, //stoimost_tarif
			RealPrice:          realPrice,   //stoimost
			WaitingTime:        waitingTime, //waiting
			WaitingPrice:       0,
			BonusPayment:       0,
			GuaranteedIncome:   0,
			ClientAllowance:    0,
			DriverUUID:         csvOrders[0][82], //adriverid
			DriverCar:          "",
			DriverTarrif:       "",
			ClientPhone:        csvOrders[0][24], //aclientphone
			ClientUUID:         csvOrders[0][1],  //clientid
			PaymentType:        paymentType,      //withcardpayment
			StoreUUID:          "",
			ProductsSum:        0,
			ProductsCount:      0,
			ProductsData:       "",
			InsertDateTime:     time.Time{},
		}
		createdYear := strconv.Itoa(createdDatetime.Year())
		createdMonth := strconv.Itoa(int(createdDatetime.Month()))
		if int(createdDatetime.Month()) < 10 {
			createdMonth = "0" + createdMonth
		}
		createdDay := strconv.Itoa(createdDatetime.Day())
		if createdDatetime.Day() < 10 {
			createdDay = "0" + createdDay
		}
		savingPath := fmt.Sprintf("%s/%s/%s/%s/", *ordersDirPath, createdYear, createdMonth, createdDay)
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
		logrus.WithFields(logrus.Fields{"path": savingPath, "name": fileName}).Info(errors.New("json file saved"))
	}
	err = file.Close()
	if err != nil {
		logrus.WithField("event", "csv file closing").Fatal(err)
	}
}

func readString(file *os.File) (string, error) {
	reader := bufio.NewReaderSize(file, 4096)
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	_, err = file.Seek(int64(-reader.Buffered()), 1)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func timeParser(strTime string) (time.Time, error) {
	if strTime == "" {
		return time.Time{}, nil
	}
	layout := "2006-01-02 15:04:05.999999-07"
	timeTime, err := time.Parse(layout, strTime)
	return timeTime, err
}

func intParser(strInt string) (int, error) {
	if strInt == "" {
		return 0, nil
	}
	if strings.Contains(strInt, ".") {
		floatInt, err := floatParser(strInt)
		return int(floatInt), err
	}
	intInt, err := strconv.Atoi(strInt)
	return intInt, err
}

func floatParser(strFloat string) (float32, error) {
	if strFloat == "" {
		return 0, nil
	}
	floatFloat, err := strconv.ParseFloat(strFloat, 32)
	return float32(floatFloat), err
}

func errorHandler(err error, jsonField string, event string, csvField string) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"json field": jsonField, "event": event, "csv field": csvField}).Fatal(err)
	}
}

func removeCharByIndex(s string, i int) string {
	c := []rune(s)
	s = string(append(c[0:i], c[i+1:]...))
	return s
}

func serviceNameHandler(orderoptionid string, orderoption [][]string) string {
	for i := 1; i < len(orderoption); i++ {
		if orderoptionid == orderoption[i][0] {
			return orderoption[i][1]
		}
	}
	return ""
}
