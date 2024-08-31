package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product_Details struct {
	ID         primitive.ObjectID
	Name       string
	Attributes []string
	Qty        int64
	Amount     uint64
	CGST       float64
	SGST       float64
}

type Bill struct {
	_id       primitive.ObjectID `bson:"_id"`
	U_id      primitive.ObjectID
	U_email   string
	Shop      Shop
	Party     Party
	Products  []Product_Details
	CreatedAt time.Time `bson:"createdAt"`
}

type Shop struct {
	S_Name    string
	S_Email   string
	S_phone   string
	S_address string
}

type Party struct {
	Name    string
	Email   string
	Phone   string
	Address string
	Pan     string
	Gstin   string
}
