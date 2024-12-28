package handlers

import (
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
)

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件
// @Tags oss
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Router /upload/file [post]
func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|UploadFile|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		ctx.Abort()
	}

	endpiont := "localhost:9000"
	accessKey := "admin"
	secretKey := "admin@123"
	bucketName := "kasha"

	minioClient, err := minio.New(endpiont, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|UploadFile|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.RequestFail, "上传失败", nil))
		ctx.Abort()
	}
	// 检查MinIO服务器是否在线
	found, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|UploadFile|%v", err)
		panic(err)
	}
	if !found {
		logx.GetLogger("ShopManage").Errorf("Handler|UploadFile|BucketNotExist")
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.RequestFail, "BucketNotExist", nil))
		ctx.Abort()
	}

	// 打开文件
	f, err := file.Open()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|UploadFile|%v", err)
		panic(err)
	}
	defer f.Close()

	// 上传文件
	_, err = minioClient.PutObject(ctx, bucketName, file.Filename, f, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|UploadFile|%v", err)
		panic(err)
	}

	// 构建文件的URL
	fileURL := fmt.Sprintf("%s//%s/%s/%s", "http:", endpiont, bucketName, file.Filename)

	ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.RequestSuccess, "上传成功", map[string]string{
		"url": fileURL,
	}))

}