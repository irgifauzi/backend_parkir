package model

type Ghcreates struct {
	GitHubAccessToken string `bson:"githubaccesstoken,omitempty" json:"githubaccesstoken,omitempty"`
	GitHubAuthorName  string `bson:"githubauthorname,omitempty" json:"githubauthorname,omitempty"`
	GitHubAuthorEmail string `bson:"githubauthoremail,omitempty" json:"githubauthoremail,omitempty"`
}