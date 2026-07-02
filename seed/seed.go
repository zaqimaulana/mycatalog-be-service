package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zaqimaulana/mycatalog-be/config"
	"github.com/zaqimaulana/mycatalog-be/models"
)

func main() {
	godotenv.Load()
	config.InitDatabase()

	products := []models.Product{
		{
			Name:        "Heineken Beer",
			Price:       45000,
			Category:    "Beer",
			Stock:       100,
			Description: "Premium lager beer from Netherlands",
			ImageURL:    "https://images.unsplash.com/photo-1608270586620-248524c67de9",
		},
		{
			Name:        "Corona Extra",
			Price:       50000,
			Category:    "Beer",
			Stock:       80,
			Description: "Light refreshing Mexican beer",
			ImageURL:    "https://images.unsplash.com/photo-1584225064536-d0fbc0a10f18",
		},
		{
			Name:        "Budweiser",
			Price:       47000,
			Category:    "Beer",
			Stock:       90,
			Description: "American classic beer",
			ImageURL:    "https://images.unsplash.com/photo-1593803431808-1a4a62a0c39d",
		},
		{
			Name:        "Guinness Stout",
			Price:       60000,
			Category:    "Beer",
			Stock:       60,
			Description: "Dark Irish stout beer",
			ImageURL:    "https://images.unsplash.com/photo-1618886614638-80e3c103d31a",
		},
		{
			Name:        "Bintang Beer",
			Price:       40000,
			Category:    "Beer",
			Stock:       120,
			Description: "Local Indonesian favorite beer",
			ImageURL:    "https://images.unsplash.com/photo-1622484212850-eb596d769edc",
		},
	}

	var created, skipped int
	for _, p := range products {
		result := config.DB.Where("name = ?", p.Name).FirstOrCreate(&p)
		if result.RowsAffected > 0 {
			created++
		} else {
			skipped++
		}
	}
	log.Printf("Seed selesai: %d ditambahkan, %d sudah ada (dilewati)", created, skipped)
}
