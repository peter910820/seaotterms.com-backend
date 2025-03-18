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

type GameRecordForUpdate struct {
	ReleaseDate time.Time
	EndDate     time.Time
	UpdateName  string
	UpdateTime  time.Time
}

// use game name to query single galgame data
func QueryGalgame(c *fiber.Ctx, db *gorm.DB) error {
	var data model.GameRecord
	// URL decoding
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	r := db.Where("name = ?", name).First(&data)
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Error(r.Error)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Info("Galgame單筆資料查詢成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// get all the galgame data for specify brand
func QueryGalgameByBrand(c *fiber.Ctx, db *gorm.DB) error {
	var data []model.GameRecord
	// URL decoding
	brand, err := url.QueryUnescape(c.Params("brand"))
	if err != nil {
		logrus.Error(err)
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

// update single galgame data (develop)
func UpdateGalgameDevelop(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData GameRecordForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	// URL decoding
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// gorm:"autoUpdateTime" can not update, so manual update update_time
	r := db.Model(&model.GameRecord{}).Where("name = ?", name).
		Select("release_date", "end_date", "update_name", "update_time").
		Updates(GameRecordForUpdate{
			ReleaseDate: clientData.ReleaseDate,
			EndDate:     clientData.EndDate,
			UpdateName:  clientData.Username,
			UpdateTime:  time.Now(),
		})
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Error(r.Error)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Infof("資料 %s 更新成功", name)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 更新成功", name),
	})
}

// insert data to galgame
func InsertGalgame(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData GameRecordForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
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
		UpdateName:  clientData.Username,
	}
	r := db.Create(&data)
	if r.Error != nil {
		logrus.Error(r.Error)
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
			logrus.Error(r.Error)
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
			logrus.Error(r.Error)
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
