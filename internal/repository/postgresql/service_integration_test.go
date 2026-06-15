package postgresql

import (
	"context"
	"testJob/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntegration_SaveServiceSuccessful(t *testing.T) {
	db := integrationDB(t)
	repo := NewStorage(db)
	ctx := context.Background()

	type Record struct {
		service_name string
		price        uint
		user_id      uuid.UUID
		start_date   time.Time
		end_date     *time.Time
	}
	id, _ := uuid.Parse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	rec := Record{
		service_name: "Yandex Plus",
		price:        400,
		user_id:      id,
		start_date:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		end_date:     nil,
	}

	// 1. Сохраняем пользователя
	user := &models.User{ID: rec.user_id}
	err := repo.SaveUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}
	var gotUserID uuid.UUID
	err = db.Get(&gotUserID, "SELECT id FROM users WHERE id = $1", rec.user_id)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Сохраняем название сервиса
	err = repo.SaveServiceName(ctx, rec.service_name)
	if err != nil {
		t.Fatal(err)
	}
	var idServiceName uint
	err = db.Get(&idServiceName, "SELECT id FROM services_name WHERE name = $1", rec.service_name)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Сохраняем подписку
	service := &models.Service{
		ServiceID: idServiceName,
		Price:     rec.price,
		UserID:    rec.user_id,
		StartDate: rec.start_date,
		EndDate:   rec.end_date,
	}
	if err := repo.SaveService(ctx, service); err != nil {
		t.Fatal(err)
	}

	// 4. Проверяем, что подписка сохранилась
	services, err := repo.GetServicesByUserID(ctx, rec.user_id)
	if err != nil {
		t.Fatal(err)
	}
	if services == nil || len(*services) == 0 {
		t.Fatal("expected at least one service, got none")
	}

	got := (*services)[0]
	if got.ServiceID != idServiceName {
		t.Errorf("service_id: want %d, got %d", idServiceName, got.ServiceID)
	}
	if got.Price != rec.price {
		t.Errorf("price: want %d, got %d", rec.price, got.Price)
	}
	if got.UserID != rec.user_id {
		t.Errorf("user_id: want %s, got %s", rec.user_id, got.UserID)
	}
	if !got.StartDate.Equal(rec.start_date) {
		t.Errorf("start_date: want %v, got %v", rec.start_date, got.StartDate)
	}
	if got.EndDate != nil {
		t.Errorf("end_date: want nil, got %v", got.EndDate)
	}
}
