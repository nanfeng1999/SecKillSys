package model

type Coupon struct {
	Id          int64    `gorm:"primary_key;auto_increment"` // 优惠券ID
	Username    string   `gorm:"type:varchar(20); not null"` // 用户名
	CouponName  string   `gorm:"type:varchar(60); not null"` // 优惠券名称
	Amount      int64	  // 金额
	Left        int64	  // 剩余数量
	Stock       int64	  // 库存
	Description string  `gorm:"type:varchar(60)"`
}

type ReqCoupon struct {
	Name			string
	Amount 			int64
	Description     string
	Stock           int64
}

type ResCoupon struct {
	Name            string  `json:"name"`
	Stock           int64   `json:"stock"`
	Description     string  `json:"description"`
}

// SellerResCoupon 商家查询优惠券时，返回的数据结构
type SellerResCoupon struct {
	ResCoupon
	Amount int64  `json:"amount"`
	Left   int64  `json:"left"`
}

// CustomerResCoupon 顾客查询优惠券时，返回的数据结构
type CustomerResCoupon struct {
	ResCoupon
}

func ParseSellerResCoupons(coupons []Coupon) []SellerResCoupon {
	var sellerCoupons []SellerResCoupon
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			SellerResCoupon{ResCoupon{coupon.CouponName, coupon.Stock, coupon.Description},
				coupon.Amount, coupon.Left})
	}
	return sellerCoupons
}

func ParseCustomerResCoupons(coupons []Coupon) []CustomerResCoupon {
	var sellerCoupons []CustomerResCoupon
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			CustomerResCoupon{ResCoupon{coupon.CouponName, coupon.Stock, coupon.Description}})
	}
	return sellerCoupons
}