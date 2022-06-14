package fieldmask_test

import (
	"testing"

	testv1 "github.com/srikrsna/fieldmask-go/internal/gen/test/v1"
	"github.com/srikrsna/fieldmask-go/internal/gen/test/v1/testv1fieldmask"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestFieldMasks(t *testing.T) {
	_, err := fieldmaskpb.New(
		&testv1.All{},
		testv1fieldmask.AllMask.E(),
		testv1fieldmask.AllMask.O().Bl(),
		testv1fieldmask.AllMask.O().By(),
		string(testv1fieldmask.AllMask.N()),
	)
	if err != nil {
		t.Fatal(err)
	}
}
