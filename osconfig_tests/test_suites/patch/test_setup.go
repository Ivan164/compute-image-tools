//  Copyright 2019 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package patch

import (
	"time"

	"github.com/GoogleCloudPlatform/compute-image-tools/osconfig_tests/utils"
	"google.golang.org/api/compute/v1"
)

type patchTestSetup struct {
	testName      string
	image         string
	startup       *compute.MetadataItems
	assertTimeout time.Duration
	machineType   string
}

var (
	windowsRecordBoot = `
while ($true) {
  $uri = 'http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/osconfig_tests/boot_count'
  $old = Invoke-RestMethod -Method GET -Uri $uri -Headers @{"Metadata-Flavor" = "Google"}
  $new = $old+1
  try {
	Invoke-RestMethod -Method PUT -Uri $uri -Headers @{"Metadata-Flavor" = "Google"} -Body $new -ErrorAction Stop
  }
  catch {
	Write-Output $_.Exception.Message
	Start-Sleep 1
    continue
  }
  break
}
`
	windowsStartup = windowsRecordBoot + utils.InstallOSConfigGooGet

	linuxRecordBoot = `
uri=http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/osconfig_tests/boot_count
old=$(curl $uri -H "Metadata-Flavor: Google" -f)
new=$(($old + 1))
curl -X PUT --data "${new}" $uri -H "Metadata-Flavor: Google"
`
	aptStartup = linuxRecordBoot + utils.InstallOSConfigDeb
	el6Startup = linuxRecordBoot + utils.InstallOSConfigYumEL6
	el7Startup = linuxRecordBoot + "yum install -y yum-utils\n" + utils.InstallOSConfigYumEL7

	windowsSetup = &patchTestSetup{
		assertTimeout: 60 * time.Minute,
		startup: &compute.MetadataItems{
			Key:   "windows-startup-script-ps1",
			Value: &windowsStartup,
		},
		machineType: "n1-standard-4",
	}
	aptSetup = &patchTestSetup{
		assertTimeout: 5 * time.Minute,
		startup: &compute.MetadataItems{
			Key:   "startup-script",
			Value: &aptStartup,
		},
		machineType: "n1-standard-2",
	}
	el6Setup = &patchTestSetup{
		assertTimeout: 5 * time.Minute,
		startup: &compute.MetadataItems{
			Key:   "startup-script",
			Value: &el6Startup,
		},
		machineType: "n1-standard-2",
	}
	el7Setup = &patchTestSetup{
		assertTimeout: 5 * time.Minute,
		startup: &compute.MetadataItems{
			Key:   "startup-script",
			Value: &el7Startup,
		},
		machineType: "n1-standard-2",
	}
	el8Setup = &patchTestSetup{
		assertTimeout: 5 * time.Minute,
		startup: &compute.MetadataItems{
			Key:   "startup-script",
			Value: &el7Startup,
		},
		machineType: "n1-standard-2",
	}
)

func imageTestSetup(mapping map[*patchTestSetup]map[string]string) (setup []*patchTestSetup) {
	for s, m := range mapping {
		for name, image := range m {
			new := patchTestSetup(*s)
			new.testName = name
			new.image = image
			setup = append(setup, &new)
		}
	}
	return
}

func headImageTestSetup() []*patchTestSetup {
	// This maps a specific patchTestSetup to test setup names and associated images.
	mapping := map[*patchTestSetup]map[string]string{
		windowsSetup: utils.HeadWindowsImages,
		el6Setup:     utils.HeadEL6Images,
		el7Setup:     utils.HeadEL7Images,
		el8Setup:     utils.HeadEL8Images,
		aptSetup:     utils.HeadAptImages,
	}

	return imageTestSetup(mapping)
}

func oldImageTestSetup() []*patchTestSetup {
	// This maps a specific patchTestSetup to test setup names and associated images.
	mapping := map[*patchTestSetup]map[string]string{
		windowsSetup: utils.OldWindowsImages,
		el6Setup:     utils.OldEL6Images,
		el7Setup:     utils.OldEL7Images,
		aptSetup:     utils.OldAptImages,
	}

	return imageTestSetup(mapping)
}

func aptHeadImageTestSetup() []*patchTestSetup {
	// This maps a specific patchTestSetup to test setup names and associated images.
	mapping := map[*patchTestSetup]map[string]string{
		aptSetup: utils.HeadAptImages,
	}

	return imageTestSetup(mapping)
}

func yumHeadImageTestSetup() []*patchTestSetup {
	// This maps a specific patchTestSetup to test setup names and associated images.
	mapping := map[*patchTestSetup]map[string]string{
		el6Setup: utils.HeadEL6Images,
		el7Setup: utils.HeadEL7Images,
		el8Setup: utils.HeadEL8Images,
	}

	return imageTestSetup(mapping)
}