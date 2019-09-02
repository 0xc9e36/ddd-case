package interfaces

import (
	"domain"
	"fmt"
	"usecases"
)

type DbHandler interface {
	Execute(statement string)
	Query(statement string) Row
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type DbUserRepo DbRepo
type DbCustomerRepo DbRepo
type DbOrderRepo DbRepo
type DbItemRepo DbRepo

func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbUserRepo"]
	return dbUserRepo
}

func (repo *DbUserRepo) Store(user usecases.User) error {
	isAdmin := "no"
	if user.IsAdmin {
		isAdmin = "yes"
	}
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO users (id, customer_id,  )
                                        VALUES ('%d', '%d', '%v')`,
		user.ID, user.Customer.ID, isAdmin))
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	customerRepo.Store(user.Customer)
	return  nil
}

func (repo *DbUserRepo) FindByID(id int) usecases.User {
	sql := fmt.Sprintf(`SELECT is_admin, customer_id FROM users WHERE id = '%d' LIMIT 1`, id)
	row := repo.dbHandler.Query(sql)
	var isAdmin string
	var customerId int
	row.Next()
	row.Scan(&isAdmin, &customerId)
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	u := usecases.User{ID: id, Customer: customerRepo.FindByID(customerId)}
	u.IsAdmin = false
	if isAdmin == "yes" {
		u.IsAdmin = true
	}
	return u
}

func NewDbCustomerRepo(dbHandlers map[string]DbHandler) *DbCustomerRepo {
	dbCustomerRepo := new(DbCustomerRepo)
	dbCustomerRepo.dbHandlers = dbHandlers
	dbCustomerRepo.dbHandler = dbHandlers["DbCustomerRepo"]
	return dbCustomerRepo
}

func (repo *DbCustomerRepo) Store(customer domain.Customer) {
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO customers (id, name)
                                        VALUES ('%d', '%v')`,
		customer.ID, customer.Name))
}

func (repo *DbCustomerRepo) FindByID(id int) domain.Customer {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT name FROM customers
                                             WHERE id = '%d' LIMIT 1`,
		id))
	var name string
	row.Next()
	row.Scan(&name)
	return domain.Customer{ID: id, Name: name}
}

func NewDbOrderRepo(dbHandlers map[string]DbHandler) *DbOrderRepo {
	dbOrderRepo := new(DbOrderRepo)
	dbOrderRepo.dbHandlers = dbHandlers
	dbOrderRepo.dbHandler = dbHandlers["DbOrderRepo"]
	return dbOrderRepo
}

func (repo *DbOrderRepo) Store(order domain.Order) error {
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO orders (id, customer_id)
                                        VALUES ('%d', '%v')`,
		order.ID, order.Customer.ID))
	for _, item := range order.Items {
		repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO items2orders (item_id, order_id)
                                            VALUES ('%d', '%d')`,
			item.ID, order.ID))
	}
	return nil
}

func (repo *DbOrderRepo) FindByID(id int) domain.Order {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT customer_id FROM orders
                                             WHERE id = '%d' LIMIT 1`,
		id))
	var customerId int
	row.Next()
	row.Scan(&customerId)
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	order := domain.Order{ID: id, Customer: customerRepo.FindByID(customerId)}
	var itemId int
	itemRepo := NewDbItemRepo(repo.dbHandlers)
	row = repo.dbHandler.Query(fmt.Sprintf(`SELECT item_id FROM items2orders
                                            WHERE order_id = '%d'`,
		order.ID))
	for row.Next() {
		row.Scan(&itemId)
		order.Add(itemRepo.FindByID(itemId))
	}
	return order
}

func NewDbItemRepo(dbHandlers map[string]DbHandler) *DbItemRepo {
	dbItemRepo := new(DbItemRepo)
	dbItemRepo.dbHandlers = dbHandlers
	dbItemRepo.dbHandler = dbHandlers["DbItemRepo"]
	return dbItemRepo
}

func (repo *DbItemRepo) Store(item domain.Item) error {
	available := "no"
	if item.Available {
		available = "yes"
	}
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO items (id, name, value, available)
                                        VALUES ('%d', '%v', '%f', '%v')`,
		item.ID, item.Name, item.Value, available))
	return nil
}

func (repo *DbItemRepo) FindByID(id int) domain.Item {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT name, value, available
                                             FROM items WHERE id = '%d' LIMIT 1`,
		id))
	var name string
	var value float64
	var available string
	row.Next()
	row.Scan(&name, &value, &available)
	item := domain.Item{ID: id, Name: name, Value: value}
	item.Available = false
	if available == "yes" {
		item.Available = true
	}
	return item
}
