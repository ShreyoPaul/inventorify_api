package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attribute struct {
	Key   string
	Value string
}

type Product struct {
	ID         primitive.ObjectID
	Name       string
	Attributes []string
	Qty        int64
	Amount     uint64
	CGST       float64
	SGST       float64
}

type Inv struct {
	ID         primitive.ObjectID
	Name       string
	Attributes []string
	Products   []Product
}

type Users struct {
	Name     string
	Email    string
	Password string
	Inv      []Inv
}

type Claims struct {
	Email string
	jwt.StandardClaims
}
