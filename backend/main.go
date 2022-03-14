package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// レスポンス用
// type Response struct {
// 	Objects []ObjInfo `json:"objects"`
// }

// レスポンス用
type ObjInfo struct {
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
}

func main() {
	e := echo.New()

	e.GET("/static", getObjects)
	e.POST("/upload", storeObject)

	e.Logger.Fatal(e.Start(":5000"))
}

// 特定バケット配下のオブジェクト一覧
func getObjects(c echo.Context) error {
	sess := createSession()
	svc := s3.New(sess)
	bucket := c.QueryParam("bucket")

	// オブジェクト取得
	res, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	// 本来はもっときちんとエラーハンドリングした方が良いが、簡単のため今回はこれで良しとする
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var objects []ObjInfo
	// オブジェクト情報をJSONにつめて返す
	for _, content := range res.Contents {
		objects = append(
			objects,
			ObjInfo{Key: *content.Key, LastModified: *content.LastModified, Size: *content.Size},
		)
	}
	return c.JSON(http.StatusOK, objects)
}

func storeObject(c echo.Context) error {
	// 画像ファイル取得
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 名前を変更
	fileModel := strings.Split(file.Filename, ".")
	fileName := uuid.New().String()
	fileExtension := fileModel[1]

	log.Println("got file", fileModel, fileName, fileExtension)

	// AWS設定
	sess := createSession()

	bucket := "static"
	objectKey := fileName + "." + fileExtension

	// アップロード
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   src,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "done")
}

// セッションを返す
func createSession() *session.Session {
	// 特に設定しなくても環境変数にセットしたクレデンシャル情報を利用して接続してくれる
	cfg := aws.Config{
		Region:           aws.String("ap-northeast-1"),
		Endpoint:         aws.String("http://minio:9000"), // コンテナ内からアクセスする場合はホストをサービス名で指定
		S3ForcePathStyle: aws.Bool(true),                  // ローカルで動かす場合は必須
	}
	return session.Must(session.NewSession(&cfg))
}
