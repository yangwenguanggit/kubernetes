/*
Copyright 2019 The Kubernetes Authors.

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

package plugins

import (
	"fmt"
	"k8s.io/api/core/v1"
)

const (
	// CinderDriverName is the name of the CSI driver for Cinder
	CinderDriverName = "cinder.csi.openstack.org"
	// CinderInTreePluginName is the name of the intree plugin for Cinder
	CinderInTreePluginName = "kubernetes.io/cinder"
)

var _ InTreePlugin = (*osCinderCSITranslator)(nil)

// osCinderCSITranslator handles translation of PV spec from In-tree Cinder to CSI Cinder and vice versa
type osCinderCSITranslator struct{}

// NewOpenStackCinderCSITranslator returns a new instance of osCinderCSITranslator
func NewOpenStackCinderCSITranslator() InTreePlugin {
	return &osCinderCSITranslator{}
}

// TranslateInTreeStorageClassParametersToCSI translates InTree Cinder storage class parameters to CSI storage class
func (t *osCinderCSITranslator) TranslateInTreeStorageClassParametersToCSI(scParameters map[string]string) (map[string]string, error) {
	return scParameters, nil
}

// TranslateInTreePVToCSI takes a PV with Cinder set from in-tree
// and converts the Cinder source to a CSIPersistentVolumeSource
func (t *osCinderCSITranslator) TranslateInTreePVToCSI(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	if pv == nil || pv.Spec.Cinder == nil {
		return nil, fmt.Errorf("pv is nil or Cinder not defined on pv")
	}

	cinderSource := pv.Spec.Cinder

	csiSource := &v1.CSIPersistentVolumeSource{
		Driver:           CinderDriverName,
		VolumeHandle:     cinderSource.VolumeID,
		ReadOnly:         cinderSource.ReadOnly,
		FSType:           cinderSource.FSType,
		VolumeAttributes: map[string]string{},
	}

	pv.Spec.Cinder = nil
	pv.Spec.CSI = csiSource
	return pv, nil
}

// TranslateCSIPVToInTree takes a PV with CSIPersistentVolumeSource set and
// translates the Cinder CSI source to a Cinder In-tree source.
func (t *osCinderCSITranslator) TranslateCSIPVToInTree(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	if pv == nil || pv.Spec.CSI == nil {
		return nil, fmt.Errorf("pv is nil or CSI source not defined on pv")
	}

	csiSource := pv.Spec.CSI

	cinderSource := &v1.CinderPersistentVolumeSource{
		VolumeID: csiSource.VolumeHandle,
		FSType:   csiSource.FSType,
		ReadOnly: csiSource.ReadOnly,
	}

	pv.Spec.CSI = nil
	pv.Spec.Cinder = cinderSource
	return pv, nil
}

// CanSupport tests whether the plugin supports a given volume
// specification from the API.  The spec pointer should be considered
// const.
func (t *osCinderCSITranslator) CanSupport(pv *v1.PersistentVolume) bool {
	return pv != nil && pv.Spec.Cinder != nil
}

// GetInTreePluginName returns the name of the intree plugin driver
func (t *osCinderCSITranslator) GetInTreePluginName() string {
	return CinderInTreePluginName
}
