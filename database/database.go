package database

import (
	"fmt"
	"os"
	"database/sql"
	 _ "github.com/go-sql-driver/mysql"
	 "github.com/jeremymwells/easyConfig"
)

type DatabaseUser struct {
	Name string
	Password string
}

type DatabaseDefinition struct {
	Name string
	Protocol string
	Address string
}

type Configuration struct {
	ReadUser DatabaseUser
	WriteUser DatabaseUser
	Dev DatabaseDefinition
	Prod DatabaseDefinition
}

type Database struct{
	WriteDb *sql.DB
	ReadDb *sql.DB
}

type DbFile struct {
	Id int64
	Sha1 [20]byte
	Md5 [16]byte
	Size int64
	ContentType string
}

type DbPastebinFile struct {
	Id int64
	File DbFile
	Url string
}

var (
	dbConfig = easyConfig.New(&Configuration{}, "../database/config.json").(*Configuration)
	currentDbDefinition = dbConfig.Dev //TODO: read command line flag
	readConnectionString = buildConnectionString(dbConfig.ReadUser)
	writeConnectionString = buildConnectionString(dbConfig.WriteUser)
)

func buildConnectionString(user DatabaseUser) string {
	
	return fmt.Sprintf("%s:%s@%s(%s)/%s?timeout=5s", 
		user.Name, user.Password, 
		currentDbDefinition.Protocol, 
		currentDbDefinition.Address, 
		currentDbDefinition.Name)
}

func Get() Database {
	readDb, err := sql.Open("mysql", readConnectionString)
	pingErr := readDb.Ping()
	if err != nil || pingErr != nil {
		fmt.Println("error opening read database: %v %v", err, pingErr)
		os.Exit(1)
	}
	
	writeDb, err := sql.Open("mysql", writeConnectionString)
	pingErr = writeDb.Ping()
	if err != nil || pingErr != nil {
		fmt.Println("error opening write database: %v %v", err, pingErr)
		os.Exit(1)
	}
	
	return Database{writeDb, readDb}
}

func (this *Database) WriteFileRecord(file *DbFile, close bool) DbFile {
	
	fmt.Println("writing a file: ", *file)
    stmt, err := this.WriteDb.Prepare("INSERT INTO file SET sha1=?, md5=?, size=?, content_type=?")
    if err != nil{
		fmt.Println("error preparing file save sql statement: ", err);		
	}
	
	res, err := stmt.Exec(fmt.Sprintf("%x", file.Sha1), fmt.Sprintf("%x", file.Md5), file.Size, file.ContentType)
    if err != nil{
		fmt.Println("error executing file save sql statement: ", err);		
	}
	
	fmt.Println("file written to db: ", res)
	id, err := res.LastInsertId()
    if err != nil{
		fmt.Println("error getting last inserted id from file save sql statement: ", err);		
	}
	fmt.Println("new id: ", id)
	
	file.Id = id
	
	if close{
		this.WriteDb.Close();
	}
	
	return *file;
}

func (this *Database) WritePastebinFile(file *DbPastebinFile, close bool) DbPastebinFile {
	
	this.WriteFileRecord(&file.File, close)
	
	fmt.Println("writing a pastebin file: ", *file)
    stmt, err := this.WriteDb.Prepare("INSERT INTO pastebin SET file_id=?, url=?")
    if err != nil{
		fmt.Println("error preparing file save sql statement: ", err);		
	}
	
	res, err := stmt.Exec(file.File.Id, file.Url)
    if err != nil{
		fmt.Println("error executing file save sql statement: ", err);		
	}
	
	fmt.Println("pastebin file written to db: ", res)
	id, err := res.LastInsertId()
    if err != nil{
		fmt.Println("error getting last inserted id from file save sql statement: ", err);		
	}
	
	fmt.Println("new id: ", id)
	file.Id = id
	
	if close{
		this.WriteDb.Close();
	}
	
	return *file;
}



// -- start a new transaction
// start transaction;
 
// -- get latest order number
// select @orderNumber := max(orderNUmber) 
// from orders;
// -- set new order number
// set @orderNumber = @orderNumber  + 1;
 
// -- insert a new order for customer 145
// insert into orders(orderNumber,
//                    orderDate,
//                    requiredDate,
//                    shippedDate,
//                    status,
//                    customerNumber)
// values(@orderNumber,
//        now(),
//        date_add(now(), INTERVAL 5 DAY),
//        date_add(now(), INTERVAL 2 DAY),
//        'In Process',
//         145);
// -- insert 2 order line items
// insert into orderdetails(orderNumber,
//                          productCode,
//                          quantityOrdered,
//                          priceEach,
//                          orderLineNumber)
// values(@orderNumber,'S18_1749', 30, '136', 1),
//       (@orderNumber,'S18_2248', 50, '55.09', 2); 
// -- commit changes    
// commit;       
 
// -- get the new inserted order
// select * from orders a 
// inner join orderdetails b on a.ordernumber = b.ordernumber
// where a.ordernumber = @ordernumber;

