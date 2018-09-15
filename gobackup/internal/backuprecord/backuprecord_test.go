package backuprecord

import "testing"

func TestGetChecksum(t *testing.T) {
	sig, err := getChecksum("../../../testdata/loremipsum.txt")
	if err != nil {
		t.Error(err)
	}
	if len(sig) <= 0 {
		t.Errorf("Signature is zero")
	}
	t.Logf("Signature len: %v", len(sig))
}

func TestFileSlightlyChangedDifferentChecksums(t *testing.T) {
	sig1, err := getChecksum("../../../testdata/loremipsum.txt")
	if err != nil {
		t.Error(err)
	}
	sig2, err := getChecksum("../../../testdata/loremipsum_singlecharmod.txt")
	if err != nil {
		t.Error(err)
	}

	if sig1 == sig2 {
		t.Errorf("Signatures do not differ")
	}
}
