package services

import (
	"context"
	"fmt"
	"log/slog"

	"firebase.google.com/go/v4/messaging"
	"github.com/zaqimaulana/mycatalog-be/config"
)

type FCMService struct{}

func NewFCMService() *FCMService { return &FCMService{} }

func (f *FCMService) SendOrderConfirmation(fcmToken string, totalAmount float64) {
	if fcmToken == "" || config.FirebaseMessaging == nil {
		return
	}
	msg := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: "Pesanan Dikonfirmasi",
			Body:  fmt.Sprintf("Pembayaran Rp %.0f berhasil. Pesananmu sedang diproses!", totalAmount),
		},
		Data: map[string]string{"type": "order_confirmation"},
	}
	if _, err := config.FirebaseMessaging.Send(context.Background(), msg); err != nil {
		slog.Warn("Gagal kirim FCM notifikasi", "error", err)
	}
}
