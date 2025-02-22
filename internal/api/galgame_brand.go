package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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

func QueryALlGalgameBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.BrandRecord

	clintData := db.Order("brand desc").Find(&data)
	if clintData.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": clintData.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

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
