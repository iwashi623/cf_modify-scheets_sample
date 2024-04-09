package cf_spreadsheets

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func init() {
	functions.HTTP("modifySheets", modifySheets)
}

func modifySheets(w http.ResponseWriter, r *http.Request) {
	// X-Access-Tokenヘッダーからアクセストークンを取得
	accessToken := r.Header.Get("X-Access-Token")
	if accessToken == "" {
		http.Error(w, "missing access token", http.StatusUnauthorized)
		return
	}
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	// Google Drive APIのクライアントを初期化
	driveService, err := createDriveService(accessToken)
	if err != nil {
		http.Error(w, "Failed to create Drive service", http.StatusInternalServerError)
		return
	}

	// 共有シートをコピー
	fileId := "1_EbOFp26lePO-ibOmlnep8e-yUpXl_YCU3bl3JqZCyk" // 共有シートのID
	fr, err := driveService.Files.Get("root").Fields("id").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	folderId := fr.Id

	_, err = copyFile(driveService, fileId, folderId)
	if err != nil {
		http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Spreadsheet copied successfully")

}
func createDriveService(accessToken string) (*drive.Service, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	client := oauth2.NewClient(ctx, ts)
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	return driveService, err
}

func copyFile(driveService *drive.Service, fileId string, folderId string) (*drive.File, error) {
	copiedFile := &drive.File{
		Parents: []string{folderId},
	}
	file, err := driveService.Files.Copy(fileId, copiedFile).Do()
	return file, err
}
