package volume

import "testing"

func TestAddFile(t *testing.T) {
	testling := New()
	t.Log("Testling filename is ", testling.Path)
	bytes, err := Add(testling, "../../../testdata/samepic.jpg")
	if err != nil {
		t.Error(err)
	}

	if bytes <= 0 {
		t.Errorf("Expected more than 0 bytes written but was %v", bytes)
	}

	_, err2 := Add(testling, "../../../testdata/London Night Lights Aerial 4K Desktop Wallpaper.jpg")
	if err2 != nil {
		t.Error(err)
	}

}
