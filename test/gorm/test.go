package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"my_item_cf_go/plugin/item_cf_big_data/cf_lib"
	"my_item_cf_go/plugin/myorm"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type Rating struct {
	// gorm.Model
	UserId    int
	MovieId   int
	Rating    float64
	Timestamp int
}

func main() {

	initConfig()

	// db, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN:                       "root:123456@tcp(127.0.0.1:3304)/test?charset=utf8&parseTime=True&loc=Local", // DSN data source name
	// 	DefaultStringSize:         256,                                                                          // string 类型字段的默认长度
	// 	DisableDatetimePrecision:  true,                                                                         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// 	DontSupportRenameIndex:    true,                                                                         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// 	DontSupportRenameColumn:   true,                                                                         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// 	SkipInitializeWithVersion: false,                                                                        // 根据当前 MySQL 版本自动配置
	// }), &gorm.Config{})

	db, err := myorm.Connect()

	if err != nil {
		panic("failed to connect database")
	}
	cf := cf_lib.GetItemCF()

	// 查询评分表总数
	var total int64
	db.Table("rating").Count(&total)
	pageSize := 1000
	// offset := (page - 1) * pageSize
	var rating []Rating

	for page := 1; page*pageSize < int(total); page++ {
		db.Table("rating").Unscoped().Scopes(Paginate(page, pageSize)).Find(&rating)
		// for i := 0; i < len(rating); i++ {
		// 	fmt.Println(rating[i])
		// }

		for _, value := range rating {
			// fmt.Println(value)
			// score, _ := strconv.ParseFloat(value.Rating, 64)
			score := value.Rating
			uid := strconv.Itoa(value.UserId)
			movieId := strconv.Itoa(value.UserId)

			if rand.Float64() <= cf.TrainSetPecent {
				cf.TestNum++
				if _, ok := cf.TrainSetRec[uid]; !ok {
					cf.TrainSetRec[uid] = make(map[string]float64)
				}
				cf.TrainSetRec[uid][movieId] = float64(score)
			} else {
				cf.TrainNum++
				if _, ok := cf.TestSet[uid]; !ok {
					cf.TestSet[uid] = make(map[string]float64)
				}
				cf.TestSet[uid][movieId] = float64(score)
			}
		}
	}
	fmt.Printf("数据分组完成，训练集包含数据%d条,测试集包含数据%d条 \n", cf.TrainNum, cf.TestNum)
}

func main1() {
	cf := cf_lib.GetItemCF()
	DataPath := "../../plugin/item_cf/ml-1m/ratings.csv"
	// fs, err := os.Open(cf.DataPath)
	fs, err := os.Open(DataPath)
	if err != nil {
		log.Fatalf("无法打开数据文件: %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	r.Read()
	rand.Seed(time.Now().UnixNano())
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("文件读取错误 ： %+v", err)
		}
		if err == io.EOF {
			break
		}
		score, _ := strconv.ParseFloat(row[2], 64)
		if rand.Float64() <= cf.TrainSetPecent {
			cf.TestNum++
			if _, ok := cf.TrainSetRec[row[0]]; !ok {
				cf.TrainSetRec[row[0]] = make(map[string]float64)
			}
			cf.TrainSetRec[row[0]][row[1]] = float64(score)
		} else {
			cf.TrainNum++
			if _, ok := cf.TestSet[row[0]]; !ok {
				cf.TestSet[row[0]] = make(map[string]float64)
			}
			cf.TestSet[row[0]][row[1]] = float64(score)
		}
	}
	fmt.Printf("数据分组完成，训练集包含数据%d条,测试集包含数据%d条 \n", cf.TrainNum, cf.TestNum)
}

func main0() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123456@tcp(127.0.0.1:3304)/item_cf?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                             // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                            // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                            // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                            // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                           // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var total int64
	db.Table("rating").Count(&total)
	pageSize := 1000
	// offset := (page - 1) * pageSize
	var rating []Rating
	for page := 1; page*pageSize < int(total); page++ {
		db.Table("rating").Unscoped().Scopes(Paginate(page, pageSize)).Find(&rating)
		fmt.Println(page)
		// for i := 0; i < len(rating); i++ {
		// 	fmt.Println(rating[i])
		// }
		for _, value := range rating {
			fmt.Println(value)
		}
	}

	// // 迁移 schema
	// db.AutoMigrate(&Product{})

	// // Create
	// db.Create(&Product{Code: "D42", Price: 100})

	// // Read
	// var product Product
	// db.First(&product, 1)                 // 根据整型主键查找
	// db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// // Update - 将 product 的 price 更新为 200
	// db.Model(&product).Update("Price", 200)
	// // Update - 更新多个字段
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// // Delete - 删除 product
	// db.Delete(&product, 1)
}

// 分页封装
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read configuration file: " + err.Error())
	}
}
