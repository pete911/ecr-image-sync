// +build gofuzzbeta

package docker

import (
	"testing"
)

func FuzzImage(f *testing.F) {

	f.Fuzz(func(t *testing.T, imageName string) {

		image, err := ToImage(imageName)
		if err != nil {
			t.Skip()
		}

		imageName2 := image.String()
		if imageName != imageName2 {
			t.Errorf("ToImage returned different image names\nbefore: %q\nafter: %q", imageName, imageName2)
		}
	})
}
