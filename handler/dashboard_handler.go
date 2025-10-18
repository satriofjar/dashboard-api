package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Booking struct {
	ID              string    `json:"id"`
	OfficeName      string    `json:"officeName"`
	RoomName        string    `json:"roomName"`
	BookingDate     time.Time `json:"bookingDate"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Participants    int       `json:"participants"`
	ListConsumption []struct {
		Name string `json:"name"`
	} `json:"listConsumption"`
}

type JenisKonsumsi struct {
	Name     string `json:"name"`
	MaxPrice int64  `json:"maxPrice"`
}

type Room struct {
	RoomName            string  `json:"roomName"`
	PersentasePemakaian float64 `json:"persentasePemakaian"`
	NominalKonsumsi     int64   `json:"nominalKonsumsi"`
	SnackSiang          int     `json:"snackSiang"`
	MakanSiang          int     `json:"makanSiang"`
	SnackSore           int     `json:"snackSore"`
}

type DashboardResponse struct {
	OfficeName string `json:"officeName"`
	Rooms      []Room `json:"rooms"`
}

func DashboardSummary(c *gin.Context) {
	periode := c.Query("periode")
	if periode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "periode id required (ex: 2024-01)",
		})
		return
	}

	bookingResp, err := http.Get("https://66876cc30bc7155dc017a662.mockapi.io/api/dummy-data/bookingList")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch booking list",
		})
		return
	}

	defer bookingResp.Body.Close()

	var bookings []Booking
	if err := json.NewDecoder(bookingResp.Body).Decode(&bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to parse booking list",
		})
		return
	}

	konsumsiReps, err := http.Get("https://6686cb5583c983911b03a7f3.mockapi.io/api/dummy-data/masterJenisKonsumsi")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch master konsumsi",
		})
		return
	}

	defer konsumsiReps.Body.Close()

	var masterKonsumsi []JenisKonsumsi
	if err := json.NewDecoder(konsumsiReps.Body).Decode(&masterKonsumsi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to parse master konsumsi",
		})
		return
	}

	priceMap := make(map[string]int64)

	for _, k := range masterKonsumsi {
		priceMap[k.Name] = k.MaxPrice
	}

	var filtered []Booking

	for _, b := range bookings {
		if b.BookingDate.Format("2006-01") == periode {
			filtered = append(filtered, b)
		}
	}

	// officeName -> roomName -> Room
	result := map[string]map[string]*Room{}

	for _, b := range filtered {
		if _, ok := result[b.OfficeName]; !ok {
			result[b.OfficeName] = map[string]*Room{}
		}

		if _, ok := result[b.OfficeName][b.RoomName]; !ok {
			result[b.OfficeName][b.RoomName] = &Room{
				RoomName: b.RoomName,
			}
		}

		room := result[b.OfficeName][b.RoomName]

		// Menambahkan jumlah pada setiap categori listConsumption dan nominalKonsumsi
		for _, cons := range b.ListConsumption {
			switch cons.Name {
			case "Snack Siang":
				room.SnackSiang += b.Participants
				room.NominalKonsumsi += int64(b.Participants) * priceMap["Snack Siang"]
			case "Makan Siang":
				room.MakanSiang += b.Participants
				room.NominalKonsumsi += int64(b.Participants) * priceMap["Makan Siang"]
			case "Snack Sore":
				room.SnackSore += b.Participants
				room.NominalKonsumsi += int64(b.Participants) * priceMap["Snack Sore"]
			}
		}

		// Persentase pemakaian diasumsikan dengan berapa lama pemakaian ruangan selama satu bulan hari kerja
		// 8 jam kerja per hari * 5 hari keja dalam seminggu * satu bulan
		// => 160 jam perbulan (20 hari kerja)
		if b.EndTime.Before(b.StartTime) {
			room.PersentasePemakaian += 0
		} else {
			duration := b.EndTime.Sub(b.StartTime).Hours()
			room.PersentasePemakaian += (duration / 160) * 100
		}
	}

	// Convert map result menjadi slice res
	var res []DashboardResponse
	for office, rooms := range result {
		rs := []Room{}
		for _, room := range rooms {
			rs = append(rs, *room)
		}

		res = append(res, DashboardResponse{
			OfficeName: office,
			Rooms:      rs,
		})

	}

	c.JSON(http.StatusOK, res)
}
