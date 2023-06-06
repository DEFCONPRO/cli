package testserver

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	flinkv2 "github.com/confluentinc/ccloud-sdk-go-v2/flink/v2"
)

func handleFcpmComputePools(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var v any

		switch r.Method {
		case http.MethodGet:
			v = flinkv2.FcpmV2ComputePoolList{Data: []flinkv2.FcpmV2ComputePool{
				{
					Id: flinkv2.PtrString("lfcp-123456"),
					Spec: &flinkv2.FcpmV2ComputePoolSpec{
						DisplayName: flinkv2.PtrString("my-compute-pool-1"),
						MaxCfu:      flinkv2.PtrInt32(1),
						Region:      flinkv2.PtrString("us-west-2"),
					},
					Status: &flinkv2.FcpmV2ComputePoolStatus{Phase: "PROVISIONED"},
				},
				{
					Id: flinkv2.PtrString("lfcp-222222"),
					Spec: &flinkv2.FcpmV2ComputePoolSpec{
						DisplayName: flinkv2.PtrString("my-compute-pool-2"),
						MaxCfu:      flinkv2.PtrInt32(1),
						Region:      flinkv2.PtrString("us-west-2"),
					},
					Status: &flinkv2.FcpmV2ComputePoolStatus{Phase: "PROVISIONED"},
				},
			}}
		case http.MethodPost:
			create := new(flinkv2.FcpmV2ComputePool)
			err := json.NewDecoder(r.Body).Decode(create)
			require.NoError(t, err)

			v = flinkv2.FcpmV2ComputePool{
				Id:     flinkv2.PtrString("lfcp-123456"),
				Spec:   create.Spec,
				Status: &flinkv2.FcpmV2ComputePoolStatus{Phase: "PROVISIONING"},
			}
		}

		err := json.NewEncoder(w).Encode(v)
		require.NoError(t, err)
	}
}

func handleFcpmComputePoolsId(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var computePool flinkv2.FcpmV2ComputePool
		id := mux.Vars(r)["id"]

		switch r.Method {
		case http.MethodGet:
			computePool = flinkv2.FcpmV2ComputePool{
				Id: flinkv2.PtrString(id),
				Spec: &flinkv2.FcpmV2ComputePoolSpec{
					DisplayName:  flinkv2.PtrString("my-compute-pool-1"),
					HttpEndpoint: flinkv2.PtrString(TestFlinkGatewayUrl.String()),
					MaxCfu:       flinkv2.PtrInt32(1),
					Region:       flinkv2.PtrString("us-west-2"),
				},
				Status: &flinkv2.FcpmV2ComputePoolStatus{Phase: "PROVISIONED"},
			}
		case http.MethodPatch:
			update := new(flinkv2.FcpmV2ComputePool)
			err := json.NewDecoder(r.Body).Decode(update)
			require.NoError(t, err)

			computePool = flinkv2.FcpmV2ComputePool{
				Id: flinkv2.PtrString(id),
				Spec: &flinkv2.FcpmV2ComputePoolSpec{
					DisplayName: flinkv2.PtrString("my-compute-pool-1"),
					MaxCfu:      flinkv2.PtrInt32(update.Spec.GetMaxCfu()),
					Region:      flinkv2.PtrString("us-west-2"),
				},
				Status: &flinkv2.FcpmV2ComputePoolStatus{Phase: "PROVISIONED"},
			}
		}

		err := json.NewEncoder(w).Encode(computePool)
		require.NoError(t, err)
	}
}

func handleFcpmRegions(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aws := flinkv2.FcpmV2Region{
			DisplayName: flinkv2.PtrString("Europe (eu-west-1)"),
			Cloud:       flinkv2.PtrString("AWS"),
			RegionName:  flinkv2.PtrString("eu-west-1"),
		}
		gcp := flinkv2.FcpmV2Region{
			DisplayName: flinkv2.PtrString("Frankfurt (europe-west3-a)"),
			Cloud:       flinkv2.PtrString("GCP"),
			RegionName:  flinkv2.PtrString("europe-west3-a"),
		}

		regions := []flinkv2.FcpmV2Region{aws, gcp}
		if r.URL.Query().Get("cloud") == "AWS" {
			regions = []flinkv2.FcpmV2Region{aws}
		}

		err := json.NewEncoder(w).Encode(flinkv2.FcpmV2RegionList{Data: regions})
		require.NoError(t, err)
	}
}