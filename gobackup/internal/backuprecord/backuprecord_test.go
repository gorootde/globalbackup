package backuprecord

import "testing"

func TestGetRsyncSignature(t *testing.T) {
	sig, err := getRsyncSignature("../../../testdata/loremipsum.txt")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Signature len: %v  Error: %v", len(sig), err)
}
