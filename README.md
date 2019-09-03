# 整洁架构


## 文章
英文: [Applying The Clean Architecture to Go applications](https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/)    
中文：[在 GO 应用中使用简明架构](https://mikespook.com/2012/09/翻译在-go-应用中使用简明架构1/#more-1440)   


## run
1. 导入 sql 数据
2. ```go run main.go```
3. 访问:  ```http://localhost:8080/orders?userID=40&orderID=60```

## 其他
文章中使用的是 sqlite3，我这里替换成了 mysql
