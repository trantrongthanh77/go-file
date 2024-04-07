package controller

import (
	"encoding/json"
	"fmt"
	"go-file/common"
	"go-file/model"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageDeleteRequest struct {
	Filename string `json:"filename"`
	//Token    string
}

func UploadImage(c *gin.Context) {
	uploader := c.GetString("username")
	if uploader == "" {
		uploader = "Anonymous User"
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("get form err: %s", err.Error()),
		})
		return
	}
	images := form.File["image"]
	var filenames []string
	for _, file := range images {
		id := uuid.New().String()
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s%s", id, ext)
		if err := c.SaveUploadedFile(file, filepath.Join(common.ImageUploadPath, filename)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("upload file err: %s", err.Error()),
			})
			return
		}
		imageObj := &model.Image{
			Filename: filename,
			Uploader: uploader,
			Time:     currentTime,
		}
		err = imageObj.Insert()
		if err != nil {
			common.SysError("failed to insert image to database: " + err.Error())
			continue
		}
		filenames = append(filenames, filename)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    filenames,
	})
}

func DeleteImage(c *gin.Context) {
	var deleteRequest ImageDeleteRequest
	err := json.NewDecoder(c.Request.Body).Decode(&deleteRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid parameter",
		})
		return
	}

	imageObj := &model.Image{
		Filename: deleteRequest.Filename,
	}
	rowsAffected := model.DB.Where("filename = ?", deleteRequest.Filename).First(&imageObj).RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "File does not exist!",
		})
		return
	}
	err = imageObj.Delete()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}
