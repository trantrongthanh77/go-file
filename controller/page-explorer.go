package controller

import (
	"context"
	"encoding/json"
	"go-file/common"
	"go-file/externalinterface/storage"
	"go-file/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PageExplorerController struct {
	GetExplorerPageOrFile (*gin.Context)
}

type PageExplorerControllerImpl struct {
	cloudinary storage.Cloudinary
}

func NewPageExplorer(cloud storage.Cloudinary) *PageExplorerControllerImpl {
	return &PageExplorerControllerImpl{
		cloudinary: cloud,
	}
}

func (c *PageExplorerControllerImpl) GetExplorerPageOrFile(ctx *gin.Context) {
	path := ctx.DefaultQuery("path", "/")
	path, _ = url.PathUnescape(path)

	folders := c.cloudinary.GetSubFolder(path)
	assets := c.cloudinary.GetAssets(path)
	files := make([]model.LocalFile, 0)
	for _, folder := range folders.Folders {
		files = append(files, model.LocalFile{
			Name:         folder.Name,
			Link:         "explorer?path=" + url.QueryEscape(folder.Path),
			Size:         "",
			IsFolder:     true,
			ModifiedTime: "",
		})
	}
	for _, asset := range assets.Assets {
		files = append(files, model.LocalFile{
			Name:         asset.Filename + "." + asset.Format,
			Link:         asset.SecureURL,
			Size:         common.Bytes2Size(int64(asset.Bytes)),
			IsFolder:     false,
			ModifiedTime: asset.CreatedAt.String()[:19],
		})

	}

	// fullPath := filepath.Join(common.ExplorerRootPath, path)
	// if !strings.HasPrefix(fullPath, common.ExplorerRootPath) {
	// 	// We may being attacked!
	// 	ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
	// 		"message":  "Only subdirectories of the specified folder can be accessed",
	// 		"option":   common.OptionMap,
	// 		"username": ctx.GetString("username"),
	// 	})
	// 	return
	// }
	// root, err := os.Stat(fullPath)
	// if err != nil {
	// 	ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
	// 		"message":  "An error occurred while processing the path, please confirm that the path is correct",
	// 		"option":   common.OptionMap,
	// 		"username": ctx.GetString("username"),
	// 	})
	// 	return
	// }
	// if root.IsDir() {
	// _, readmeFileLink, err := getData(path, fullPath)
	// if err != nil {
	// 	ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
	// 		"message":  err.Error(),
	// 		"option":   common.OptionMap,
	// 		"username": ctx.GetString("username"),
	// 	})
	// 	return
	// }
	var pathLinks []string
	var paths []string
	if path != "/" {
		paths = strings.Split(path, "/")
		pathLinks = make([]string, len(paths))
		for i := 0; i < len(paths); i++ {
			pathLinks[i] = "explorer?path=" + strings.Join(paths[:i+1], "/")
		}
	}

	ctx.HTML(http.StatusOK, "explorer.html", gin.H{
		"message":        "",
		"option":         common.OptionMap,
		"username":       ctx.GetString("username"),
		"files":          files,
		"readmeFileLink": "",
		"pathLinks":      pathLinks,
		"paths":          paths,
	})
	// } else {
	// 	ctx.File(filepath.Join(common.ExplorerRootPath, path))
	// }
}

func getDataFromFS(path string, fullPath string) (localFilesPtr *[]model.LocalFile, readmeFileLink string, err error) {
	var localFiles []model.LocalFile
	var tempFiles []model.LocalFile
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return
	}
	if path != "/" {
		parts := strings.Split(path, "/")
		// Add the special item: ".." which means parent dir
		if len(parts) > 0 {
			parts = parts[:len(parts)-1]
		}
		parentPath := strings.Join(parts, "/")
		parentFile := model.LocalFile{
			Name:         "..",
			Link:         "explorer?path=" + url.QueryEscape(parentPath),
			Size:         "",
			IsFolder:     true,
			ModifiedTime: "",
		}
		localFiles = append(localFiles, parentFile)
		path = strings.Trim(path, "/") + "/"
	} else {
		path = ""
	}
	for _, f := range files {
		link := "explorer?path=" + url.QueryEscape(path+f.Name())
		file := model.LocalFile{
			Name:         f.Name(),
			Link:         link,
			Size:         common.Bytes2Size(f.Size()),
			IsFolder:     f.Mode().IsDir(),
			ModifiedTime: f.ModTime().String()[:19],
		}
		if file.IsFolder {
			localFiles = append(localFiles, file)
		} else {
			tempFiles = append(tempFiles, file)
		}
		if f.Name() == "README.md" {
			readmeFileLink = link
		}
	}
	localFiles = append(localFiles, tempFiles...)
	localFilesPtr = &localFiles
	return
}

func getData(path string, fullPath string) (localFilesPtr *[]model.LocalFile, readmeFileLink string, err error) {
	if !common.ExplorerCacheEnabled {
		return getDataFromFS(path, fullPath)
	} else {
		ctx := context.Background()
		rdb := common.RDB
		key := "cacheExplorer:" + fullPath
		n, _ := rdb.Exists(ctx, key).Result()
		if n <= 0 {
			// Cache doesn't exist
			localFilesPtr, readmeFileLink, err = getDataFromFS(path, fullPath)
			if err != nil {
				return
			}
			// Start a coroutine to update cache
			go func() {
				var values []string
				for _, f := range *localFilesPtr {
					s, err := json.Marshal(f)
					if err != nil {
						return
					}
					values = append(values, string(s))
				}
				rdb.RPush(ctx, key, values)
				rdb.Expire(ctx, key, time.Duration(common.ExplorerCacheTimeout)*time.Second)
			}()
		} else {
			// Cache existed, use cached data
			var localFiles []model.LocalFile
			file := model.LocalFile{}
			for _, s := range rdb.LRange(ctx, key, 0, -1).Val() {
				err = json.Unmarshal([]byte(s), &file)
				if err != nil {
					return
				}
				if file.Name == "README.md" {
					readmeFileLink = file.Link
				}
				localFiles = append(localFiles, file)
			}
			localFilesPtr = &localFiles
		}
	}
	return
}
