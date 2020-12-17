package model

import (
	"encoding/json"
	"fmt"

	"github.com/pmoule/go2hal/hal"

	"github.com/jinzhu/gorm"
)

type TopicCollection []Topic

type Topic struct {
	CommonModelFields
	Position        GPSPosition `gorm:"embedded;embeddedPrefix:position_" json:"position"`
	Threshold       uint        `json:"threshold"`
	SubscriptionID  uint        `sql:"type:bigint REFERENCES subscriptions(id) ON DELETE CASCADE" json:"-"`
	Covid19RegionID uint        `json:"-"`
	Events          []Event     `json:"-"`
}

type GPSPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (p *GPSPosition) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}{}
	all := struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Latitude == nil {
		err = fmt.Errorf("latitude is missing")
	} else if required.Longitude == nil {
		err = fmt.Errorf("longitude is missing")
	} else {
		err = json.Unmarshal(data, &all)
		p.Longitude = all.Longitude
		p.Latitude = all.Latitude
	}
	return err
}

func NewTopic(position GPSPosition, threshold, subID uint, cov19RegID uint) (Topic, error) {
	t := Topic{
		Position:        position,
		Threshold:       threshold,
		SubscriptionID:  subID,
		Covid19RegionID: cov19RegID,
	}
	err := t.Store()
	return t, err
}

func (t *Topic) Store() error {
	return db.Save(&t).Error
}

func (t *Topic) Update(position GPSPosition, threshold uint, cov19RegID uint) error {
	t.Position = position
	t.Threshold = threshold
	t.Covid19RegionID = cov19RegID
	return t.Store()
}

func GetTopic(tID, sID uint) (*Topic, error) {
	t := &Topic{}
	err := db.Where("id = ? AND subscription_id = ?", tID, sID).First(t).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func GetTopicsBySubscriptionID(sID uint) (TopicCollection, error) {
	var tops TopicCollection
	err := db.Where("subscription_id = ?", sID).Find(&tops).Error
	return tops, err
}

func GetTopicsWithThresholdAlert(c Covid19Region) (TopicCollection, error) {
	var tops TopicCollection
	err := db.Where("covid19_region_id = ? AND threshold <= ?", c.ID, c.Incidence).Find(&tops).Error
	return tops, err
}

func (t *Topic) Delete() error {
	return db.Unscoped().Delete(t).Error
}

func (tc TopicCollection) ToHAL() hal.Resource {
	root := hal.NewResourceObject()

	// ToDo

	return root
}

func (t Topic) ToHAL() hal.Resource {
	root := hal.NewResourceObject()

	// ToDo

	return root
}
