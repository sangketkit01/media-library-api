package handlers

import (
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/token"
	"github.com/sangketkit01/media-library-api/internal/util"
)

func (h *Handler) UploadFile(c *fiber.Ctx) error {
	p := c.Locals("payload")
	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	contentType := c.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return fiber.NewError(fiber.StatusBadRequest, "we only accept multipart/form-data")
	}

	user, err := h.Store.GetUserByID(c.Context(), pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve user")
	}

	userFolder := filepath.Join("../../uploads", user.ID.String())
	if err := os.MkdirAll(userFolder, os.ModePerm); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user folder")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid multipart form")
	}

	files := form.File["files"]
	if len(files) > 20 {
		return fiber.NewError(fiber.StatusBadRequest, "too many files (max: 20)")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var saveErrors []string

	for _, file := range files {
		if file.Size > 100<<20 { 
			saveErrors = append(saveErrors, fmt.Sprintf("file %s too large", file.Filename))
			continue
		}

		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			dst := filepath.Join(userFolder, uniqueName)

			if err := c.SaveFile(file, dst); err != nil {
				mu.Lock()
				saveErrors = append(saveErrors, fmt.Sprintf("failed to save %s: %v", file.Filename, err))
				mu.Unlock()
				return
			}

			ext := filepath.Ext(file.Filename)
			mimeTypeByExtension := mime.TypeByExtension(ext)

			arg := db.CreateMediaFileParams{
				UserID: pgtype.UUID{
					Bytes: user.ID.Bytes,
					Valid: true,
				},
				Filename: uniqueName,
				FileType: mimeTypeByExtension,
				Size:     file.Size,
			}

			if _, err := h.Store.CreateMediaFile(c.Context(), arg); err != nil {
				mu.Lock()
				saveErrors = append(saveErrors, fmt.Sprintf("DB save error for %s: %v", file.Filename, err))
				mu.Unlock()
			}
		}(file)
	}

	wg.Wait()

	if len(saveErrors) > 0 {
		return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
			"message": "Some files failed to upload",
			"errors":  saveErrors,
		})
	}

	return c.JSON(fiber.Map{"message": "Saved all files successfully!"})
}

func (h *Handler) AssignMediaToGroup(c *fiber.Ctx) error {
	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid media id")
	}

	groupID, err := uuid.Parse(c.Params("group_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid group id")
	}

	p := c.Locals("payload")
	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	_, err = h.Store.GetUserByID(c.Context(), pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive user data.")
	}

	arg := db.AssignMediaToGroupParams{
		ID: pgtype.UUID{
			Bytes: mediaID,
			Valid: true,
		},

		GroupID: pgtype.UUID{
			Bytes: groupID,
			Valid: true,
		},
	}

	err = h.Store.AssignMediaToGroup(c.Context(), arg)
	if err != nil {
		util.RouteCustomError(err, c.Path())

		return fiber.NewError(fiber.StatusInternalServerError, "failed to assign media to a group")
	}

	return c.JSON(fiber.Map{"message": "Assign media to a group successfully."})
}

func (h *Handler) GetCurrentUserMedia(c *fiber.Ctx) error {
	p := c.Locals("payload")
	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	user, err := h.Store.GetUserByID(c.Context(), pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive user data.")
	}

	groupQuery := c.Query("group_id")
	if groupQuery != "" {
		groupID, err := uuid.Parse(groupQuery)
		if err != nil {
			util.RouteCustomError(err, c.Path())
			return fiber.NewError(fiber.StatusInternalServerError, "invalid group id.")
		}

		medias, err := h.Store.ListMediaByGroup(c.Context(), db.ListMediaByGroupParams{
			UserID: pgtype.UUID{
				Bytes: user.ID.Bytes,
				Valid: true,
			},

			GroupID: pgtype.UUID{
				Bytes: groupID,
				Valid: true,
			},
		})

		if err != nil {
			util.RouteCustomError(err, c.Path())
			return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive media data.")
		}

		return c.JSON(medias)
	}

	medias, err := h.Store.ListMediaByUser(c.Context(), pgtype.UUID{
		Bytes: user.ID.Bytes,
		Valid: true,
	})

	if err != nil {
		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive media data.")
	}

	return c.JSON(medias)
}

func (h *Handler) DownloadMedia(c *fiber.Ctx) error {
	p := c.Locals("payload")
	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	user, err := h.Store.GetUserByID(c.Context(), pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive user data.")
	}

	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid media id")
	}

	media, err := h.Store.GetMediaFileByID(c.Context(), pgtype.UUID{
		Bytes: mediaID,
		Valid: true,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "media not found")
		}

		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive media.")
	}

	if media.UserID != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "You are not allowed to download other's file")
	}

	filepath := filepath.Join("../../uploads", media.UserID.String(), media.Filename)
	log.Println("file path =", filepath)
	if err := c.Download(filepath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to download media file")
	}
	return nil
}
