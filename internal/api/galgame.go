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

type GameRecordForClient struct {
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	ReleaseDate time.Time `json:"releaseDate"`
	AllAges     bool      `json:"allAges"`
	EndDate     time.Time `json:"endDate"`
	Username    string    `json:"username"`
}

// get all the galgame data for specify brand
func QueryGalgame(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.GameRecord
	// URL decoding
	brand, err := url.QueryUnescape(c.Params("brand"))
	if err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	r := db.Where("brand = ?", brand).Find(&data)
	if r.Error != nil {
		logrus.Fatal(r.Error.Error())
	}
	// if data not exist, retrun a empty struct
	logrus.Info("Galgame多筆資料查詢成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// insert data to galgame
func InsertGalgame(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData GameRecordForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	data := model.GameRecord{
		Name:        clientData.Name,
		Brand:       clientData.Brand,
		ReleaseDate: clientData.ReleaseDate,
		AllAges:     clientData.AllAges,
		EndDate:     clientData.EndDate,
		InputName:   clientData.Username,
	}
	r := db.Create(&data)
	if r.Error != nil {
		logrus.Errorf("%s\n", r.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Infof("資料 %s 創建成功", clientData.Name)

	/* --------------------------------- */
	/* --------------------------------- */

	// update brand info
	var brandData model.BrandRecord

	r = db.Where("brand = ?", clientData.Brand).First(&brandData)
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

	annotation := "待攻略"
	if brandData.Completed+1 == brandData.Total {
		annotation = "制霸"
	}

	// gorm:"autoUpdateTime" can not update, so manual update update_time
	r = db.Model(&model.BrandRecord{}).Where("brand = ?", clientData.Brand).
		Select("completed", "annotation", "update_name", "update_time").
		Updates(BrandRecordForUpdate{
			Completed:  brandData.Completed + 1,
			Annotation: annotation,
			UpdateName: clientData.Username,
			UpdateTime: time.Now(),
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 創建成功", clientData.Name),
	})
}
