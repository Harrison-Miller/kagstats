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

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

type PlayerAvatarResponse struct {
	Large string `json:"large"`
}

func GetPlayerAvatar(player *models.Player) error {
	path := fmt.Sprintf("%s/%s/avatar", apiPath, player.Username)

	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "error getting player avatar")
	}
	defer resp.Body.Close()

	avatarResp := PlayerAvatarResponse{}

	err = json.NewDecoder(resp.Body).Decode(&avatarResp)
	if err != nil {
		return errors.Wrap(err, "error parsing player avatar response")
	}

	// set it to blank so old ones get unset
	player.Avatar = ""

	if avatarResp.Large == "" {
		return nil
	}

	// check that the link is valid
	resp2, err := client.Get(avatarResp.Large)
	if err != nil {
		return errors.Wrap(err, "error checking player avatar")
	}
	defer resp2.Body.Close()


	if resp2.StatusCode < 200 || resp2.StatusCode >= 300 {
		return fmt.Errorf("player avatar link for %s returned error code: %d", player.Username, resp2.StatusCode)
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

	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "error getting player tier")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("player tier for %s returned error code: %d", player.Username, resp.StatusCode)
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

	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "error getting player info")
	}
	defer resp.Body.Close()

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

func ValidateToken(username string, token string) error {
	// For development using two users on the site
	/*if username == "Henry" &&  token == "foobar" {
		return nil
	}*/

	path := fmt.Sprintf("https://api.kag2d.com/player/%s/token/%s", username, token)

	resp, err := client.Get(path)
	if err != nil {
		return errors.Wrap(err, "error authenticating player")
	}
	defer resp.Body.Close()
	return nil
}
