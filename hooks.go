package gogitlab

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	project_url_hooks = "/projects/:id/hooks"          // Get list of project hooks
	project_url_hook  = "/projects/:id/hooks/:hook_id" // Get single project hook
)

type Hook struct {
	ID                    int    `json:"id,omitempty"`
	URL                   string `json:"url,omitempty"`
	CreatedAt             string `json:"created_at,omitempty"`
	PushEvents            *bool  `json:"push_events,omitempty"`
	IssuesEvents          *bool  `json:"issues_events,omitempty"`
	MergeRequestsEvents   *bool  `json:"merge_requests_events,omitempty"`
	TagPushEvents         *bool  `json:"tag_push_events,omitempty"`
	NoteEvents            *bool  `json:"note_events,omitempty"`
	BuildEvents           *bool  `json:"build_events,omitempty"`
	PipelineEvents        *bool  `json:"pipeline_events,omitempty"`
	WikiPageEvents        *bool  `json:"wiki_page_events,omitempty"`
	EnableSSLVerification *bool  `json:"enable_ssl_verification,omitempty"`
}

/*
Get list of project hooks.

    GET /projects/:id/hooks

Parameters:

    id The ID of a project

*/
func (g *Gitlab) ProjectHooks(id string) ([]*Hook, error) {
	var hooks []*Hook

	url, opaque := g.ResourceUrlRaw(project_url_hooks, map[string]string{":id": id})

	contents, err := g.buildAndExecRequestRaw("GET", url, opaque, nil)
	if err != nil {
		return hooks, err
	}

	err = json.Unmarshal(contents, &hooks)

	return hooks, err
}

/*
Get single project hook.

    GET /projects/:id/hooks/:hook_id

Parameters:

    id      The ID of a project
    hook_id The ID of a hook

*/
func (g *Gitlab) ProjectHook(id string, hook_id int) (*Hook, error) {
	hook := new(Hook)
	url, opaque := g.ResourceUrlRaw(project_url_hook, map[string]string{
		":id":      id,
		":hook_id": strconv.Itoa(hook_id),
	})

	contents, err := g.buildAndExecRequestRaw("GET", url, opaque, nil)
	if err != nil {
		return hook, err
	}

	err = json.Unmarshal(contents, &hook)

	return hook, err
}

/*
Add new project hook.

    POST /projects/:id/hooks

Parameters:

    id                    The ID or NAMESPACE/PROJECT_NAME of a project
    hook_url              The hook URL
    push_events           Trigger hook on push events
    issues_events         Trigger hook on issues events
    merge_requests_events Trigger hook on merge_requests events

*/
func (g *Gitlab) AddProjectHook(id string, hook *Hook) error {
	url, opaque := g.ResourceUrlRaw(project_url_hooks, map[string]string{":id": id})

	body := buildHookQuery(hook)
	_, err := g.buildAndExecRequestRaw("POST", url, opaque, []byte(body))

	return err
}

/*
Edit existing project hook.

    PUT /projects/:id/hooks/:hook_id

Parameters:

    id                    The ID or NAMESPACE/PROJECT_NAME of a project
    hook_id               The ID of a project hook
    hook_url              The hook URL
    push_events           Trigger hook on push events
    issues_events         Trigger hook on issues events
    merge_requests_events Trigger hook on merge_requests events

*/
func (g *Gitlab) EditProjectHook(id string, hook *Hook) error {
	url, opaque := g.ResourceUrlRaw(project_url_hook, map[string]string{
		":id":      id,
		":hook_id": strconv.Itoa(hook.ID),
	})

	body := buildHookQuery(hook)
	_, err := g.buildAndExecRequestRaw("PUT", url, opaque, []byte(body))

	return err
}

/*
Remove hook from project.

    DELETE /projects/:id/hooks/:hook_id

Parameters:

    id      The ID or NAMESPACE/PROJECT_NAME of a project
    hook_id The ID of hook to delete

*/
func (g *Gitlab) RemoveProjectHook(id string, hook_id int) error {
	url, opaque := g.ResourceUrlRaw(project_url_hook, map[string]string{
		":id":      id,
		":hook_id": strconv.Itoa(hook_id),
	})

	_, err := g.buildAndExecRequestRaw("DELETE", url, opaque, nil)

	return err
}

/*
Build HTTP query to add or edit hook
*/
func buildHookQuery(hook *Hook) string {
	v := url.Values{}

	hv := reflect.ValueOf(hook).Elem()
	ht := hv.Type()
	for i := 0; i < ht.NumField(); i++ {
		cf := ht.Field(i)
		if cf.Name == "ID" || cf.Name == "CreatedAt" {
			continue
		}
		val := hv.FieldByName(cf.Name)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
		}
		v.Set(strings.Split(cf.Tag.Get("json"), ",")[0], fmt.Sprintf("%v", val))
	}
	return v.Encode()
}
