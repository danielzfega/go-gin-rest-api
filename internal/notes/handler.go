package notes

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    svc *NoteService
}

func NewHandler(svc *NoteService) *Handler {
    return &Handler{svc: svc}
}

func RegisterRoutes(r *gin.Engine, svc *NoteService) {
    h := NewHandler(svc)
    g := r.Group("/api/notes")
    g.POST("", h.create)
    g.GET("", h.list)
    g.GET(":id", h.get)
    g.PUT(":id", h.update)
    g.DELETE(":id", h.delete)
}

type createRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content"`
}

type updateRequest struct {
    Title   *string `json:"title"`
    Content *string `json:"content"`
}

func (h *Handler) create(c *gin.Context) {
    var req createRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    n, err := h.svc.Create(req.Title, req.Content)
    if err != nil {
        status := http.StatusBadRequest
        if err.Error() == "validation_title_required" {
            status = http.StatusBadRequest
        }
        c.JSON(status, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, n)
}

func (h *Handler) list(c *gin.Context) {
    items, err := h.svc.List()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "list_failed"})
        return
    }
    c.JSON(http.StatusOK, items)
}

func (h *Handler) get(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
        return
    }
    n, err := h.svc.Get(id)
    if err != nil {
        if err.Error() == "not_found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "get_failed"})
        return
    }
    c.JSON(http.StatusOK, n)
}

func (h *Handler) update(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
        return
    }
    var req updateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    n, err := h.svc.Update(id, req.Title, req.Content)
    if err != nil {
        if err.Error() == "validation_title_required" {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err.Error() == "not_found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
        return
    }
    c.JSON(http.StatusOK, n)
}

func (h *Handler) delete(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
        return
    }
    err = h.svc.Delete(id)
    if err != nil {
        if err.Error() == "not_found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
        return
    }
    c.Status(http.StatusNoContent)
}

