/*
 Copyright 2022. The KubeVela Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package util

import (
	"encoding/json"

	wfContext "github.com/kubevela/workflow/pkg/context"
	"github.com/kubevela/workflow/pkg/cue/model"
	"github.com/kubevela/workflow/pkg/cue/model/value"
	monitorContext "github.com/kubevela/workflow/pkg/monitor/context"
	"github.com/kubevela/workflow/pkg/types"
)

const (
	// ProviderName is provider name for install.
	ProviderName = "util"
)

type provider struct{}

func (p *provider) PatchK8sObject(ctx monitorContext.Context, wfCtx wfContext.Context, v *value.Value, act types.Action) error {
	val, err := v.LookupValue("value")
	if err != nil {
		return err
	}
	pv, err := v.LookupValue("patch")
	if err != nil {
		return err
	}
	base, err := model.NewBase(val.CueValue())
	if err != nil {
		return err
	}
	if err = base.Unify(pv.CueValue()); err != nil {
		return v.FillObject(err, "err")
	}

	workload, err := base.Unstructured()
	if err != nil {
		return v.FillObject(err, "err")
	}
	return v.FillObject(workload.Object, "result")
}

// String convert byte to string
func (p *provider) String(ctx monitorContext.Context, wfCtx wfContext.Context, v *value.Value, act types.Action) error {
	b, err := v.LookupValue("bt")
	if err != nil {
		return err
	}
	s, err := b.CueValue().Bytes()
	if err != nil {
		return err
	}
	return v.FillObject(string(s), "str")
}

// Log print cue value in log
func (p *provider) Log(ctx monitorContext.Context, wfCtx wfContext.Context, v *value.Value, act types.Action) error {
	data, err := v.LookupValue("data")
	if err != nil {
		return err
	}
	logCtx := ctx.Fork("cue logs")
	if s, err := data.GetString(); err == nil {
		logCtx.Info(s)
		return nil
	}
	var tmp interface{}
	if err := data.UnmarshalTo(&tmp); err != nil {
		return err
	}
	b, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	logCtx.Info(string(b))
	return nil
}

// Install register handlers to provider discover.
func Install(p types.Providers) {
	prd := &provider{}
	p.Register(ProviderName, map[string]types.Handler{
		"patch-k8s-object": prd.PatchK8sObject,
		"string":           prd.String,
		"log":              prd.Log,
	})
}
