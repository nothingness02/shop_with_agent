package Order

import (
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusShipped   OrderStatus = "shipped"   // 已发货
	OrderStatusDelivered OrderStatus = "delivered" // 已送达
	OrderStatusCompleted OrderStatus = "completed" // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

type Order struct {
	gorm.Model
	OrderID        string      `gorm:"uniqueIndex;size:64"` // 订单号
	UserID         uint        `gorm:"index"`               // 用户ID
	TotalAmount    float64     `gorm:"type:decimal(10,2)"`  // 订单总金额
	DiscountAmount float64     `gorm:"type:decimal(10,2)"`  // 优惠金额
	ShippingFee    float64     `gorm:"type:decimal(10,2)"`  // 运费
	ActualAmount   float64     `gorm:"type:decimal(10,2)"`  // 实付金额
	Status         OrderStatus `gorm:"size:32;index"`       // 订单状态
	//PaymentMethod    PaymentMethod `gorm:"size:32"`            // 支付方式
	//PaymentTime      *time.Time   // 支付时间
	//ShippingTime     *time.Time   // 发货时间
	//CompletionTime   *time.Time   // 完成时间
	//CancellationTime *time.Time   // 取消时间
	// 收货信息
	ShippingName    string `gorm:"size:100"` // 收货人姓名
	ShippingPhone   string `gorm:"size:20"`  // 收货人电话
	ShippingAddress string `gorm:"size:255"` // 收货地址
	ShippingZipCode string `gorm:"size:10"`  // 邮政编码
	Tsv             string `gorm:"type:tsvector;index:,type:gin;->"`
	//// 关联关系
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	//PaymentRecords   []PaymentRecord `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"index"`              // 订单ID
	ProductID   uint    `gorm:"index"`              // 商品ID
	ProductName string  `gorm:"size:100"`           // 商品名称（快照）
	ProductImg  string  `gorm:"size:500"`           // 商品图片（快照）
	Price       float64 `gorm:"type:decimal(10,2)"` // 单价
	Quantity    int     // 数量
	Subtotal    float64 `gorm:"type:decimal(10,2)"` // 小计
}
