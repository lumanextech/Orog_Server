package defi_quotation_v1

import (
	xhttp "github.com/zeromicro/x/http"
	"net/http"

	"github.com/simance-ai/smdx/app/internal/logic/defi_quotation_v1"
	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetXTopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetXTopRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := defi_quotation_v1.NewGetXTopLogic(r.Context(), svcCtx)
		resp, err := l.GetXTop(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
