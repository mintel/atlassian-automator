// Code generated by interfacer; DO NOT EDIT

package confluenceclient

import (
	"io"
	"net/http"
	"net/url"

	goconfluence "github.com/virtomize/confluence-go-api"
)

// Client is an interface generated for "github.com/virtomize/confluence-go-api.API".
type Client interface {
	AddLabels(string, *[]goconfluence.Label) (*goconfluence.Labels, error)
	AnonymousUser() (*goconfluence.User, error)
	Auth(*http.Request)
	CreateContent(*goconfluence.Content) (*goconfluence.Content, error)
	CurrentUser() (*goconfluence.User, error)
	DelContent(string) (*goconfluence.Content, error)
	DeleteLabel(string, string) (*goconfluence.Labels, error)
	GetAllSpaces(goconfluence.AllSpacesQuery) (*goconfluence.AllSpaces, error)
	GetAttachments(string) (*goconfluence.Search, error)
	GetBlueprintTemplates(goconfluence.TemplateQuery) (*goconfluence.TemplateSearch, error)
	GetChildPages(string) (*goconfluence.Search, error)
	GetComments(string) (*goconfluence.Search, error)
	GetContent(goconfluence.ContentQuery) (*goconfluence.ContentSearch, error)
	GetContentByID(string, goconfluence.ContentQuery) (*goconfluence.Content, error)
	GetContentTemplates(goconfluence.TemplateQuery) (*goconfluence.TemplateSearch, error)
	GetContentVersion(string) (*goconfluence.ContentVersionResult, error)
	GetHistory(string) (*goconfluence.History, error)
	GetLabels(string) (*goconfluence.Labels, error)
	GetWatchers(string) (*goconfluence.Watchers, error)
	Request(*http.Request) ([]byte, error)
	Search(goconfluence.SearchQuery) (*goconfluence.Search, error)
	SendAllSpacesRequest(*url.URL, string) (*goconfluence.AllSpaces, error)
	SendContentAttachmentRequest(*url.URL, string, io.Reader, map[string]string) (*goconfluence.Search, error)
	SendContentRequest(*url.URL, string, *goconfluence.Content) (*goconfluence.Content, error)
	SendContentVersionRequest(*url.URL, string) (*goconfluence.ContentVersionResult, error)
	SendHistoryRequest(*url.URL, string) (*goconfluence.History, error)
	SendLabelRequest(*url.URL, string, *[]goconfluence.Label) (*goconfluence.Labels, error)
	SendSearchRequest(*url.URL, string) (*goconfluence.Search, error)
	SendUserRequest(*url.URL, string) (*goconfluence.User, error)
	SendWatcherRequest(*url.URL, string) (*goconfluence.Watchers, error)
	UpdateAttachment(string, string, string, io.Reader) (*goconfluence.Search, error)
	UpdateContent(*goconfluence.Content) (*goconfluence.Content, error)
	UploadAttachment(string, string, io.Reader) (*goconfluence.Search, error)
	User(string) (*goconfluence.User, error)
	VerifyTLS(bool)
}
