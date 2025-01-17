// Copyright 2021 The Audit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this File except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package custom

import (
	"fmt"
	"sort"

	"github.com/operator-framework/audit/pkg"
	"github.com/operator-framework/audit/pkg/reports/bundles"
)

type Migrated struct {
	Name            string
	Kinds           []string
	Bundles         []string
	Channels        []string
	BundlesMigrated []string
	AllBundles      []bundles.Column
}

type NotMigrated struct {
	Name            string
	Kinds           []string
	Channels        []string
	Bundles         []string
	BundlesMigrated []string
	AllBundles      []bundles.Column
}

type APIDashReport struct {
	ImageName   string
	ImageID     string
	ImageHash   string
	ImageBuild  string
	Migrated    []Migrated
	NotMigrated []NotMigrated
	GeneratedAt string
}

// NewAPIDashReport returns the structure to render the Deprecate API custom dashboard
// nolint:dupl
func NewAPIDashReport(bundlesReport bundles.Report) *APIDashReport {
	apiDash := APIDashReport{}
	apiDash.ImageName = bundlesReport.Flags.IndexImage
	apiDash.ImageID = bundlesReport.IndexImageInspect.ID
	apiDash.ImageBuild = bundlesReport.IndexImageInspect.Created
	apiDash.GeneratedAt = bundlesReport.GenerateAt

	mapPackagesWithBundles := MapBundlesPerPackage(bundlesReport)
	migrated := MapPkgsComplyingWithDeprecateAPI122(mapPackagesWithBundles)
	notMigrated := make(map[string][]bundles.Column)
	for key := range mapPackagesWithBundles {
		if len(migrated[key]) == 0 {
			notMigrated[key] = mapPackagesWithBundles[key]
		}
	}

	for k, bundles := range migrated {
		kinds, channels, bundlesNotMigrated, bundlesMigrated := getReportValues(bundles)
		apiDash.Migrated = append(apiDash.Migrated, Migrated{
			Name:            k,
			Kinds:           pkg.GetUniqueValues(kinds),
			Channels:        pkg.GetUniqueValues(channels),
			Bundles:         bundlesNotMigrated,
			BundlesMigrated: bundlesMigrated,
			AllBundles:      bundles,
		})
	}

	for k, bundles := range notMigrated {
		kinds, channels, bundlesNotMigrated, bundlesMigrated := getReportValues(bundles)
		apiDash.NotMigrated = append(apiDash.NotMigrated, NotMigrated{
			Name:            k,
			Kinds:           pkg.GetUniqueValues(kinds),
			Channels:        pkg.GetUniqueValues(channels),
			Bundles:         bundlesNotMigrated,
			BundlesMigrated: bundlesMigrated,
			AllBundles:      bundles,
		})
	}

	return &apiDash

}

func getReportValues(bundles []bundles.Column) ([]string, []string, []string, []string) {
	var kinds []string
	var channels []string
	for _, b := range bundles {
		kinds = append(kinds, b.KindsDeprecateAPIs...)
	}
	for _, b := range bundles {
		channels = append(channels, b.Channels...)
	}
	var bundlesNotMigrated []string
	var bundlesMigrated []string
	for _, b := range bundles {
		if len(b.KindsDeprecateAPIs) > 0 {
			bundlesNotMigrated = append(bundlesNotMigrated, buildBundleString(b))
		} else {
			bundlesMigrated = append(bundlesMigrated, buildBundleString(b))
		}
	}

	sort.Slice(bundlesNotMigrated[:], func(i, j int) bool {
		return bundlesNotMigrated[i] < bundlesNotMigrated[j]
	})

	sort.Slice(bundlesMigrated[:], func(i, j int) bool {
		return bundlesMigrated[i] < bundlesMigrated[j]
	})

	return kinds, channels, bundlesNotMigrated, bundlesMigrated
}

func buildBundleString(b bundles.Column) string {
	return fmt.Sprintf("%s - (label=%s,max=%s,channels=%s,head:%s,defaultChannel:%s, deprecated:%s)",
		b.BundleName,
		b.OCPLabel,
		GetMaxOCPValue(b),
		pkg.GetUniqueValues(b.Channels),
		pkg.GetYesOrNo(b.IsHeadOfChannel),
		pkg.GetYesOrNo(b.IsFromDefaultChannel),
		pkg.GetYesOrNo(b.IsDeprecated),
	)
}
