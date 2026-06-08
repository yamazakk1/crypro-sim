package handler

import (
    "encoding/json"
    "log"
    "net/http"

    pbAsset "crypto-simulator/pkg/pb/asset"

    "google.golang.org/grpc/status"
)

func (h *GatewayHandler) HandleListAssets(w http.ResponseWriter, r *http.Request) {
    log.Println("gateway: HandleListAssets called")
    resp, err := h.AssetServiceClient.ListAssets(r.Context(), &pbAsset.ListAssetsRequest{})
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleListAssets failed: %v", err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    var assets []Asset
    for _, a := range resp.GetAssets() {
        assets = append(assets, Asset{
            Id:       a.GetId(),
            Symbol:   a.GetSymbol(),
            FullName: a.GetFullname(),
            IsActive: a.GetIsActive(),
        })
    }

    log.Printf("gateway: HandleListAssets success, returned %d assets", len(assets))
    writeJson(w, http.StatusOK, ListAssetsResponse{Assets: assets})
}

func (h *GatewayHandler) HandleGetAsset(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    log.Printf("gateway: HandleGetAsset called, id=%s", id)
    if id == "" {
        log.Println("gateway: HandleGetAsset: empty id")
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
        return
    }

    resp, err := h.AssetServiceClient.GetAsset(r.Context(), &pbAsset.GetAssetRequest{Id: id})
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleGetAsset failed for id=%s: %v", id, err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    log.Printf("gateway: HandleGetAsset success for id=%s", id)
    writeJson(w, http.StatusOK, Asset{
        Id:       resp.GetId(),
        Symbol:   resp.GetSymbol(),
        FullName: resp.GetFullname(),
        IsActive: resp.GetIsActive(),
    })
}

func (h *GatewayHandler) HandleCreateAsset(w http.ResponseWriter, r *http.Request) {
    log.Println("gateway: HandleCreateAsset called")
    if !h.isAdmin(r) {
        log.Println("gateway: HandleCreateAsset: forbidden (not admin)")
        writeJson(w, http.StatusForbidden, ErrorResponse{Error: "admin only"})
        return
    }

    var req CreateAssetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("gateway: HandleCreateAsset: invalid body: %v", err)
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
        return
    }

    if req.Symbol == "" || req.FullName == "" {
        log.Println("gateway: HandleCreateAsset: empty fields")
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "symbol and full_name are required"})
        return
    }

   resp, err := h.AssetServiceClient.CreateAsset(r.Context(), &pbAsset.CreateAssetRequest{
        Symbol:       req.Symbol,
        Fullname:     req.FullName,
        InitialPrice: float32(req.InitialPrice),  // ← должно быть!
    })
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleCreateAsset failed: %v", err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    log.Printf("gateway: HandleCreateAsset success: id=%s, symbol=%s", resp.GetId(), resp.GetSymbol())
    writeJson(w, http.StatusCreated, CreateAssetResponse{
        Id:       resp.GetId(),
        Symbol:   resp.GetSymbol(),
        FullName: resp.GetFullname(),
        IsActive: resp.GetIsActive(),
    })
}

func (h *GatewayHandler) HandleUpdateAsset(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    log.Printf("gateway: HandleUpdateAsset called, id=%s", id)
    if !h.isAdmin(r) {
        log.Println("gateway: HandleUpdateAsset: forbidden (not admin)")
        writeJson(w, http.StatusForbidden, ErrorResponse{Error: "admin only"})
        return
    }

    if id == "" {
        log.Println("gateway: HandleUpdateAsset: empty id")
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
        return
    }

    var req UpdateAssetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("gateway: HandleUpdateAsset: invalid body: %v", err)
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
        return
    }

    _, err := h.AssetServiceClient.UpdateAsset(r.Context(), &pbAsset.UpdateAssetRequest{
        Id:       id,
        Fullname: req.FullName,
    })
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleUpdateAsset failed for id=%s: %v", id, err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    log.Printf("gateway: HandleUpdateAsset success for id=%s", id)
    writeJson(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *GatewayHandler) HandleDeactivateAsset(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    log.Printf("gateway: HandleDeactivateAsset called, id=%s", id)
    if !h.isAdmin(r) {
        log.Println("gateway: HandleDeactivateAsset: forbidden (not admin)")
        writeJson(w, http.StatusForbidden, ErrorResponse{Error: "admin only"})
        return
    }

    if id == "" {
        log.Println("gateway: HandleDeactivateAsset: empty id")
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
        return
    }

    _, err := h.AssetServiceClient.DeactivateAsset(r.Context(), &pbAsset.DeactivateAssetRequest{Id: id})
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleDeactivateAsset failed for id=%s: %v", id, err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    log.Printf("gateway: HandleDeactivateAsset success for id=%s", id)
    writeJson(w, http.StatusOK, map[string]string{"status": "deactivated"})
}

func (h *GatewayHandler) HandleActivateAsset(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    log.Printf("gateway: HandleActivateAsset called, id=%s", id)
    if !h.isAdmin(r) {
        log.Println("gateway: HandleActivateAsset: forbidden (not admin)")
        writeJson(w, http.StatusForbidden, ErrorResponse{Error: "admin only"})
        return
    }

    if id == "" {
        log.Println("gateway: HandleActivateAsset: empty id")
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
        return
    }

    _, err := h.AssetServiceClient.ActivateAsset(r.Context(), &pbAsset.ActivateAssetRequest{Id: id})
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleActivateAsset failed for id=%s: %v", id, err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    log.Printf("gateway: HandleActivateAsset success for id=%s", id)
    writeJson(w, http.StatusOK, map[string]string{"status": "activated"})
}

func (h *GatewayHandler) HandleDeleteAsset(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    log.Printf("gateway: HandleDeleteAsset called, id=%s", id)
    if !h.isAdmin(r) {
        log.Println("gateway: HandleDeleteAsset: forbidden (not admin)")
        writeJson(w, http.StatusForbidden, ErrorResponse{Error: "admin only"})
        return
    }

    if id == "" {
        log.Println("gateway: HandleDeleteAsset: empty id")
        writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
        return
    }

    resp, err := h.AssetServiceClient.DeleteAsset(r.Context(), &pbAsset.DeleteAssetRequest{Id: id})
    if err != nil {
        st := status.Convert(err)
        log.Printf("gateway: HandleDeleteAsset failed for id=%s: %v", id, err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    log.Printf("gateway: HandleDeleteAsset success for id=%s", id)
    writeJson(w, http.StatusOK, DeleteAssetResponse{Success: resp.GetSuccess()})
}