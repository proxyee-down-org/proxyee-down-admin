package models

import (
	"common"
	"time"
)

type RecommendSoft struct {
	Id         int64     `json:"id"`
	Title      string    `json:"title"`
	Preview    string    `json:"preview"`
	Url        string    `json:"url"`
	Status     int       `json:"status"`
	Rank       int       `json:"rank"`
	CreateTime time.Time `json:"createTime"`
}

func SelectRecommendSoft() (*[]RecommendSoft, error) {
	var list []RecommendSoft
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select title,preview,url from recommend_soft where status = 1 order by rank")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var recommendSoft RecommendSoft
		err = rows.Scan(&recommendSoft.Title, &recommendSoft.Preview, &recommendSoft.Url)
		if err != nil {
			return nil, err
		} else {
			list = append(list, recommendSoft)
		}
	}
	return &list, nil
}
