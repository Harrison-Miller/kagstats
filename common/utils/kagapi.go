package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/pkg/errors"
)

const patreonPath = "https://api.kag2d.com/patreon/tier"
const apiPath = "https://api.kag2d.com/v1/player"

type PlayerAvatarResponse struct {
	Large string `json:"large"`
}

func GetPlayerAvatar(player *models.Player) error {
	path := fmt.Sprintf("%s/%s/avatar", apiPath, player.Username)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "error gettitng player avatar")
	}

	avatarResp := PlayerAvatarResponse{}

	err = json.NewDecoder(resp.Body).Decode(&avatarResp)
	if err != nil {
		return errors.Wrap(err, "error parsing player avatar response")
	}

	// set it to blank so old ones get unset
	player.Avatar = ""

	// check that the link is valid
	resp, err = client.Get(avatarResp.Large)
	if err != nil {
		return errors.Wrap(err, "error checking player avatar")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("player avatar link returned error code: %d", resp.StatusCode)
	}

	player.Avatar = avatarResp.Large

	return nil
}

type PlayerTierResponse struct {
	PlayerTier struct {
		Tier int64 `json:"tier"`
	} `json:"playerTier"`
}

func GetPlayerTier(player *models.Player) error {
	path := fmt.Sprintf("%s/%s", patreonPath, player.Username)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "eror getting player tier")
	}

	tierResp := PlayerTierResponse{}

	err = json.NewDecoder(resp.Body).Decode(&tierResp)
	if err != nil {
		return errors.Wrap(err, "error parsing player tier response")
	}

	player.Tier = tierResp.PlayerTier.Tier

	return nil
}

type PlayerInfoResponse struct {
	PlayerInfo struct {
		OldGold    bool   `json:"old_gold"`
		Registered string `json:"registered"`
		Role       int64  `json:"role"`
	} `json:"playerInfo"`
}

func GetPlayerInfo(player *models.Player) error {
	path := fmt.Sprintf("%s/%s", apiPath, player.Username)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "error getting player info")
	}

	infoResp := PlayerInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&infoResp)
	if err != nil {
		return errors.Wrap(err, "error parsing player info response")
	}

	player.OldGold = infoResp.PlayerInfo.OldGold
	player.Registered = infoResp.PlayerInfo.Registered
	player.Role = infoResp.PlayerInfo.Role

	return nil
}
