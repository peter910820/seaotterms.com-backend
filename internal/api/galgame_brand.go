package api

import (
	"fmt"
	"net/url"

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

// query all galgamebrand data
func QueryALlGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.BrandRecord

	r := db.Order("brand desc").Find(&data)
	if r.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	r := db.Where("brand = ?", brand).First(&data)
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Debugf("%v", data)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// insert data to galgamebrand
func InsertGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData BrandRecordForClient
	if err := c.BodyParser(&clientData); err != nil {
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
	result := db.Create(&data)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 創建成功", clientData.Brand),
	})
}
