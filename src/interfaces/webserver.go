package interfaces

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	usecases2 "usecases"
)

type OrderInteractor interface {
	Items(userID, orderID int) ([]usecases2.Item, error)
	Add(userID, orderID, itemID int) error
}

type WebServiceHandler struct {
	OrderInteractor OrderInteractor
}

func (w WebServiceHandler) ShowOrder(res http.ResponseWriter, req *http.Request) {
	userID, _ := strconv.Atoi(req.FormValue("userID"))
	orderID, _ := strconv.Atoi(req.FormValue("orderID"))


	items, _ := w.OrderInteractor.Items(userID, orderID)

	for _, item := range items {
		io.WriteString(res, fmt.Sprintf("item id: %d\n", item.ID))
		io.WriteString(res, fmt.Sprintf("item name: %v\n", item.Name))
		io.WriteString(res, fmt.Sprintf("item value: %f\n", item.Value))
	}
}
