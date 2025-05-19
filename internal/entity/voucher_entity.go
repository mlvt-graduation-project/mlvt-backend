package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VoucherCode struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code        string             `json:"code" bson:"code"`
	Token       uint64             `json:"token" bson:"token"`
	MaxUsage    uint               `json:"max_usage" bson:"max_usage"`
	UsedCount   uint               `json:"used_count" bson:"used_count"`
	ExpiredTime time.Time          `json:"expired_time" bson:"expired_time"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type GetAllVoucherRequest struct {
	Status         string `form:"status"`
	SortBy         string `form:"sortBy"`
	SearchKey      string `form:"searchKey"`
	Sort           string `form:"sort"`
	SearchCriteria string `form:"searchCriteria"`
	Offset         int    `form:"from"`
	Limit          int    `form:"to"`
}

type GetAllVoucherResponse struct {
	Vouchers []VoucherCode `json:"vouchers"`
	TotalCount int64 `json:"total_count"`
}