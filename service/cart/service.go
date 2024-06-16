package cart

import (
	"fmt"

	"github.com/dclouisDan/ecommerce-api/types"
)

func getCartItemsID(items []types.CartCheckoutItem) ([]int, error) {
	itemIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity == 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductId)
		}
		itemIDs[i] = item.ProductId
	}

	return itemIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartCheckoutItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}
	// check if all products are in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	// calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce quantity of products in the db
	for _, item := range items {
		product := productMap[item.ProductId]
		product.Quantity -= item.Quantity
		if err := h.productStore.UpdateProduct(product); err != nil {
			return 0, 0, err
		}
	}
	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}
	// create order items
	for _, item := range items {
		err := h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductId].Price,
		})
		if err != nil {
			return 0, 0, err
		}
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(items []types.CartCheckoutItem, productMap map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}
	for _, item := range items {
		product, ok := productMap[item.ProductId]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductId)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is not available in the quantity requested.", item.ProductId)
		}
	}
	return nil
}

func calculateTotalPrice(items []types.CartCheckoutItem, productMap map[int]types.Product) float64 {
	var total float64
	for _, item := range items {
		total += productMap[item.ProductId].Price * float64(item.Quantity)
	}
	return total
}
