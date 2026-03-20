package eventrepo

import (
	"context"

	"gorm.io/gorm"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetByID(
	ctx context.Context,
	req dto.EventGetByIDRequest,
) (dto.EventGetByIDResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.GetByID", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Event
	if err := db.Where("id = ?", string(req.ID)).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.EventGetByIDResponse{}, repository.ErrNotFound
		}
		return dto.EventGetByIDResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.EventGetByIDResponse{Event: *eventconvert.ModelToEvent(&m)}), nil
}
