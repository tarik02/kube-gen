package kubegen

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strconv"
	"text/template"

	kapi "k8s.io/api/core/v1"
)

var Funcs = template.FuncMap{
	"allPodsReady":  allPodsReady,
	"anyPodReady":   anyPodReady,
	"dir":           dirList,
	"exists":        exists,
	"groupBy":       groupBy,
	"groupByKeys":   groupByKeys,
	"groupByMulti":  groupByMulti,
	"intersect":     intersect,
	"isPodReady":    isPodReady,
	"isValidJson":   isValidJSON,
	"pathJoin":      filepath.Join,
	"pathJoinSlice": pathJoinSlice,
	"parseBool":     strconv.ParseBool,
	"readyPods":     readyPods,
	"shell":         execShell,
	"when":          when,
	"where":         where,
	"whereExist":    whereExist,
	"whereNotExist": whereNotExist,
	"whereAny":      whereAny,
	"whereAll":      whereAll,
}

func pathJoinSlice(input []string) string {
	return filepath.Join(input...)
}

func unmarshalJSON(input string) (any, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return nil, err
	}
	return v, nil
}

func isValidJSON(input string) bool {
	_, err := unmarshalJSON(input)
	return err == nil
}

type ShellResult struct {
	Success bool
	Stdout  string
	Stderr  string
}

func execShell(cs string) *ShellResult {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd := exec.Command(shellExe, shellArg, cs)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := &ShellResult{
		Success: err == nil,
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
	}
	return res
}

func isPodReady(i any) bool {
	if p, ok := i.(kapi.Pod); ok {
		return IsPodReady(&p)
	} else if p, ok := i.(*kapi.Pod); ok {
		return IsPodReady(p)
	}
	return false
}

func allPodsReady(pods []kapi.Pod) bool {
	for _, p := range pods {
		if !isPodReady(p) {
			return false
		}
	}
	return true
}

func anyPodReady(pods []kapi.Pod) bool {
	for _, p := range pods {
		if isPodReady(p) {
			return true
		}
	}
	return false
}

func readyPods(pods []kapi.Pod) []kapi.Pod {
	var ready []kapi.Pod
	for _, p := range pods {
		if isPodReady(p) {
			ready = append(ready, p)
		}
	}
	return ready
}
