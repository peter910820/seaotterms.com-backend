package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"seaotterms.com-backend/internal/model"
)

type BrandRecordForClient struct {
	Brand       string `json:"brand"`
	Username    string `json:"username"`
	Completed   int    `json:"completed"`
	Total       int    `json:"total"`
	Dissolution bool   `json:"dissolution"`
}

type BrandRecordForUpdate struct {
	Brand       string
	Completed   int
	Total       int
	Annotation  string
	Dissolution bool
	UpdateName  string
	UpdateTime  time.Time
}

// query all galgamebrand data
func QueryAllGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.BrandRecord

	r := db.Order("brand asc").Find(&data)
	if r.Error != nil {
		logrus.Errorf("%s\n", r.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Info("GalgameBrand全部資料查詢成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// use brand name to query single galgamebrand data
func QueryGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data model.BrandRecord
	// URL decoding
	brand, err := url.QueryUnescape(c.Params("brand"))
	if err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	r := db.Where("brand = ?", brand).First(&data)
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Errorf("%s\n", r.Error.Error())
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Info("GalgameBrand單筆資料查詢成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// insert data to galgamebrand
func InsertGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData BrandRecordForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	annotation := "待攻略"
	if clientData.Completed == clientData.Total {
		annotation = "制霸"
	}

	data := model.BrandRecord{
		Brand:       clientData.Brand,
		Completed:   clientData.Completed,
		Total:       clientData.Total,
		Annotation:  annotation,
		Dissolution: clientData.Dissolution,
		InputName:   clientData.Username,
		UpdateName:  clientData.Username,
	}
	r := db.Create(&data)
	if r.Error != nil {
		logrus.Errorf("%s\n", r.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Infof("資料 %s 創建成功", clientData.Brand)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 創建成功", clientData.Brand),
	})
}

// update single galgamebrand data
func UpdateGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData BrandRecordForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	// URL decoding
	brand, err := url.QueryUnescape(c.Params("brand"))
	if err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	annotation := "待攻略"
	if clientData.Completed == clientData.Total {
		annotation = "制霸"
	}

	// gorm:"autoUpdateTime" can not update, so manual update update_time
	r := db.Model(&model.BrandRecord{}).Where("brand = ?", brand).
		Select("brand", "completed", "total", "annotation", "dissolution", "update_name", "update_time").
		Updates(BrandRecordForUpdate{
			Brand:       clientData.Brand,
			Completed:   clientData.Completed,
			Total:       clientData.Total,
			Annotation:  annotation,
			Dissolution: clientData.Dissolution,
			UpdateName:  clientData.Username,
			UpdateTime:  time.Now(),
		})
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Errorf("%s\n", r.Error.Error())
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Infof("資料 %s 更新成功", brand)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 更新成功", brand),
	})
}
