package controller

import (
	"fmt"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/ghupload"
	"github.com/whatsauth/itmodel"
)

func PostUploadGithub(w http.ResponseWriter, r *http.Request) {
	var respn itmodel.Response

	fmt.Println("Starting file upload process")

	// Parse the form file
	_, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error parsing form file:", err)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	// Get the folder parameter
	folder := helper.GetParam(r)
	var pathFile string
	if folder != "" {
		pathFile = folder + "/" + header.Filename
	} else {
		pathFile = header.Filename
	}

	// Fetch GitHub credentials from the database
	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		fmt.Println("Error fetching GitHub credentials:", err)
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	// Upload the file to GitHub
	content, _, err := ghupload.GithubUpload(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, header, "parkirgratis", "parkirgratis.github.io", pathFile, false)
	if err != nil {
		fmt.Println("Error uploading file to GitHub:", err)
		respn.Info = "gagal upload github"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusEarlyHints, respn) // Changed `content` to `respn`
		return
	}

	// Check if content is nil
	if content == nil || content.Content == nil {
		fmt.Println("Error: content or content.Content is nil")
		respn.Response = "Error uploading file"
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = *content.Content.Name
	respn.Response = *content.Content.Path
	helper.WriteJSON(w, http.StatusOK, respn)
	fmt.Println("File upload process completed successfully")
}
